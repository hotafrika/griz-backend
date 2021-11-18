package main

import (
	"context"
	"database/sql"
	"github.com/hotafrika/griz-backend/internal/server/app"
	"github.com/hotafrika/griz-backend/internal/server/app/authtoken"
	"github.com/hotafrika/griz-backend/internal/server/app/password"
	"github.com/hotafrika/griz-backend/internal/server/app/token"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/api"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/cache/inmemory"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/database/sqlite"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	bindAddr := ":8081"
	reqTimeout := 10 * time.Second
	parseTimeout := 20 * time.Second
	logLevel := zerolog.TraceLevel
	authTokenTTL := 1800 * time.Second
	hashTTL := 1800 * time.Second
	socialLinkTTL := 1800 * time.Second
	encryptionPassString := "abc"
	encryptionAuthTokenString := "abc"
	encryptionHashString := "1234567812345678" // 16symbols

	// ENV parsing
	ba, ok := os.LookupEnv("BACKEND_ADDRESS")
	if ok {
		bindAddr = ba
	}
	rt, ok := os.LookupEnv("REQUEST_TIMEOUT")
	if ok {
		rti, err := strconv.Atoi(rt)
		if err == nil {
			reqTimeout = time.Duration(rti) * time.Second
		}
	}
	prt, ok := os.LookupEnv("PARSE_REQUEST_TIMEOUT")
	if ok {
		prti, err := strconv.Atoi(prt)
		if err == nil {
			parseTimeout = time.Duration(prti) * time.Second
		}
	}
	ll, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		lli, err := strconv.Atoi(ll)
		if err == nil {
			if lli < -1 || lli > 6 {
				log.Fatal("LOG_LEVEL must be between -1 and 6")
			}
			logLevel = zerolog.Level(lli)
		}
	}
	catt, ok := os.LookupEnv("CACHE_AUTH_TOKEN_TTL")
	if ok {
		catti, err := strconv.Atoi(catt)
		if err == nil {
			authTokenTTL = time.Duration(catti) * time.Second
		}
	}
	cht, ok := os.LookupEnv("CACHE_HASH_TTL")
	if ok {
		chti, err := strconv.Atoi(cht)
		if err == nil {
			hashTTL = time.Duration(chti) * time.Second
		}
	}
	cslt, ok := os.LookupEnv("CACHE_SOCIAL_LINK_TTL")
	if ok {
		cslti, err := strconv.Atoi(cslt)
		if err == nil {
			socialLinkTTL = time.Duration(cslti) * time.Second
		}
	}
	pek, ok := os.LookupEnv("PASSWORD_ENCRYPTION_KEY")
	if ok {
		encryptionPassString = pek
	}
	atek, ok := os.LookupEnv("AUTH_TOKEN_ENCRYPTION_KEY")
	if ok {
		encryptionAuthTokenString = atek
	}
	hek, ok := os.LookupEnv("HASH_ENCRYPTION_KEY")
	if ok {
		if len(hek) != 16 {
			log.Fatal("len of HASH_ENCRYPTION_KEY is not 16 bytes")
		}
		encryptionAuthTokenString = hek
	}

	// Work with SQL
	dbDriver := "sqlite3"
	dbString := "db.sqlite3"
	dbd, ok := os.LookupEnv("DB_DRIVER")
	if ok {
		dbDriver = dbd
	}
	dbs, ok := os.LookupEnv("DB_CONNECTION_STRING")
	if ok {
		dbString = dbs
	}
	db, err := sql.Open(dbDriver, dbString)
	if err != nil {
		panic(err)
	}

	logger := zlog.Level(logLevel)
	cache := inmemory.NewCache()

	// Inmemory repos
	//codeRepo := inmemory2.NewCodeRepository()
	//userRepo := inmemory2.NewUserRepository()

	// SQL repos
	codeRepo := sqlite.NewCodeRepository(db)
	userRepo := sqlite.NewUserRepository(db)

	passEncryptor := password.NewEncryptorByString(encryptionPassString)
	authTokenEncryptor := authtoken.NewJWTFromString(encryptionAuthTokenString, authTokenTTL)
	hashEncryptor, err := token.NewAES(encryptionHashString)
	if err != nil {
		log.Fatal("unable to initialize hash encryptor")
	}


	service := app.NewCodeService(
		authTokenTTL,
		hashTTL,
		socialLinkTTL,
		&logger,
		cache,
		codeRepo,
		userRepo,
		passEncryptor,
		authTokenEncryptor,
		hashEncryptor,
	)

	rest := api.NewRest(bindAddr, reqTimeout, parseTimeout, &logger, service)
	err = rest.Start()
	if err != nil {
		log.Fatalf("error with server: %v", err)
	}
}

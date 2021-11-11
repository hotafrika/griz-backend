package main

import (
	"github.com/hotafrika/griz-backend/internal/server/app"
	"github.com/hotafrika/griz-backend/internal/server/app/authtoken"
	"github.com/hotafrika/griz-backend/internal/server/app/password"
	"github.com/hotafrika/griz-backend/internal/server/app/token"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/api"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/cache/inmemory"
	inmemory2 "github.com/hotafrika/griz-backend/internal/server/infrastructure/database/inmemory"
	"github.com/rs/zerolog"
	"log"
	"time"
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

	logger := zerolog.Logger{}
	logger = logger.Level(logLevel)
	cache := inmemory.NewCache()
	codeRepo := inmemory2.NewCodeRepository()
	userRepo := inmemory2.NewUserRepository()
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
		logger,
		cache,
		codeRepo,
		userRepo,
		passEncryptor,
		authTokenEncryptor,
		hashEncryptor,
	)

	rest := api.NewRest(bindAddr, reqTimeout, parseTimeout, logger, service)
	err = rest.Start()
	if err != nil {
		log.Fatal("error with server")
	}
}

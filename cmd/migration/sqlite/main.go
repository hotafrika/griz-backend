package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/hotafrika/griz-backend/internal/server/app/password"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/database/sqlite"
	"github.com/pressly/goose/v3"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	dbDriver := "sqlite3"
	dbString := "db.sqlite3"
	encryptionPassString := "abc"
	initialUser := "user1"
	initialPassword := "password"
	dbd, ok := os.LookupEnv("DB_DRIVER")
	if ok {
		dbDriver = dbd
	}
	dbs, ok := os.LookupEnv("DB_CONNECTION_STRING")
	if ok {
		dbString = dbs
	}
	pek, ok := os.LookupEnv("PASSWORD_ENCRYPTION_KEY")
	if ok {
		encryptionPassString = pek
	}

	db, err := sql.Open(dbDriver, dbString)
	if err != nil {
		panic(err)
	}

	goose.SetDialect(dbDriver)
	goose.SetBaseFS(embedMigrations)

	if err = goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	// Create initial user (seed)
	userRepo := sqlite.NewUserRepository(db)
	passEncryptor := password.NewEncryptorByString(encryptionPassString)
	p, _ := passEncryptor.EncodeString(initialPassword)
	_, err = userRepo.Create(context.TODO(), entities.User{
		Username: initialUser,
		Password: string(p),
	})
	if err != nil {
		fmt.Println("error during creating user", err)
	}
}

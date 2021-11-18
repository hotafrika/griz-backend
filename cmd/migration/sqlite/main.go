package main

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
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

	goose.SetDialect(dbDriver)
	goose.SetBaseFS(embedMigrations)

	if err = goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}

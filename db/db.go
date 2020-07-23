package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var err error

func Init(dbPath string) error {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	return nil
}

func GetDb() *sql.DB {
	// todo nil?
	return db
}

func Close() {
	if db != nil {
		db.Close()
	}
}

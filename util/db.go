package util

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB  *sql.DB
	err error
)

func init() {
	DB, err = sql.Open("sqlite3", "file:netvigil.db")
	if err != nil {
		panic(err)
	}
}

package util

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	ID         string
	RemoteAddr string
	TIX        string
	Risk       int
	Reason     string
	Location   string
}

var (
	DB *sql.DB
)

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "file:record.db")
	if err != nil {
		panic(err)
	}
	DB.Exec("CREATE TABLE IF NOT EXISTS records(id TEXT PRIMARY KEY, tix TEXT, risk INTEGER, reason TEXT, location TEXT, netstat TEXT)")
	DB.Exec("INSERT INTO records(id, tix, risk, reason, location, netstat) VALUES('1', 'tix', 1, 'reason', 'location', 'netstat')")
}

// func GetRecord(id string) (Record, error) {
// 	data, err := db.Get([]byte(id), nil)
// 	if err != nil {
// 		return Record{}, err
// 	}
// 	return Record{IP: id, TIX: string(data)}, nil
// }

// func PutRecord(r Record) {
// 	db.Put(r.IP, []byte(r.TIX), nil)
// }

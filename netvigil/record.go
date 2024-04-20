package netvigil

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type RiskLevel int
type ConfidenceLevel int

type Record struct {
	Time       string
	LocalAddr  string
	RemoteAddr string
	TIX        string
	Location   string
	Reason     string
	Executable string
	Risk       RiskLevel
	Confidence ConfidenceLevel
}

const (
	Unknown RiskLevel = iota
	Safe
	Normal
	Suspicious
	Malicious
)

const (
	Low ConfidenceLevel = iota
	Medium
	High
)

var (
	DB *sql.DB
)

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "file:netvigil.db")
	if err != nil {
		panic(err)
	}
	DB.Exec("CREATE TABLE IF NOT EXISTS records (time TEXT PRIMARY KEY, local_addr TEXT, remote_addr TEXT, tix TEXT, location TEXT, reason TEXT, executable TEXT, risk INTEGER, confidence INTEGER)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_remote_addr ON records (remote_addr)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_tix ON records (tix)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_risk ON records (risk)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_confidence ON records (confidence)")
}

func (r Record) Save() error {
	if r.Risk == Suspicious {
		fmt.Printf("\x1B[33mSuspicious threat detected: %s —▸ %s\x1B[0m\n", r.Executable, r.RemoteAddr)
	} else if r.Risk == Malicious {
		fmt.Printf("\x1B[31mMalicious threat detected: %s —▸ %s\x1B[0m\n", r.Executable, r.RemoteAddr)
	}
	_, err := DB.Exec("INSERT INTO records (time, local_addr, remote_addr, tix, location, reason, executable, risk, confidence) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", r.Time, r.LocalAddr, r.RemoteAddr, r.TIX, r.Location, r.Reason, r.Executable, r.Risk, r.Confidence)
	return err
}

// func GetRecordsByRemoteAddr(remoteAddr string) ([]*Record, error) {
// 	rows, err := DB.Query("SELECT time, local_addr, remote_addr, tix, location, reason, executable, risk, confidence FROM records WHERE remote_addr = ?", remoteAddr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var records []*Record
// 	for rows.Next() {
// 		r := &Record{RemoteAddr: remoteAddr}
// 		err := rows.Scan(&r.Time, &r.LocalAddr, &r.RemoteAddr, &r.TIX, &r.Location, &r.Reason, &r.Executable, &r.Risk, &r.Confidence)
// 		if err != nil {
// 			return nil, err
// 		}
// 		records = append(records, r)
// 	}
// 	return records, nil
// }

func GetSortedRecords(sortBy string, limit int, page int) ([]*Record, error) {
	accepted := []string{"time", "remote_addr", "tix", "risk", "confidence"}
	matched := false
	sortBy = strings.ToLower(sortBy)
	if sortBy == "" {
		sortBy = "time"
	}
	for _, a := range accepted {
		if a == sortBy {
			matched = true
			break
		}
	}
	if !matched {
		return nil, errors.New("invalid sort key")
	}
	rows, err := DB.Query("SELECT time, local_addr, remote_addr, tix, location, reason, executable, risk, confidence FROM records ORDER BY "+sortBy+" LIMIT ? OFFSET ?", limit, limit*(page-1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Record
	for rows.Next() {
		r := &Record{}
		err := rows.Scan(&r.Time, &r.LocalAddr, &r.RemoteAddr, &r.TIX, &r.Location, &r.Reason, &r.Executable, &r.Risk, &r.Confidence)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

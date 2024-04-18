package util

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type RiskLevel int
type ConfidenceLevel int

type Record struct {
	ID         string
	LocalAddr  string
	RemoteAddr string
	TIX        string
	Location   string
	Reason     string
	Risk       RiskLevel
	Confidence ConfidenceLevel
}

const (
	Unknown RiskLevel = iota
	Safe
	Normal
	Suspicious
	Malicious
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
	DB.Exec("CREATE TABLE IF NOT EXISTS records (id TEXT PRIMARY KEY, local_addr TEXT, remote_addr TEXT, tix TEXT, location TEXT, reason TEXT, risk INTEGER, confidence INTEGER)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_remote_addr ON records (remote_addr)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_tix ON records (tix)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_risk ON records (risk)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_confidence ON records (confidence)")
}

func (r Record) Save() error {
	_, err := DB.Exec("INSERT INTO records (id, local_addr, remote_addr, tix, location, reason, risk, confidence) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", r.ID, r.LocalAddr, r.RemoteAddr, r.TIX, r.Location, r.Reason, r.Risk, r.Confidence)
	return err
}

func GetRecordByID(id string) (*Record, error) {
	row := DB.QueryRow("SELECT local_addr, remote_addr, tix, location, reason, risk, confidence FROM records WHERE id = ?", id)
	r := &Record{ID: id}
	err := row.Scan(&r.LocalAddr, &r.RemoteAddr, &r.TIX, &r.Location, &r.Reason, &r.Risk, &r.Confidence)
	return r, err
}

func GetRecordsByRemoteAddr(remoteAddr string) ([]*Record, error) {
	rows, err := DB.Query("SELECT id, local_addr, tix, location, reason, risk, confidence FROM records WHERE remote_addr = ?", remoteAddr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Record
	for rows.Next() {
		r := &Record{RemoteAddr: remoteAddr}
		err := rows.Scan(&r.ID, &r.LocalAddr, &r.TIX, &r.Location, &r.Reason, &r.Risk, &r.Confidence)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

func GetSortedRecords(sortBy string, limit int, page int) ([]*Record, error) {
	if sortBy != "remote_addr" && sortBy != "tix" && sortBy != "risk" && sortBy != "confidence" {
		return nil, errors.New("invalid sort key")
	}
	rows, err := DB.Query("SELECT id, local_addr, remote_addr, tix, location, reason, risk, confidence FROM records ORDER BY "+sortBy+" LIMIT ? OFFSET ?", limit, limit*(page-1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Record
	for rows.Next() {
		r := &Record{}
		err := rows.Scan(&r.ID, &r.LocalAddr, &r.RemoteAddr, &r.TIX, &r.Location, &r.Reason, &r.Risk, &r.Confidence)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

package netvigil

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type RiskLevel int
type ConfidenceLevel int

type Record struct {
	Time       int64
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
	DB.Exec("CREATE TABLE IF NOT EXISTS records (time INTEGER KEY, local_addr TEXT, remote_addr TEXT, tix TEXT, location TEXT, reason TEXT, executable TEXT, risk INTEGER, confidence INTEGER)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_remote_addr ON records (remote_addr)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_tix ON records (tix)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_risk ON records (risk)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_confidence ON records (confidence)")
}

func (r Record) Save() error {
	_, err := DB.Exec("INSERT INTO records (time, local_addr, remote_addr, tix, location, reason, executable, risk, confidence) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", r.Time, r.LocalAddr, r.RemoteAddr, r.TIX, r.Location, r.Reason, r.Executable, r.Risk, r.Confidence)
	return err
}

func GetRecords(limit int, page int) ([]*Record, error) {
	rows, err := DB.Query("SELECT time, local_addr, remote_addr, tix, location, reason, executable, risk, confidence FROM records LIMIT ? OFFSET ?", limit, limit*(page-1))
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

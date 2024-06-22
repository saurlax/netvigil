package util

import (
	"fmt"
	"strings"
)

type RiskLevel int
type CredibilityLevel int

// Threat Intelligence Record
type Threat struct {
	ID          int64            `json:"id"`
	Time        int64            `json:"time"`
	IP          string           `json:"ip"`
	TIC         string           `json:"tic"`
	Reason      string           `json:"reason"`
	Risk        RiskLevel        `json:"risk"`
	Credibility CredibilityLevel `json:"credibility"`
}

const (
	Safe       RiskLevel = iota // Confirmed to be safe, usually used for whitelists
	Normal                      // Normal IP
	Suspicious                  // Suspected of attack, IP addresses that require special attention
	Malicious                   // Confirmed as malicious IP
)

const (
	Low    CredibilityLevel = iota // Low credibility, usually used for public intelligence sources
	Medium                         // Medium credibility, usually used for private intelligence sources
	High                           // High credibility, usually used for internal intelligence sources
)

func init() {
	DB.Exec("CREATE TABLE IF NOT EXISTS threats (time INTEGER, ip TEXT, tic TEXT, reason TEXT, risk INTEGER, credibility INTEGER)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_time ON threats (time)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_ip ON threats (ip)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_tic ON threats (tic)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_risk ON threats (risk)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_credibility ON threats (credibility)")
}

func (t *Threat) Save() error {
	_, err := DB.Exec("INSERT INTO threats (time, ip, tic, reason, risk, credibility) VALUES (?, ?, ?, ?, ?, ?)", t.Time, t.IP, t.TIC, t.Reason, t.Risk, t.Credibility)
	return err
}

func GetThreats(limit int, page int) ([]*Threat, error) {
	offset := limit * (page - 1)
	rows, err := DB.Query("SELECT ROWID, time, ip, tic, reason, risk, credibility FROM threats ORDER BY time DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threats []*Threat
	for rows.Next() {
		var t Threat
		err := rows.Scan(&t.ID, &t.Time, &t.IP, &t.TIC, &t.Reason, &t.Risk, &t.Credibility)
		if err != nil {
			return nil, err
		}
		threats = append(threats, &t)
	}
	return threats, nil
}

func GetThreatsByIPs(ips []string) ([]*Threat, error) {
	if len(ips) == 0 {
		return nil, nil
	}
	// Prepare the IN clause with the correct number of placeholders
	placeholders := strings.Repeat("?,", len(ips)-1) + "?"
	query := fmt.Sprintf("SELECT ROWID, time, ip, tic, reason, risk, credibility FROM threats WHERE ip IN (%s)", placeholders)

	// Convert ips to a slice of interface{} for DB.Query
	args := make([]any, len(ips))
	for i, ip := range ips {
		args[i] = ip
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threats []*Threat
	for rows.Next() {
		var t Threat
		err := rows.Scan(&t.ID, &t.Time, &t.IP, &t.TIC, &t.Reason, &t.Risk, &t.Credibility)
		if err != nil {
			return nil, err
		}
		threats = append(threats, &t)
	}
	return threats, nil
}

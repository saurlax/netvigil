package util

type RiskLevel int
type CredibilityLevel int

// Threat Intelligence Record
type Threat struct {
	Time        int64
	IP          string
	TIC         string           // Intelligence source
	Reason      string           // The reason for the corresponding risk value
	Risk        RiskLevel        // The risk level of IP
	Credibility CredibilityLevel // The credibility level of intelligence
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
	rows, err := DB.Query("SELECT time, ip, tic, reason, risk, credibility FROM threats ORDER BY time DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threats []*Threat
	for rows.Next() {
		var t Threat
		err := rows.Scan(&t.Time, &t.IP, &t.TIC, &t.Reason, &t.Risk, &t.Credibility)
		if err != nil {
			return nil, err
		}
		threats = append(threats, &t)
	}
	return threats, nil
}

func GetThreatsByIPs(ips []string) ([]*Threat, error) {
	rows, err := DB.Query("SELECT time, ip, tic, reason, risk, credibility FROM threats WHERE ip IN (?)", ips)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threats []*Threat
	for rows.Next() {
		var t Threat
		err := rows.Scan(&t.Time, &t.IP, &t.TIC, &t.Reason, &t.Risk, &t.Credibility)
		if err != nil {
			return nil, err
		}
		threats = append(threats, &t)
	}
	return threats, nil
}

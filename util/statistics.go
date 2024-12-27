package util

import "time"

type Statistics struct {
	Time                   time.Time
	RiskUnknownCount       int64 `json:"risk_unknown_count"`
	RiskSafeCount          int64 `json:"risk_safe_count"`
	RiskNormalCount        int64 `json:"risk_normal_count"`
	RiskSuspiciousCount    int64 `json:"risk_suspicious_count"`
	RiskMaliciousCount     int64 `json:"risk_malicious_count"`
	CredibilityLowCount    int64 `json:"credibility_low_count"`
	CredibilityMediumCount int64 `json:"credibility_medium_count"`
	CredibilityHighCount   int64 `json:"credibility_high_count"`
}

var (
	statsPeriod = int64(60 * 60 * 24) // 1 day
)

// Increment the statistics if exists, otherwise insert a new record
func (stats *Statistics) Update() error {
	_, err := DB.Exec(`
	INSERT INTO statistics (
			time, 
			risk_unknown_count,
			risk_safe_count,
			risk_normal_count,
			risk_suspicious_count,
			risk_malicious_count,
			credibility_low_count,
			credibility_medium_count,
			credibility_high_count
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(time) DO UPDATE SET
			risk_unknown_count = risk_unknown_count + excluded.risk_unknown_count,
			risk_safe_count = risk_safe_count + excluded.risk_safe_count,
			risk_normal_count = risk_normal_count + excluded.risk_normal_count,
			risk_suspicious_count = risk_suspicious_count + excluded.risk_suspicious_count,
			risk_malicious_count = risk_malicious_count + excluded.risk_malicious_count,
			credibility_low_count = credibility_low_count + excluded.credibility_low_count,
			credibility_medium_count = credibility_medium_count + excluded.credibility_medium_count,
			credibility_high_count = credibility_high_count + excluded.credibility_high_count`,
		stats.Time.Unix()/statsPeriod,
		stats.RiskUnknownCount,
		stats.RiskSafeCount,
		stats.RiskNormalCount,
		stats.RiskSuspiciousCount,
		stats.RiskMaliciousCount,
		stats.CredibilityLowCount,
		stats.CredibilityMediumCount,
		stats.CredibilityHighCount,
	)
	return err
}

func init() {
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS statistics (time INTEGER UNIQUE, risk_unknown_count INTEGER, risk_safe_count INTEGER, risk_normal_count INTEGER, risk_suspicious_count INTEGER, risk_malicious_count INTEGER, credibility_low_count INTEGER, credibility_medium_count INTEGER, credibility_high_count INTEGER)")
	if err != nil {
		panic("Failed to create table statistics:" + err.Error())
	}
	_, err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_time ON statistics (time)")
	if err != nil {
		panic("Failed to create index idx_time on statistics:" + err.Error())
	}
}

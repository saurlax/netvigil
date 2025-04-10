package util

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

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

// GetSevenDayThreatPieChart 计算近七日威胁数量的饼图数据
func GetSevenDayThreatPieChart(DB *sql.DB) ([]map[string]interface{}, error) {
	now := time.Now().Unix()
	startDay := now / statsPeriod
	endDay := startDay - 6

	rows, err := DB.Query(`
		SELECT 
			COALESCE(SUM(risk_unknown_count), 0) AS riskUnkown,
			COALESCE(SUM(risk_safe_count), 0) AS riskSafe,
			COALESCE(SUM(risk_normal_count), 0) AS riskNormal,
			COALESCE(SUM(risk_suspicious_count), 0) AS riskSuspicious,
			COALESCE(SUM(risk_malicious_count), 0) AS riskMalicious
		FROM statistics
		WHERE time BETWEEN ? AND ?
	`, endDay, startDay)

	if err != nil {
		log.Println("Error fetching threat data from statistics:", err)
		return nil, fmt.Errorf("error fetching threat data from statistics: %v", err)
	}

	defer rows.Close()

	var riskUnknown, riskSafe, riskNormal, riskSuspicious, riskMalicious int64
	if rows.Next() {
		err := rows.Scan(&riskUnknown, &riskSafe, &riskNormal, &riskSuspicious, &riskMalicious)
		if err != nil {
			log.Println("Error scanning threat data:", err)
			return nil, fmt.Errorf("error scanning threat data: %v", err)
		}
	}

	// 格式化为适合饼图显示的数据格式
	pieChartData := []map[string]interface{}{
		{"name": "未知", "value": riskUnknown},
		{"name": "安全", "value": riskSafe},
		{"name": "正常", "value": riskNormal},
		{"name": "可疑", "value": riskSuspicious},
		{"name": "恶意", "value": riskMalicious},
	}
	return pieChartData, nil
}

// 可疑及以上威胁度统计数据
type SuspiciousFrequencyData struct {
	Time                     int64 `json:"time"`
	SuspiciousAboveFrequency int64 `json:"suspiciousAboveFrequency"`
}

// 地理位置统计数据
type GeoLocationData struct {
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Count     int64   `json:"count"`
}

func GetGeoLocationFrequency(DB *sql.DB) ([]GeoLocationData, error) {
	rows, err := DB.Query(`
		SELECT location, latitude, longitude, COUNT(*) as count 
		FROM netstats 
		GROUP BY location, latitude, longitude
		ORDER BY count DESC
	`)

	if err != nil {
		log.Println("Error fetching geo location threat frequency:", err)
		return nil, fmt.Errorf("error fetching geo location threat frequency: %v", err)
	}
	defer rows.Close()

	var results []GeoLocationData
	for rows.Next() {
		var data GeoLocationData
		if err := rows.Scan(&data.Name, &data.Latitude, &data.Longitude, &data.Count); err != nil {
			log.Println("Error scanning geo location threat data:", err)
			return nil, fmt.Errorf("error scanning geo location threat data: %v", err)
		}
		results = append(results, data)
	}
	return results, nil
}

// GetTicCount 获取TIC计数
func GetTicCount(DB *sql.DB) (map[string]int64, error) {
	rows, err := DB.Query(`
		SELECT tic, COUNT(*) as frequency
		FROM threats
		GROUP BY tic
		ORDER BY frequency DESC
	`)
	if err != nil {
		log.Println("Error fetching tic count:", err)
		return nil, fmt.Errorf("error fetching tic count: %v", err)
	}
	defer rows.Close()

	ticCount := make(map[string]int64)

	for rows.Next() {
		var ticName string
		var frequency int64
		err := rows.Scan(&ticName, &frequency)
		if err != nil {
			log.Println("Error scanning tic count data:", err)
			return nil, fmt.Errorf("error scanning tic count data: %v", err)
		}
		ticCount[ticName] = frequency
	}

	return ticCount, nil
}

// GetSuspiciousAboveFrequency 获取可疑及以上威胁度的频率
func GetSuspiciousAboveFrequency(DB *sql.DB) (map[int64]int64, error) {
	rows, err := DB.Query(`
		SELECT time, COALESCE(risk_suspicious_count, 0) + COALESCE(risk_malicious_count, 0) AS suspiciousAboveFrequency
		FROM statistics
		WHERE time >= (SELECT MAX(time) FROM statistics) - 7
		ORDER BY time DESC
	`)
	if err != nil {
		log.Println("Error fetching suspicious above frequency:", err)
		return nil, fmt.Errorf("error fetching suspicious above frequency: %v", err)
	}
	defer rows.Close()

	results := make(map[int64]int64)
	for rows.Next() {
		var time, suspiciousAboveFrequency int64
		err := rows.Scan(&time, &suspiciousAboveFrequency)
		if err != nil {
			log.Println("Error scanning suspicious above frequency data:", err)
			return nil, fmt.Errorf("error scanning suspicious above frequency data: %v", err)
		}

		// 转换时间戳为秒级
		results[time*statsPeriod] = suspiciousAboveFrequency
	}

	return results, nil
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

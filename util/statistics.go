package util

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

type Statistics struct {
	ID                    int64 `json:"id"`
	Time                  int64 `json:"time"` // Calculated based on timestamp / stat_period
	Risk_unknown_count    int64 `json:"risk_unknown_count"`
	Risk_safe_count       int64 `json:"risk_safe_count"`
	Risk_normal_count     int64 `json:"risk_normal_count"`
	Risk_suspicious_count int64 `json:"risk_suspicious_count"`
	Risk_Malicious_count  int64 `json:"risk_Malicious_count"`
	Credibility_low_count int64 `json:"credibility_low_count"`
	// To Add
}

var (
	StatPeriod int
	StatDB     *sql.DB
)

func init() {
	var err error
	StatDB, err = sql.Open("sqlite3", "file:statistics.db")
	if err != nil {
		log.Panicln("Failed to open statistics database:", err)
	}

	_, err = StatDB.Exec(`
		CREATE TABLE IF NOT EXISTS statistics (
			id INTEGER,
			time INTEGER,
			risk_unknown_count INTEGER,
			risk_safe_count INTEGER,
			risk_normal_count INTEGER,
			risk_suspicious_count INTEGER,
			risk_Malicious_count INTEGER,
			credibility_low_count INTEGER
			-- To Add
		)
	`)
	if err != nil {
		log.Println("Failed to create table:", err)
	}
	StatDB.Exec("CREATE INDEX IF NOT EXISTS idx_time ON statistics (time)")
}

func (s *Statistics) Save() error {
	_, err := StatDB.Exec(`
		INSERT INTO statistics (id, time, risk_unknown_count, risk_safe_count, risk_normal_count,risk_suspicious_count, risk_Malicious_count, credibility_low_count) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`,
		s.ID,
		s.Time,
		s.Risk_unknown_count,
		s.Risk_safe_count,
		s.Risk_normal_count,
		s.Risk_suspicious_count,
		s.Risk_Malicious_count,
		s.Credibility_low_count)
	return err
}

// 定期更新统计数据
func StartStatisticsJob() {
	go func() {
		statPeriod := viper.GetInt("stat_period")
		if statPeriod <= 0 {
			statPeriod = 60
			fmt.Println("The stat_period has been set as 60 seconds")
		}

		for {
			time.Sleep(time.Duration(statPeriod) * time.Second)

			timestamp := time.Now().Unix()
			timeSlot := timestamp / int64(statPeriod)

			stat := &Statistics{
				Time:                  timeSlot,
				Risk_unknown_count:    calculateRiskUnknownCount(),
				Risk_safe_count:       calculateRiskSafeCount(),
				Risk_normal_count:     calculateRiskNormalCount(),
				Risk_suspicious_count: calculateRiskSuspiciousCount(),
				Risk_Malicious_count:  calculateRiskMaliciousCount(),
				Credibility_low_count: calculateCredibilityLowCount(),
			}

			if err := stat.Save(); err != nil {
				log.Println("Error saving statistics:", err)
			} else {
				log.Printf("Statistics for time %d saved successfully\n", timeSlot)
			}
		}
	}()
}

func calculateRiskUnknownCount() int64 {
	var count int64
	DB.QueryRow(`SELECT COUNT(*) FROM threats WHERE risk_level = ?`, "Unknown").Scan(&count)
	return count
}

func calculateRiskSafeCount() int64 {
	var count int64
	DB.QueryRow(`SELECT COUNT(*) FROM threats WHERE risk_level = ?`, "Safe").Scan(&count)
	return count
}

func calculateRiskNormalCount() int64 {
	var count int64
	DB.QueryRow(`SELECT COUNT(*) FROM threats WHERE risk_level = ?`, "Normal").Scan(&count)
	return count
}

func calculateRiskSuspiciousCount() int64 {
	var count int64
	DB.QueryRow(`SELECT COUNT(*) FROM threats WHERE risk_level = ?`, "Suspicious").Scan(&count)
	return count
}

func calculateRiskMaliciousCount() int64 {
	var count int64
	DB.QueryRow(`SELECT COUNT(*) FROM threats WHERE risk_level = ?`, "Malicious").Scan(&count)
	return count
}

func calculateCredibilityLowCount() int64 {
	var count int64
	DB.QueryRow(`SELECT COUNT(*) FROM threats WHERE credibility_level = ?`, "Low").Scan(&count)
	return count
}

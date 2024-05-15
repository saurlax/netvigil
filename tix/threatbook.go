package tix

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/saurlax/netvigil/netvigil"
)

type ThreatBook struct {
	APIKey string `yaml:"apikey"`
}

type ThreatBookResult struct {
	ResponseCode int    `json:"response_code"`
	VerBoseMsg   string `json:"verbose_msg"`
	Data         map[string]struct {
		Basic struct {
			Location struct {
				Country  string `json:"country"`
				Province string `json:"province"`
				City     string `json:"city"`
			} `json:"location"`
		} `json:"basic_info"`
		Judgments       []string `json:"judgments"`
		Severity        string   `json:"severity"`
		ConfidenceLevel string   `json:"confidence_level"`
	} `json:"data"`
}

func (t *ThreatBook) Check(netstats []netstat.SockTabEntry) []netvigil.Record {
	var records []netvigil.Record
	var resource []string
	for _, v := range netstats {
		if !v.RemoteAddr.IP.IsPrivate() {
			resource = append(resource, v.RemoteAddr.IP.String())
		}
	}
	if len(resource) == 0 {
		return records
	}
	res, err := http.PostForm("https://api.threatbook.cn/v3/scene/ip_reputation", url.Values{
		"apikey":   {t.APIKey},
		"resource": resource,
	})
	if err != nil {
		fmt.Println("[Threatbook] Failed to request:", err)
		return records
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("[Threatbook] Failed to read response:", err)
		return records
	}
	var result ThreatBookResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("[Threatbook] Failed to unmarshal response:", err)
		return records
	}
	if result.ResponseCode != 0 {
		fmt.Printf("[Threatbook] Abnormal response (%v): %v", result.ResponseCode, result.VerBoseMsg)
		return records
	}
	for _, e := range netstats {
		for k, v := range result.Data {
			if e.RemoteAddr.IP.String() == k {
				var risk netvigil.RiskLevel
				switch v.Severity {
				default:
				case "low":
					risk = netvigil.Unknown
				case "info":
					risk = netvigil.Safe
				case "medium":
					risk = netvigil.Suspicious
				case "high", "critical":
					risk = netvigil.Malicious
				}
				var confidence netvigil.ConfidenceLevel
				switch v.ConfidenceLevel {
				default:
				case "low":
					confidence = netvigil.Low
				case "medium", "high":
					confidence = netvigil.Medium
				}

				records = append(records, netvigil.Record{
					Time:       time.Now().UnixMilli(),
					LocalAddr:  e.LocalAddr.String(),
					RemoteAddr: e.RemoteAddr.String(),
					TIX:        "ThreatBook",
					Reason:     strings.Join(v.Judgments, ", "),
					Executable: e.Process.Name,
					Risk:       risk,
					Confidence: confidence,
					Location:   fmt.Sprintf("%s %s %s", v.Basic.Location.Country, v.Basic.Location.Province, v.Basic.Location.City),
				})
				break
			}
		}
	}
	return records
}

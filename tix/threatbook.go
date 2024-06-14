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
	"github.com/saurlax/netvigil/util"
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
		} `json:"basic"`
		Judgments       []string `json:"judgments"`
		Severity        string   `json:"severity"`
		ConfidenceLevel string   `json:"confidence_level"`
	} `json:"data"`
}

func (t *ThreatBook) Check(netstats []netstat.SockTabEntry) []util.Record {
	var records []util.Record
	var resource []string
	for _, v := range netstats {
		if !v.RemoteAddr.IP.IsPrivate() {
			resource = append(resource, v.RemoteAddr.IP.String())
		}
	}
	if len(resource) == 0 {
		return records
	}
	// request
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
	// parse
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
	// match
	for _, e := range netstats {
		for k, v := range result.Data {
			if e.RemoteAddr.IP.String() == k {
				var risk util.RiskLevel
				switch v.Severity {
				default:
				case "low":
					risk = util.Unknown
				case "info":
					risk = util.Safe
				case "medium":
					risk = util.Suspicious
				case "high", "critical":
					risk = util.Malicious
				}
				var confidence util.ConfidenceLevel
				switch v.ConfidenceLevel {
				default:
				case "low":
					confidence = util.Low
				case "medium", "high":
					confidence = util.Medium
				}

				records = append(records, util.Record{
					Time:       time.Now().UnixMilli(),
					LocalIP:    e.LocalAddr.IP.String(),
					LocalPort:  int(e.LocalAddr.Port),
					RemoteIP:   e.RemoteAddr.IP.String(),
					RemotePort: int(e.RemoteAddr.Port),
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

package tic

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/saurlax/netvigil/util"
)

type Threatbook struct {
	APIKey string
}

type ThreatbookResponse struct {
	ResponseCode int    `json:"response_code"`
	VerBoseMsg   string `json:"verbose_msg"`
	Data         map[string]struct {
		Judgments       []string `json:"judgments"`
		Severity        string   `json:"severity"`
		ConfidenceLevel string   `json:"confidence_level"`
	} `json:"data"`
}

func (t *Threatbook) Check(ips []string) []*util.Threat {
	var threats []*util.Threat
	var resource []string

	for _, ip := range ips {
		if !net.ParseIP(ip).IsPrivate() {
			resource = append(resource, ip)
		}
	}

	resp, err := http.PostForm("https://api.threatbook.cn/v3/scene/ip_reputation", url.Values{
		"apikey":   {t.APIKey},
		"resource": resource,
	})
	if err != nil {
		log.Println("[Threatbook] Failed to request:", err)
		return threats
	}
	defer resp.Body.Close()

	var res ThreatbookResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Println("[Threatbook] Failed to decode response:", err)
		return threats
	}
	if res.ResponseCode != 0 {
		log.Printf("[Threatbook] Abnormal response (%v): %v\n", res.ResponseCode, res.VerBoseMsg)
	}

	for ip, data := range res.Data {
		var risk util.RiskLevel
		var credibility util.CredibilityLevel

		switch data.Severity {
		case "info", "low":
			risk = util.Normal
		case "high", "critical":
			risk = util.Malicious
		case "medium":
			risk = util.Suspicious
		default:
			risk = util.Normal
		}

		switch data.ConfidenceLevel {
		case "low":
			credibility = util.Low
		case "medium", "high":
			credibility = util.Medium
		default:
			credibility = util.Low
		}

		threats = append(threats, &util.Threat{
			Time:        time.Now().UnixMilli(),
			IP:          ip,
			TIC:         "Threatbook",
			Reason:      strings.Join(data.Judgments, ", "),
			Risk:        risk,
			Credibility: credibility,
		})
	}

	return threats
}

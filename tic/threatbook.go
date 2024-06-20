package tix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/saurlax/netvigil/util"
)

type Threatbook struct {
	APIKey string
}

type ThreatbookRequest struct {
	APIKey   string   `json:"apikey"`
	Resource []string `json:"resource"`
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

func (t *Threatbook) Check(ips []string) []util.Threat {
	var threats []util.Threat
	var resource []string

	for _, ip := range ips {
		if !net.ParseIP(ip).IsPrivate() {
			resource = append(resource, ip)
		}
	}

	req, err := json.Marshal(ThreatbookRequest{
		APIKey:   t.APIKey,
		Resource: resource,
	})
	resp, err := http.Post("https://api.threatbook.cn/v3/scene/ip_reputation", "application/json", bytes.NewBuffer(req))
	if err != nil {
		fmt.Println("[Threatbook] Failed to request:", err)
		return threats
	}
	defer resp.Body.Close()

	var res ThreatbookResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Println("[Threatbook] Failed to decode response:", err)
		return threats
	}
	if res.ResponseCode != 0 {
		fmt.Printf("[Threatbook] Abnormal response (%v): %v", res.ResponseCode, res.VerBoseMsg)
		return threats
	}

	for ip, data := range res.Data {
		var risk util.RiskLevel
		var credibility util.CredibilityLevel

		switch data.Severity {
		case "info":
		case "low":
			risk = util.Normal
		case "medium":
			risk = util.Suspicious
		case "high":
		case "critical":
			risk = util.Malicious
		}

		switch data.ConfidenceLevel {
		case "low":
			credibility = util.Low
		case "medium":
		case "high":
			credibility = util.High
		}

		threats = append(threats, util.Threat{
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

package tix

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

type ThreatBook struct {
	APIKey string `yaml:"apikey"`
}

type ThreatBookResult struct {
	ResponseCode int                           `json:"response_code"`
	IPs          map[string]ThreatBookResultIP `json:"ips"`
}

type ThreatBookResultIP struct {
	Severity  string                `json:"severity"`
	Judgments string                `json:"judgments"`
	ASN       ThreatBookResultIPASN `json:"asn"`
}

type ThreatBookResultIPASN struct {
	Info string `json:"info"`
}

func request(apikey string, resource []string) []IPRecord {
	var records []IPRecord
	res, err := http.PostForm("https://api.threatbook.cn/v3/scene/ip_reputation", url.Values{
		"apikey":   {apikey},
		"resource": resource,
	})
	if err != nil {
		log.Println("[Threatbook] Failed to request:", err)
		return records
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("[Threatbook] Failed to read response:", err)
		return records
	}
	var result ThreatBookResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("[Threatbook] Failed to unmarshal response:", err)
		return records
	}
	if result.ResponseCode != 0 {
		log.Println("[Threatbook] Abnormal response code:", result.ResponseCode)
	}
	for k, v := range result.IPs {
		var risk Risk
		switch v.Severity {
		case "info":
			risk = Safe
		case "low":
			risk = Low
		case "medium":
			risk = Suspicious
		case "high", "critical":
			risk = Malicious
		default:
			risk = Safe
		}
		records = append(records, IPRecord{
			IP:          net.ParseIP(k),
			Risk:        risk,
			Reason:      v.Judgments,
			Description: v.ASN.Info,
			ConfirmedBy: "Threatbook",
		})
	}
	return records
}

func (t *ThreatBook) CheckIPs(ips []net.IP) []IPRecord {
	records := make([]IPRecord, len(ips))
	var resource []string
	for _, v := range ips {
		if !v.IsPrivate() {
			resource = append(resource, v.String())
		}
	}
	var result []IPRecord
	if t.APIKey != "" {
		result = request(t.APIKey, resource)
	}
Loop:
	for i, ip := range ips {
		for _, v := range result {
			if v.IP.Equal(ip) {
				records[i] = v
				continue Loop
			}
		}
		records[i] = IPRecord{IP: ip, Risk: Safe, ConfirmedBy: "Threatbook"}
	}
	return records
}

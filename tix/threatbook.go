package tix

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/saurlax/netvigil/netvigil"
)

type ThreatBook struct {
	APIKey string `yaml:"apikey"`
}

type ThreatBookResult struct {
	ResponseCode int                           `json:"response_code"`
	VerBoseMsg   string                        `json:"verbose_msg"`
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

func request(apikey string, resource []string) []netvigil.Record {
	var records []netvigil.Record
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
		log.Printf("[Threatbook] Abnormal response (%v): %v", result.ResponseCode, result.VerBoseMsg)
	}
	for k, v := range result.IPs {
		var risk netvigil.RiskLevel
		switch v.Severity {
		case "info":
			risk = netvigil.Safe
		case "low":
			risk = netvigil.Unknown
		case "medium":
			risk = netvigil.Suspicious
		case "high", "critical":
			risk = netvigil.Malicious
		default:
			risk = netvigil.Safe
		}
		records = append(records, netvigil.Record{
			RemoteAddr: net.ParseIP(k).String(),
			Risk:       risk,
			Reason:     v.Judgments,
			Location:   v.ASN.Info,
		})
	}
	return records
}

func (t *ThreatBook) Check(ips []net.IP) []netvigil.Record {
	records := make([]netvigil.Record, len(ips))
	var resource []string
	for _, v := range ips {
		if !v.IsPrivate() {
			resource = append(resource, v.String())
		}
	}
	var result []netvigil.Record
	if t.APIKey != "" {
		result = request(t.APIKey, resource)
	}
Loop:
	for i, ip := range ips {
		for _, v := range result {
			if v.RemoteAddr == ip.String() {
				records[i] = v
				continue Loop
			}
		}
		records[i] = netvigil.Record{RemoteAddr: ip.String(), Risk: netvigil.Safe, TIX: "Threatbook"}
	}
	return records
}

package tic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/saurlax/netvigil/util"
)

type Netvigil struct {
	Server string
	APIKey string
}

type NetvigilRequest struct {
	APIKey string   `json:"apikey"`
	IPs    []string `json:"ips"`
}

func (t *Netvigil) Check(netstats []*util.Netstat) []util.Result {
	var results []util.Result
	var threats []util.Threat
	var ips []string

	for _, netstat := range netstats {
		ips = append(ips, netstat.DstIP)
	}

	requestBody, err := json.Marshal(NetvigilRequest{
		APIKey: t.APIKey,
		IPs:    ips,
	})
	if err != nil {
		log.Println("[Netvigil] Failed to marshal request:", err)
		return results
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/check", t.Server), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("[Netvigil] Failed to request:", err)
		return results
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Println("[Netvigil] Failed to request:", string(bodyBytes))
		return results
	}

	err = json.NewDecoder(resp.Body).Decode(&threats)
	if err != nil {
		log.Println("[Netvigil] Failed to decode response:", err)
		return results
	}

	for _, n := range netstats {
		for _, t := range threats {
			if n.DstIP == t.IP {
				results = append(results, util.Result{
					Time:    n.Time,
					IP:      n.DstIP,
					Netstat: n,
					Threat:  &t,
				})
			}
		}
	}

	return results
}

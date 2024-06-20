package tic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/saurlax/netvigil/util"
)

type Netvigil struct {
	Server string
	Token  string
}

type NetvigilRequest struct {
	Token string   `json:"token"`
	IPs   []string `json:"ips"`
}

type NetvigilResponse []util.Threat

func (t *Netvigil) Check(ips []string) []util.Threat {
	var threats []util.Threat
	requestBody, err := json.Marshal(NetvigilRequest{
		Token: t.Token,
		IPs:   ips,
	})
	if err != nil {
		fmt.Println("[Netvigil] Failed to marshal request:", err)
		return threats
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/check", t.Server), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("[Netvigil] Failed to request:", err)
		return threats
	}
	defer resp.Body.Close()

	var res NetvigilResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Println("[Netvigil] Failed to decode response:", err)
		return threats
	}

	return res
}

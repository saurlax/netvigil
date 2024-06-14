package tix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cakturk/go-netstat/netstat"
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

type NetvigilResponse struct {
	Records []util.Record `json:"records"`
}

func (t *Netvigil) Check(netstats []netstat.SockTabEntry) []util.Record {
	var records []util.Record
	var ips []string
	for _, v := range netstats {
		ips = append(ips, v.RemoteAddr.IP.String())
	}
	if len(ips) == 0 {
		return records
	}

	requestBody, err := json.Marshal(NetvigilRequest{
		Token: t.Token,
		IPs:   ips,
	})
	if err != nil {
		fmt.Println("[Netvigil] Failed to marshal request:", err)
		return records
	}

	res, err := http.Post(fmt.Sprintf("%s/api/check", t.Server), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("[Netvigil] Failed to request:", err)
		return records
	}
	defer res.Body.Close()

	var response NetvigilResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println("[Netvigil] Failed to decode response:", err)
		return records
	}

	return response.Records
}

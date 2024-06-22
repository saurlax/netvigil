package tic

import (
	"fmt"
	"net"

	"github.com/saurlax/netvigil/util"
	"github.com/spf13/viper"
)

// Threat Intelligence Center
type TIC interface {
	Check(ips []string) []*util.Threat
}

var tics = make([]TIC, 0)

// Create a TIC instance with config
func Create(m map[string]any) TIC {
	switch m["type"] {
	case "local":
		blacklist := make([]net.IP, 0)
		for _, v := range m["blacklist"].([]any) {
			blacklist = append(blacklist, net.ParseIP(v.(string)))
		}
		return &Local{
			Blacklist: blacklist,
		}
	case "threatbook":
		return &Threatbook{
			APIKey: m["apikey"].(string),
		}
	case "netvigil":
		return &Netvigil{
			Server: m["server"].(string),
			Token:  m["token"].(string),
		}
	default:
		return nil
	}
}

// Check all IPs with all TICs created
func CheckIPs(ips []string) []*util.Threat {
	threats, _ := util.GetThreatsByIPs(ips)
	ips2check := make([]string, 0)
Loop:
	for _, ip := range ips {
		for _, t := range threats {
			if ip == t.IP {
				continue Loop
			}
		}
		ips2check = append(ips2check, ip)
	}
	if len(ips2check) == 0 {
		return threats
	}
	for _, tic := range tics {
		for _, t := range tic.Check(ips2check) {
			t.Save()
			threats = append(threats, t)
		}
	}
	return threats
}

// Check all captured IPs
func Check() {
	ips := make([]string, 0)
Loop:
	for {
		select {
		case ip := <-util.IPs:
			ips = append(ips, ip)
		default:
			break Loop
		}
	}
	CheckIPs(ips)
}

func init() {
	config := viper.Get("tix").([]any)
	for _, v := range config {
		m, ok := v.(map[string]any)
		if !ok {
			break
		}
		tix := Create(m)
		if tix != nil {
			fmt.Printf("[TIC] %s created\n", m["type"])
			tics = append(tics, tix)
		}
	}
}

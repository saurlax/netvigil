package tic

import (
	"fmt"
	"net"

	"github.com/saurlax/netvigil/util"
	"github.com/spf13/viper"
)

// Threat Intelligence Center
type TIC interface {
	Check(ips []string) []util.Threat
}

var tics = make([]TIC, 0)

// Create a TIC instance with config like
//
// ```
// [[tix]]
// type = "local"
// xxx = "xxx"
// ```
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
func Check(ips []string) []util.Threat {
	threats := make([]util.Threat, 0)
	for _, tic := range tics {
		threats = append(threats, tic.Check(ips)...)
	}
	return threats
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

package tic

import (
	"log"
	"net"

	"github.com/google/gopacket/layers"
	"github.com/saurlax/netvigil/util"
	"github.com/spf13/viper"
)

// Threat Intelligence Center
type TIC interface {
	Check(ips []string) []*util.Threat
}

var tics = make([]TIC, 0)

// create a TIC instance with config
func create(m map[string]any) TIC {
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
func CheckAll(ips []string) []*util.Threat {
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
		case packet := <-util.Packets:
			ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
			if ipv4Layer != nil {
				ip := ipv4Layer.(*layers.IPv4)
				ips = append(ips, ip.DstIP.String())
			}
			ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
			if ipv6Layer != nil {
				ip := ipv6Layer.(*layers.IPv6)
				ips = append(ips, ip.DstIP.String())
			}
		default:
			break Loop
		}
	}
	CheckAll(ips)
}

func init() {
	config := viper.Get("tic").([]any)
	for _, v := range config {
		m, ok := v.(map[string]any)
		if !ok {
			break
		}
		tic := create(m)
		if tic != nil {
			log.Printf("[TIC] %s created\n", m["type"])
			tics = append(tics, tic)
		}
	}
}

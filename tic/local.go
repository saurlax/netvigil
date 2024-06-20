package tic

import (
	"net"
	"time"

	"github.com/saurlax/netvigil/util"
)

type Local struct {
	Blacklist []net.IP
	WhiteList []net.IP
}

func (t *Local) Check(ips []string) []util.Threat {
	threats := make([]util.Threat, 0)
	for _, ip := range ips {
		for _, black := range t.Blacklist {
			if black.Equal(net.ParseIP(ip)) {
				threats = append(threats, util.Threat{
					Time:        time.Now().UnixMilli(),
					IP:          ip,
					TIC:         "Local",
					Reason:      "Blacklisted",
					Risk:        util.Malicious,
					Credibility: util.High,
				})
			}
		}
		for _, white := range t.WhiteList {
			if white.Equal(net.ParseIP(ip)) {
				threats = append(threats, util.Threat{
					Time:        time.Now().UnixMilli(),
					IP:          ip,
					TIC:         "Local",
					Reason:      "Whitelisted",
					Risk:        util.Safe,
					Credibility: util.High,
				})
			}
		}
	}
	return threats
}

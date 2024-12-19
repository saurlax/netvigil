package tic

import (
	"log"

	"github.com/saurlax/netvigil/util"
)

type Local struct{}

func (t *Local) Check(netstats []*util.Netstat) []util.Result {
	var results []util.Result
	var ips []string

	for _, n := range netstats {
		ips = append(ips, n.DstIP)
	}

	threats, err := util.GetThreatsByIPs(ips)
	if err != nil {
		log.Println("Error getting threats:", err)
		return results
	}

	for _, n := range netstats {
		for _, t := range threats {
			if n.DstIP == t.IP {
				results = append(results, util.Result{
					Time:    n.Time,
					IP:      n.DstIP,
					Netstat: n,
					Threat:  t,
				})
			}
		}
	}

	return results
}

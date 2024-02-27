package tix

import "net"

type Local struct {
	Blacklist []net.IP `yaml:"blacklist"`
	Whitelist []net.IP `yaml:"whitelist"`
}

func (t *Local) CheckIPs(ips []net.IP) []IPRecord {
	records := make([]IPRecord, len(ips))
Loop:
	for i, ip := range ips {
		for _, v := range t.Blacklist {
			if v.Equal(ip) {
				records[i] = IPRecord{IP: ip, Risk: Malicious, Reason: "Blacklisted", ConfirmedBy: "Local"}
				continue Loop
			}
		}
		for _, v := range t.Whitelist {
			if v.Equal(ip) {
				records[i] = IPRecord{IP: ip, Risk: Whitelist, Reason: "Whitelisted", ConfirmedBy: "Local"}
				continue Loop
			}
		}
		records[i] = IPRecord{IP: ip, Risk: Safe, ConfirmedBy: "Local"}
	}
	return records
}

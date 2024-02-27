package tix

import "net"

type Local struct {
	Blacklist []net.IP
	Whitelist []net.IP
}

func (t *Local) CheckIPs(ips []net.IP) []IPRecord {
	records := make([]IPRecord, len(ips))
	for i, ip := range ips {
		for _, v := range t.Blacklist {
			if v.Equal(ip) {
				records[i] = IPRecord{IP: ip, Risk: 3, Description: "Blacklisted", ConfirmedBy: "Local"}
				continue
			}
		}
		for _, v := range t.Whitelist {
			if v.Equal(ip) {
				records[i] = IPRecord{IP: ip, Risk: 0, Description: "Whitelisted", ConfirmedBy: "Local"}
				continue
			}
		}
		records[i] = IPRecord{IP: ip, Risk: 1, ConfirmedBy: "Local"}
	}
	return records
}

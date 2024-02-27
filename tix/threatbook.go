package tix

import "net"

type ThreatBook struct {
	APIKey string `yaml:"apikey"`
}

func (t *ThreatBook) CheckIPs(ips []net.IP) []IPRecord {
	records := make([]IPRecord, len(ips))

	for i, ip := range ips {
		records[i] = IPRecord{IP: ip, Risk: 1, ConfirmedBy: "ThreatBook"}
	}
	return records
}

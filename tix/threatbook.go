package tix

import "net"

type ThreatBook struct {
	APIKey string
}

func (t *ThreatBook) CheckIPs(ips []net.IP) []IPRecord {
	return nil
}

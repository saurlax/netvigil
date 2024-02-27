package tix

import "net"

type ThreatBook struct {
	APIKey string `yaml:"apikey"`
}

func (t *ThreatBook) CheckIPs(ips []net.IP) []IPRecord {
	return nil
}

package tix

import "net"

// Threat Intelligence Center
type TIX interface {
	CheckIPs(ips []net.IP) []IPRecord
}
type IPRecord struct {
	IP   net.IP
	Risk int
	// 0 safe
	// 1 unknown / low
	// 1 suspicious
	// 2 malicious
	Description string
	// ASN    string
	// TODO: geoip
	ConfirmedBy string
}

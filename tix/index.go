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
	// 1 unknown
	// 2 suspicious
	// 3 malicious
	Description string
	// ASN    string
	// TODO: geoip
	ConfirmedBy string
}

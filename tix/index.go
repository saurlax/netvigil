package tix

import "net"

// Threat Intelligence Center
type TIX interface {
	CheckIPs(ips []net.IP) []IPRecord
}

type Risk int

const (
	Safe       Risk = 0
	Low        Risk = 1
	Suspicious Risk = 2
	Malicious  Risk = 3
	Whitelist  Risk = 10
)

type IPRecord struct {
	IP          net.IP `json:"ip"`
	Risk        Risk   `json:"risk"`
	Reason      string `json:"reason"`
	Description string `json:"description"`
	// ASN    string
	// TODO: geoip
	ConfirmedBy string `json:"confirmed_by"`
}

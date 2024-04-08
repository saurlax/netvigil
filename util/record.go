package util

import "net"

type Record struct {
	IP       net.IP
	TIX      string
	Risk     int
	Reason   string
	Location string
	Netstat  Netstat
}

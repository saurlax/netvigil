package tix

import (
	"github.com/cakturk/go-netstat/netstat"
	"github.com/saurlax/net-vigil/util"
)

type Local struct {
	Blacklist []string
	Whitelist []string
}

func (t *Local) Check(netstats []netstat.SockTabEntry) []util.Record {
	return nil
}

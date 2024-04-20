package tix

import (
	"net"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/saurlax/net-vigil/util"
)

type Local struct {
	Blacklist []net.IP
}

func (t *Local) Check(netstats []netstat.SockTabEntry) []util.Record {
	if len(t.Blacklist) == 0 {
		return nil
	}
	records := make([]util.Record, 0)
	for _, e := range netstats {
		for _, banned := range t.Blacklist {
			if e.RemoteAddr.IP.Equal(banned) {
				records = append(records, util.Record{
					LocalAddr:  e.LocalAddr.String(),
					RemoteAddr: e.RemoteAddr.String(),
					TIX:        "Local",
					Reason:     "Blacklisted",
					Executable: e.Process.Name,
					Risk:       util.Malicious,
					Confidence: util.High,
				})
			}
		}
	}
	return records
}

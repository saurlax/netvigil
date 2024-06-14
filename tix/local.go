package tix

import (
	"net"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/saurlax/netvigil/util"
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
	Loop:
		for _, banned := range t.Blacklist {
			if e.RemoteAddr.IP.Equal(banned) {
				records = append(records, util.Record{
					Time:       time.Now().UnixMilli(),
					LocalIP:    e.LocalAddr.IP.String(),
					LocalPort:  int(e.LocalAddr.Port),
					RemoteIP:   e.RemoteAddr.IP.String(),
					RemotePort: int(e.RemoteAddr.Port),
					TIX:        "Local",
					Reason:     "Blacklisted",
					Executable: e.Process.Name,
					Risk:       util.Malicious,
					Confidence: util.High,
				})
				break Loop
			}
		}
		records = append(records, util.Record{
			Time:       time.Now().UnixMilli(),
			LocalIP:    e.LocalAddr.IP.String(),
			LocalPort:  int(e.LocalAddr.Port),
			RemoteIP:   e.RemoteAddr.IP.String(),
			RemotePort: int(e.RemoteAddr.Port),
			TIX:        "Local",
			Reason:     "",
			Executable: e.Process.Name,
			Risk:       util.Unknown,
			Confidence: util.Low,
		})
	}
	return records
}

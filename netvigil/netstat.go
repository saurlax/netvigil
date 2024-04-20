package netvigil

import (
	"fmt"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/keybase/go-ps"
	"github.com/spf13/viper"
)

var NetstatCh = make(chan netstat.SockTabEntry, 1000)

var cache = make(map[string]time.Time)
var filter = func(s *netstat.SockTabEntry) bool {
	return s.State == netstat.Established && !s.RemoteAddr.IP.IsLoopback()
}

func capture() {
	for e, t := range cache {
		if time.Since(t) > 60*time.Second {
			delete(cache, e)
		}
	}

	tcps, _ := netstat.TCPSocks(filter)
	tcp6s, _ := netstat.TCP6Socks(filter)
	udps, _ := netstat.UDPSocks(filter)
	udp6s, _ := netstat.UDP6Socks(filter)
	entries := append(append(append(tcps, tcp6s...), udps...), udp6s...)

	for _, e := range entries {
		key := fmt.Sprintf("%s-%s", e.LocalAddr.String(), e.RemoteAddr.String())
		if _, ok := cache[key]; ok {
			continue
		}
		cache[key] = time.Now()
		proc, err := ps.FindProcess(e.Process.Pid)
		if err == nil && proc != nil {
			path, err := proc.Path()
			if err == nil {
				e.Process.Name = path
			}
		}
		select {
		case NetstatCh <- e:
		default:
			// break when channel is full
			return
		}
	}
}

func init() {
	go func() {
		for {
			time.Sleep(viper.GetDuration("capture_interval"))
			capture()
		}
	}()
}

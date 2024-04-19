package util

import (
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/keybase/go-ps"
	"github.com/spf13/viper"
)

var NetstatCh chan netstat.SockTabEntry

var filter = func(s *netstat.SockTabEntry) bool {
	return s.State == netstat.Established && !s.RemoteAddr.IP.IsLoopback()
}

func capture() {
	tcps, _ := netstat.TCPSocks(filter)
	tcp6s, _ := netstat.TCP6Socks(filter)
	udps, _ := netstat.UDPSocks(filter)
	udp6s, _ := netstat.UDP6Socks(filter)

	entries := append(append(append(tcps, tcp6s...), udps...), udp6s...)

	for _, e := range entries {
		proc, err := ps.FindProcess(e.Process.Pid)
		if err == nil {
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

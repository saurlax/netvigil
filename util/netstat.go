package util

import (
	"fmt"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/keybase/go-ps"
	"github.com/spf13/viper"
)

type Netstat struct {
	Time       int64  `json:"time"`
	LocalIP    string `json:"local_ip"`
	LocalPort  uint16 `json:"local_port"`
	RemoteIP   string `json:"remote_ip"`
	RemotePort uint16 `json:"remote_port"`
	Executable string `json:"executable"`
	Location   string `json:"location"`
}

var NetstatCh = make(chan Netstat, viper.GetInt("buffer_size"))

// netstat obtained from the system at different times will be duplicated, using a cache to deduplicate
//
// key: LocalIP:LcoalPort-RemoteIP:RemotePort, value: time
var cache = make(map[string]time.Time)

// In theory, no TIC can detect loopback addresses
var filter = func(s *netstat.SockTabEntry) bool {
	return s.State == netstat.Established && !s.RemoteAddr.IP.IsLoopback()
}

func capture() {
	for e, t := range cache {
		if time.Since(t) > 60*time.Second {
			delete(cache, e)
		}
	}

	// get all tcp and udp sockets
	tcps, _ := netstat.TCPSocks(filter)
	tcp6s, _ := netstat.TCP6Socks(filter)
	udps, _ := netstat.UDPSocks(filter)
	udp6s, _ := netstat.UDP6Socks(filter)
	entries := append(append(append(tcps, tcp6s...), udps...), udp6s...)

	for _, e := range entries {
		key := fmt.Sprintf("%s-%s", e.LocalAddr.String(), e.RemoteAddr.String())
		// continue if the entry is already in the cache
		if _, ok := cache[key]; ok {
			continue
		}
		cache[key] = time.Now()

		n := Netstat{
			Time:       time.Now().UnixMilli(),
			LocalIP:    e.LocalAddr.IP.String(),
			LocalPort:  e.LocalAddr.Port,
			RemoteIP:   e.RemoteAddr.IP.String(),
			RemotePort: e.RemoteAddr.Port,
		}
		// get the executable path of the process
		if e.Process != nil {
			proc, err := ps.FindProcess(e.Process.Pid)
			if err == nil && proc != nil {
				path, err := proc.Path()
				if err == nil {
					n.Executable = path
				}
			}
		}
		select {
		case NetstatCh <- n:
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

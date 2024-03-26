package util

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/keybase/go-ps"
)

type NetStat struct {
	LocalAddr  net.IP
	LocalPort  uint16
	RemoteAddr net.IP
	RemotePort uint16
	ExeName    string
	ExePath    string
	Pid        int
}

var NetstatCh chan NetStat

func capture() {
	accepted := func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Established && !s.RemoteAddr.IP.IsLoopback()
	}
	for {
		var err error
		time.Sleep(time.Duration(config.CaptureInterval) * time.Second)
		tcps, err := netstat.TCPSocks(accepted)
		if err != nil {
			log.Println("Failed to get tcp socks", err)
		}
		tcp6s, err := netstat.TCP6Socks(accepted)
		if err != nil {
			log.Println("Failed to get tcp6 socks", err)
		}
		udps, err := netstat.UDPSocks(accepted)
		if err != nil {
			log.Println("Failed to get udp socks", err)
		}
		udp6s, err := netstat.UDP6Socks(accepted)
		if err != nil {
			log.Println("Failed to get udp6 socks", err)
		}
		tabs := append(append(append(tcps, tcp6s...), udps...), udp6s...)
		log.Println("Captured", len(tabs), "sockets")
	Loop:
		for _, e := range tabs {
			proc, err := ps.FindProcess(int(e.Process.Pid))
			if err != nil {
				fmt.Println("Failed to determine process:", err)
				continue
			}
			path, _ := proc.Path()

			select {
			case sockets <- NetStatData{
				LocalAddr:  e.LocalAddr.IP,
				LocalPort:  e.LocalAddr.Port,
				RemoteAddr: e.RemoteAddr.IP,
				RemotePort: e.RemoteAddr.Port,
				ExeName:    e.Process.Name,
				Pid:        e.Process.Pid,
				// ExePath:    path,
			}:
			default:
				break Loop
			}
		}
	}
}

func init() {
	println("init")
}

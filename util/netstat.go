package util

import (
	"fmt"
	"net"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/keybase/go-ps"
)

type Netstat struct {
	LocalAddr  net.IP
	LocalPort  uint16
	RemoteAddr net.IP
	RemotePort uint16
	ExeName    string
	ExePath    string
	Pid        int
}

var NetstatCh chan Netstat

var filter = func(s *netstat.SockTabEntry) bool {
	return s.State == netstat.Established && !s.RemoteAddr.IP.IsLoopback()
}

func capture() {
	tcps, _ := netstat.TCPSocks(filter)
	tcp6s, _ := netstat.TCP6Socks(filter)
	udps, _ := netstat.UDPSocks(filter)
	udp6s, _ := netstat.UDP6Socks(filter)

	tabs := append(append(append(tcps, tcp6s...), udps...), udp6s...)
	fmt.Println("Captured %d sockets", len(tabs))

	for _, e := range tabs {
		proc, err := ps.FindProcess(e.Process.Pid)
		path := ""
		if err == nil {
			path, err = proc.Path()
		}
		select {
		case NetstatCh <- Netstat{
			LocalAddr:  e.LocalAddr.IP,
			LocalPort:  e.LocalAddr.Port,
			RemoteAddr: e.RemoteAddr.IP,
			RemotePort: e.RemoteAddr.Port,
			ExeName:    e.Process.Name,
			Pid:        e.Process.Pid,
			ExePath:    path,
		}:
		default:
			// break when channel is full
			return
		}
	}
}

func init() {
	go func() {
		for {
			time.Sleep(time.Duration(Config.CaptureInterval) * time.Second)
			capture()
		}
	}()
}

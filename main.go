package main

import (
	"fmt"
	"net"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/mitchellh/go-ps"
)

var (
	packets = make(chan NetStatData)
)

type NetStatData struct {
	LocalAddr  net.IP
	RemoteAddr net.IP
	Executable string
}

func capture(c chan<- NetStatData) {
	time.Sleep(1 * time.Second)
	for {
		tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
			return s.State == netstat.Established
		})
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		for _, e := range tabs {
			proc, err := ps.FindProcess(int(e.Process.Pid))
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			c <- NetStatData{
				LocalAddr:  e.LocalAddr.IP,
				RemoteAddr: e.RemoteAddr.IP,
				Executable: proc.Executable(),
			}
			// FIXME: not absolute path
		}
	}
}

func main() {
	go capture(packets)
	for data := range packets {
		fmt.Printf(data.LocalAddr.String() + " -> " + data.RemoteAddr.String() + " : " + data.Executable + "\n")
	}
}

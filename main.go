package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/mitchellh/go-ps"
	"github.com/saurlax/net-vigil/tix"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/yaml.v3"
)

var (
	packets = make(chan NetStatData)
	config  Config
	db      *leveldb.DB
)

type NetStatData struct {
	LocalAddr  net.IP
	LocalPort  uint16
	RemoteAddr net.IP
	RemotePort uint16
	Executable string
}

type Config struct {
	CaptureInterval int `yaml:"capture_interval"`
	CheckInterval   int `yaml:"check_interval"`
	ThreatBook      tix.ThreatBook
}

func capture() {
	for {
		time.Sleep(time.Duration(config.CaptureInterval) * time.Second)
		tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
			return s.State == netstat.Established && !s.RemoteAddr.IP.IsLoopback()
		})
		if err != nil {
			log.Println("Cannot get TCPSocks", err)
			continue
		}

		for _, e := range tabs {
			proc, err := ps.FindProcess(int(e.Process.Pid))
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			packets <- NetStatData{
				LocalAddr:  e.LocalAddr.IP,
				LocalPort:  e.LocalAddr.Port,
				RemoteAddr: e.RemoteAddr.IP,
				RemotePort: e.RemoteAddr.Port,
				Executable: proc.Executable(),
			}
			// FIXME: not absolute path
		}
	}
}

func check() {
	for {
		time.Sleep(time.Duration(config.CheckInterval) * time.Second)
		c := <-packets
		println(c.RemoteAddr.String())
	}
}

func init() {
	// config
	data, _ := os.ReadFile("config.yml")
	yaml.Unmarshal(data, &config)
	if config.CaptureInterval <= 0 {
		config.CaptureInterval = 10
	}
	if config.CheckInterval <= 0 {
		config.CheckInterval = 60
	}

	// db
	db, _ = leveldb.OpenFile("db", nil)
	defer db.Close()
}

func main() {
	log.Println("Service started")
	println(config.CaptureInterval, config.CheckInterval)
	go capture()
	check()
}

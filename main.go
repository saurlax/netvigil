package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/mitchellh/go-ps"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/yaml.v3"
)

var (
	config  Config = Config{}
	packets        = make(chan NetStatData)
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
	CaptureInterval int
	CheckInterval   int
}

type IPRecord struct {
	IP   net.IP
	Risk int
	// 0 safe
	// 1 unknown / low
	// 1 suspicious
	// 2 malicious
	Description string
	// ASN    string
	// TODO: geoip
	ConfirmedBy string
}

// Threat Intelligence Center
type TIX interface {
	CheckIPs(ips []net.IP) []IPRecord
}

func capture() {
	time.Sleep(1 * time.Second)
	for {
		tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
			return s.State == netstat.Established
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
		time.Sleep(1 * time.Second)
		c := <-packets
		println(c.RemoteAddr.String())
	}
}

func init() {
	// config
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	// db
	db, err = leveldb.OpenFile("db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func main() {
	log.Println("Service started")
	println(config.CaptureInterval, config.CheckInterval)
	go capture()
	check()
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/keybase/go-ps"
	"github.com/saurlax/net-vigil/tix"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/yaml.v3"
)

var (
	sockets chan NetStatData
	config  Config
	db      *leveldb.DB
)

type NetStatData struct {
	LocalAddr  net.IP
	LocalPort  uint16
	RemoteAddr net.IP
	RemotePort uint16
	ExeName    string
	ExePath    string
	Pid        int
}

type Config struct {
	CaptureInterval int            `yaml:"capture_interval"`
	CheckInterval   int            `yaml:"check_interval"`
	Buffer          int            `yaml:"buffer"`
	Local           tix.Local      `yaml:"local"`
	ThreatBook      tix.ThreatBook `yaml:"threatbook"`
}

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
				ExePath:    path,
			}:
			default:
				break Loop
			}
		}
	}
}

func check() {
	for {
		time.Sleep(time.Duration(config.CheckInterval) * time.Second)
		log.Println("Checking")
		var ips []net.IP
	Loop:
		for {
			select {
			case i := <-sockets:
				data, _ := db.Get(i.RemoteAddr, nil)
				if data != nil {
					var record tix.IPRecord
					err := json.Unmarshal(data, &record)
					if err != nil {
						log.Println("Failed to unmarshal record:", i.RemoteAddr, err)
						continue
					}
					if record.Risk > 1 {
						detected(record)
					}
					continue
				}
				// dedup
				for _, v := range ips {
					if v.Equal(i.RemoteAddr) {
						continue Loop
					}
				}
				ips = append(ips, i.RemoteAddr)
			default:
				break Loop
			}
		}

		records := config.ThreatBook.CheckIPs(ips)
		// TODO: more tix
		for i, v := range config.Local.CheckIPs(ips) {
			if v.Risk > records[i].Risk {
				records[i] = v
			}
		}

		// store
		for _, record := range records {
			if record.Risk > 1 {
				detected(record)
			}
			value, err := json.Marshal(record)
			if err != nil {
				log.Println("Failed to marshal record:", record.IP, err)
				continue
			}
			db.Put(record.IP, value, nil)
		}
	}
}

func detected(record tix.IPRecord) {
	log.Println("Threat detected:", record.IP, record.Reason, record.ConfirmedBy)
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

	buffer := config.Buffer
	if buffer <= 0 {
		buffer = 100
	}
	sockets = make(chan NetStatData, buffer)
}

func main() {
	log.Println("Service started")
	go capture()
	check()
}

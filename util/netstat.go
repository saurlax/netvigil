package util

import (
	"fmt"
	"net"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/keybase/go-ps"
	"github.com/spf13/viper"
)

type Netstat struct {
	ID         int64  `json:"id"`
	Time       int64  `json:"time"`
	LocalIP    string `json:"localIP"`
	LocalPort  uint16 `json:"localPort"`
	RemoteIP   string `json:"remoteIP"`
	RemotePort uint16 `json:"remotePort"`
	Executable string `json:"executable"`
	Location   string `json:"location"`
}

var IPs chan string

// netstat obtained from the system at different times will be duplicated, using a cache to deduplicate
//
// key: LocalIP:LcoalPort-RemoteIP:RemotePort, value: time
var cache = make(map[string]time.Time)

// In theory, no TIC can detect loopback addresses
var filter = func(s *netstat.SockTabEntry) bool {
	return s.State == netstat.Established && !s.RemoteAddr.IP.IsLoopback()
}

// Capture network traffic information
func Capture() {
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

		// get the location of the IP address
		ip := net.ParseIP(n.RemoteIP)
		record, err := GeoLiteCity.Lookup(ip)
		if err == nil {
			countryName := record.Country.Names["zh-CN"]
			if countryName == "" {
				countryName = record.Country.Names["en"]
			}
			cityName := record.City.Names["zh-CN"]
			if cityName == "" {
				cityName = record.City.Names["en"]
			}
			n.Location = countryName
			if cityName != "" {
				if n.Location != "" {
					n.Location += " "
				}
				n.Location += cityName
			}
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

		n.Save()

		select {
		case IPs <- n.RemoteIP:
		default:
			// break if the channel is full
			return
		}
	}
}

func init() {
	IPs = make(chan string, viper.GetInt("buffer_size"))
	DB.Exec("CREATE TABLE IF NOT EXISTS netstats (time INTEGER, local_ip TEXT, local_port INTEGER, remote_ip TEXT, remote_port INTEGER, executable TEXT, location TEXT)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_time ON netstats (time)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_local_ip ON netstats (local_ip)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_local_port ON netstats (local_port)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_remote_ip ON netstats (remote_ip)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_remote_port ON netstats (remote_port)")
}

func (n *Netstat) Save() error {
	_, err := DB.Exec("INSERT INTO netstats (time, local_ip, local_port, remote_ip, remote_port, executable, location) VALUES (?, ?, ?, ?, ?, ?, ?)", n.Time, n.LocalIP, n.LocalPort, n.RemoteIP, n.RemotePort, n.Executable, n.Location)
	return err
}

func GetNetstats(limit int, page int) ([]*Netstat, error) {
	offset := limit * (page - 1)
	rows, err := DB.Query("SELECT ROWID, time, local_ip, local_port, remote_ip, remote_port, executable, location FROM netstats ORDER BY time DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var netstats []*Netstat
	for rows.Next() {
		var n Netstat
		err := rows.Scan(&n.ID, &n.Time, &n.LocalIP, &n.LocalPort, &n.RemoteIP, &n.RemotePort, &n.Executable, &n.Location)
		if err != nil {
			return nil, err
		}
		netstats = append(netstats, &n)
	}
	return netstats, nil
}

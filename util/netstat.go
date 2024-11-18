package util

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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

// Capture network traffic information
func Capture() {

	for e, t := range cache {
		if time.Since(t) > 60*time.Second {
			delete(cache, e)
		}
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	// Open pcap device
	for _, device := range devices {
		wg.Add(1) // Increment the WaitGroup counter for each device

		// Launch a goroutine to capture packets for each device
		go func(device pcap.Interface) {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine finishes

			handle, err := pcap.OpenLive(device.Name, 1600, true, -1)
			if err != nil {
				log.Printf("Error opening device %s: %v\n", device.Name, err)
				return
			}
			defer handle.Close()

			// Set a timeout for each device capture (2 seconds)
			timer := time.NewTimer(2 * time.Second)
			defer timer.Stop()

			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			log.Printf("Starting capture on device: %s\n", device.Name)

			// Capture packets until the timeout is reached
			for {
				select {
				case packet := <-packetSource.Packets():
					ipLayer := packet.Layer(layers.LayerTypeIPv4)
					if ipLayer != nil {
						ip := ipLayer.(*layers.IPv4)

						// Skip loopback addresses (127.0.0.1)
						if ip.SrcIP.IsLoopback() || ip.DstIP.IsLoopback() {
							continue
						}

						key := fmt.Sprintf("%s-%s", ip.SrcIP.String(), ip.DstIP.String())

						// Continue if the entry is already in the cache
						if _, ok := cache[key]; ok {
							continue
						}

						log.Printf("Processing connection: %s\n", key)
						cache[key] = time.Now()

						n := Netstat{
							Time:       time.Now().UnixMilli(),
							LocalIP:    ip.SrcIP.String(),
							LocalPort:  0,
							RemoteIP:   ip.DstIP.String(),
							RemotePort: 0,
						}

						// Get the location of the IP address
						ipAddr := net.ParseIP(n.RemoteIP)
						record, err := GeoLiteCity.Lookup(ipAddr)
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

						n.Save()

						// Send the remote IP address to the channel
						select {
						case IPs <- n.RemoteIP:
						default:
							// If the channel is full, exit the goroutine
							return
						}
					}
				// Timeout, stop capturing and move to the next device
				case <-timer.C:
					log.Printf("Capture timeout reached for device %s\n", device.Name)
					return
				}
			}
		}(device) // Pass the device to the goroutine
	}

	// Wait for all goroutines to finish before exiting the function
	wg.Wait()
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

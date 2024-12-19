package util

import (
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Netstat struct {
	ID         int64  `json:"id"`
	Time       int64  `json:"time"`
	SrcIP      string `json:"srcIP"`
	SrcPort    uint16 `json:"srcPort"`
	DstIP      string `json:"dstIP"`
	DstPort    uint16 `json:"dstPort"`
	Executable string `json:"executable"`
	Location   string `json:"location"`
	Packet     gopacket.Packet
}

var (
	Netstats chan Netstat = make(chan Netstat, 1024)
)

// Capture network traffic information
func capture(ps *gopacket.PacketSource) {
	for packet := range ps.Packets() {
		n := Netstat{
			Time:   packet.Metadata().Timestamp.Unix(),
			Packet: packet,
		}

		ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ipv4Layer != nil {
			if ipv4Layer.(*layers.IPv4).DstIP.IsLoopback() {
				continue
			}
			n.SrcIP = ipv4Layer.(*layers.IPv4).SrcIP.String()
			n.DstIP = ipv4Layer.(*layers.IPv4).DstIP.String()
		} else {
			ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
			if ipv6Layer != nil {
				if ipv6Layer.(*layers.IPv6).DstIP.IsLoopback() {
					continue
				}
				n.SrcIP = ipv6Layer.(*layers.IPv6).SrcIP.String()
				n.DstIP = ipv6Layer.(*layers.IPv6).DstIP.String()
			} else {
				continue
			}
		}

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			n.SrcPort = uint16(tcpLayer.(*layers.TCP).SrcPort)
			n.DstPort = uint16(tcpLayer.(*layers.TCP).DstPort)
		} else {
			udpLayer := packet.Layer(layers.LayerTypeUDP)
			if udpLayer != nil {
				n.SrcPort = uint16(udpLayer.(*layers.UDP).SrcPort)
				n.DstPort = uint16(udpLayer.(*layers.UDP).DstPort)
			} else {
				continue
			}
		}

		// Get the location of the IP address
		ipAddr := net.ParseIP(n.DstIP)
		record, _ := GeoLiteCity.Lookup(ipAddr)
		if record != nil {
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
		Netstats <- n
	}
}

func init() {
	DB.Exec("CREATE TABLE IF NOT EXISTS netstats (time INTEGER, src_ip TEXT, src_port INTEGER, dst_ip TEXT, dst_port INTEGER, executable TEXT, location TEXT)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_time ON netstats (time)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_src_ip ON netstats (src_ip)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_src_port ON netstats (src_port)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_dst_ip ON netstats (dst_ip)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_dst_port ON netstats (dst_port)")

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, dev := range devices {
		for _, addr := range dev.Addresses {
			if !addr.IP.IsLoopback() && !addr.IP.IsUnspecified() {
				handle, err := pcap.OpenLive(dev.Name, 1600, true, -1)
				if err != nil {
					log.Fatalf("Error opening device %s: %v\n", dev.Name, err)
				} else {
					log.Printf("Capturing on device: %s\n", dev.Name)
				}
				packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
				go capture(packetSource)
				break
			}
		}
	}
}

func (n *Netstat) Save() error {
	_, err := DB.Exec("INSERT INTO netstats (time, src_ip, src_port, dst_ip, dst_port, executable, location) VALUES (?, ?, ?, ?, ?, ?, ?)", n.Time, n.SrcIP, n.SrcPort, n.DstIP, n.DstPort, n.Executable, n.Location)
	return err
}

func GetNetstats(limit int, page int) ([]*Netstat, error) {
	offset := limit * (page - 1)
	rows, err := DB.Query("SELECT ROWID, time, src_ip, src_port, dst_ip, dst_port, executable, location FROM netstats ORDER BY time DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var netstats []*Netstat
	for rows.Next() {
		var n Netstat
		err := rows.Scan(&n.ID, &n.Time, &n.SrcIP, &n.SrcPort, &n.DstIP, &n.DstPort, &n.Executable, &n.Location)
		if err != nil {
			return nil, err
		}
		netstats = append(netstats, &n)
	}
	return netstats, nil
}

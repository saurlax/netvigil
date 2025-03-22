package util

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/viper"
)

type Netstat struct {
	ID         int64   `json:"id"`
	Time       int64   `json:"time"`
	SrcIP      string  `json:"srcIP"`
	SrcPort    uint16  `json:"srcPort"`
	DstIP      string  `json:"dstIP"`
	DstPort    uint16  `json:"dstPort"`
	Executable string  `json:"executable"`
	Location   string  `json:"location"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Packet     gopacket.Packet
}

var (
	Netstats chan Netstat = make(chan Netstat, 1024)
)

// Capture network traffic information
func capture(ps *gopacket.PacketSource) {
	for packet := range ps.Packets() {
		n := Netstat{
			Time:   packet.Metadata().Timestamp.UnixMilli(),
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

		// 获取发起程序信息
		if runtime.GOOS == "windows" {
			n.Executable = findExecutableWindows(n.SrcPort)
		} else if runtime.GOOS == "linux" {
			n.Executable = findExecutableLinux(n.SrcPort)
		} else {
			n.Executable = "The OS has not been supported yet!"
		}

		// Get the location of the IP address
		ipAddr := net.ParseIP(n.DstIP)
		record, _ := GeoLiteCity.Lookup(ipAddr)
		if record != nil {
			n.Latitude = record.Location.Latitude
			n.Longitude = record.Location.Longitude

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

func findExecutableWindows(port uint16) string {
	// 执行 netstat 命令
	output, err := exec.Command("cmd", "/C", "netstat -ano").Output()
	if err != nil {
		return "Error to execute command:" + err.Error()
	}

	// 解析 netstat 输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 5 && strings.Contains(fields[1], fmt.Sprintf(":%d", port)) {
			// 获取 PID
			pid := fields[len(fields)-1]
			// 获取进程名称
			cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %s", pid))
			procOutput, err := cmd.Output()
			if err != nil {
				return ""
			}

			procLines := strings.Split(string(procOutput), "\n")
			if len(procLines) > 3 {
				return strings.Fields(procLines[3])[0] // 返回程序名称
			}
		}
	}
	return ""
}

func findExecutableLinux(port uint16) string {
	// 打开 /proc/net/tcp 文件
	file, err := os.Open("/proc/net/tcp")
	if err != nil {
		return "Error to execute command:" + err.Error()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // 跳过第一行（表头）

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		// 解析本地地址和端口
		localAddr := fields[1]
		parts := strings.Split(localAddr, ":")
		if len(parts) != 2 {
			continue
		}

		localPort, err := strconv.ParseUint(parts[1], 16, 16)
		if err != nil {
			continue
		}

		// 如果找到匹配的端口
		if uint16(localPort) == port {
			inode := fields[9]
			return getLinuxExecutableByInode(inode)
		}
	}
	return ""
}

func getLinuxExecutableByInode(inode string) string {
	// 遍历 /proc 目录
	procDir, err := os.Open("/proc")
	if err != nil {
		return "Error to Open proc:" + err.Error()
	}
	defer procDir.Close()

	pids, err := procDir.Readdirnames(-1)
	if err != nil {
		return "Error to Readdirnames:" + err.Error()
	}

	for _, pid := range pids {
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}

		fdDir := fmt.Sprintf("/proc/%s/fd", pid)
		fdFiles, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fdFiles {
			link, err := os.Readlink(fmt.Sprintf("%s/%s", fdDir, fd.Name()))
			if err != nil || !strings.Contains(link, inode) {
				continue
			}

			// 获取程序名称
			cmdPath := fmt.Sprintf("/proc/%s/comm", pid)
			cmd, err := os.ReadFile(cmdPath)
			if err != nil {
				return ""
			}

			return strings.TrimSpace(string(cmd))
		}
	}

	return ""
}

func init() {
	DB.Exec(`CREATE TABLE IF NOT EXISTS netstats (
		time INTEGER, 
		src_ip TEXT, 
		src_port INTEGER, 
		dst_ip TEXT, 
		dst_port INTEGER, 
		executable TEXT, 
		location TEXT,
		latitude REAL,
		longitude REAL
)`)
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_time ON netstats (time)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_src_ip ON netstats (src_ip)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_src_port ON netstats (src_port)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_dst_ip ON netstats (dst_ip)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_dst_port ON netstats (dst_port)")

	if viper.GetBool("capture") {
		devices, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatal(err)
		}

		for _, dev := range devices {
			for _, addr := range dev.Addresses {
				if !addr.IP.IsLoopback() && !addr.IP.IsUnspecified() {
					handle, err := pcap.OpenLive(dev.Name, 1600, true, -1)
					if err != nil {
						log.Fatalf("Error opening device %s (%s): %v\n", dev.Name, dev.Description, err)
					} else {
						log.Printf("Capturing on device: %s (%s)\n", dev.Name, dev.Description)
					}
					packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
					go capture(packetSource)
					break
				}
			}
		}
	} else {
		log.Println("Network traffic capture is disabled")
	}
}

func (n *Netstat) Save() error {
	_, err := DB.Exec("INSERT INTO netstats (time, src_ip, src_port, dst_ip, dst_port, executable, location, latitude, longitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		n.Time, n.SrcIP, n.SrcPort, n.DstIP, n.DstPort, n.Executable, n.Location, n.Latitude, n.Longitude)
	return err
}

func GetNetstats(limit int, page int) ([]*Netstat, int, error) {
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM netstats").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}

	offset := limit * (page - 1)
	rows, err := DB.Query("SELECT ROWID, time, src_ip, src_port, dst_ip, dst_port, executable, location FROM netstats ORDER BY time DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Println("Failed to query netstats:", err)
		return nil, 0, err
	}
	defer rows.Close()

	var netstats []*Netstat
	for rows.Next() {
		var n Netstat
		err := rows.Scan(&n.ID, &n.Time, &n.SrcIP, &n.SrcPort, &n.DstIP, &n.DstPort, &n.Executable, &n.Location)
		if err != nil {
			return nil, 0, err
		}
		netstats = append(netstats, &n)
	}
	return netstats, total, nil
}

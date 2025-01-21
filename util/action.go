package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

// 获取操作系统类型
func getOS() string {
	return runtime.GOOS
}

// / Windows: 检查防火墙规则是否已存在
func firewallRuleExistsWindows(ip string, direction string) bool {
	checkCmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name=netvigil_block_"+direction+"_"+ip)
	output, err := checkCmd.Output()
	if err != nil {
		log.Printf("Error checking Windows firewall rule: %v\n", err)
		return false
	}
	return strings.Contains(string(output), "Rule Name: netvigil_block_"+direction+"_"+ip)
}

// Linux: 检查 iptables 规则是否已存在
func firewallRuleExistsLinux(ip string, direction string) bool {
	var checkCmd *exec.Cmd
	if direction == "in" {
		checkCmd = exec.Command("iptables", "-C", "INPUT", "-s", ip, "-j", "DROP")
	} else {
		checkCmd = exec.Command("iptables", "-C", "OUTPUT", "-d", ip, "-j", "DROP")
	}
	err := checkCmd.Run()
	return err == nil // 如果命令成功执行，说明规则已存在
}

func AddFireWall(ip string) {
	if getOS() == "windows" {
		if !firewallRuleExistsWindows(ip, "in") {
			in := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
				"name=netvigil_block_in_"+ip, "dir=in", "action=block", "remoteip="+ip)
			in.Stdout = os.Stdout
			in.Stderr = os.Stderr
			if err := in.Run(); err != nil {
				log.Printf("Failed to add inbound firewall rule for %s: %v\n", ip, err)
			} else {
				log.Printf("Inbound firewall rule added for %s\n", ip)
			}
		} else {
			log.Printf("Inbound firewall rule for %s already exists, skipping...\n", ip)
		}

		if !firewallRuleExistsWindows(ip, "out") {
			out := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
				"name=netvigil_block_out_"+ip, "dir=out", "action=block", "remoteip="+ip)
			out.Stdout = os.Stdout
			out.Stderr = os.Stderr
			if err := out.Run(); err != nil {
				log.Printf("Failed to add outbound firewall rule for %s: %v\n", ip, err)
			} else {
				log.Printf("Outbound firewall rule added for %s\n", ip)
			}
		} else {
			log.Printf("Outbound firewall rule for %s already exists, skipping...\n", ip)
		}
	} else { // Linux
		if !firewallRuleExistsLinux(ip, "in") {
			in := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
			if err := in.Run(); err != nil {
				log.Printf("Failed to add inbound iptables rule for %s: %v\n", ip, err)
			} else {
				log.Printf("Inbound iptables rule added for %s\n", ip)
			}
		} else {
			log.Printf("Inbound iptables rule for %s already exists, skipping...\n", ip)
		}

		if !firewallRuleExistsLinux(ip, "out") {
			out := exec.Command("iptables", "-A", "OUTPUT", "-d", ip, "-j", "DROP")
			if err := out.Run(); err != nil {
				log.Printf("Failed to add outbound iptables rule for %s: %v\n", ip, err)
			} else {
				log.Printf("Outbound iptables rule added for %s\n", ip)
			}
		} else {
			log.Printf("Outbound iptables rule for %s already exists, skipping...\n", ip)
		}
	}
}

func suspiciousAction(n Netstat) {
	AddFireWall(n.DstIP)
	log.Printf("\x1B[33mSuspicious threat detected: %s → %s\x1B[0m\n", n.Executable, n.DstIP)
	beeep.Notify("Suspicious threat detected!", fmt.Sprintf("%s → %s", n.Executable, n.DstIP), "")
}

func maliciousAction(n Netstat) {
	AddFireWall(n.DstIP)
	log.Printf("\x1B[31mMalicious threat detected: %s → %s\x1B[0m\n", n.Executable, n.DstIP)
	beeep.Notify("Malicious threat detected!", fmt.Sprintf("%s → %s", n.Executable, n.DstIP), "")
}

func Action(results []*Result) {
	stats := Statistics{
		Time:                   time.Now(),
		RiskUnknownCount:       0,
		RiskSafeCount:          0,
		RiskNormalCount:        0,
		RiskSuspiciousCount:    0,
		RiskMaliciousCount:     0,
		CredibilityLowCount:    0,
		CredibilityMediumCount: 0,
		CredibilityHighCount:   0,
	}

	for _, r := range results {
		if r.Threat == nil {
			stats.RiskUnknownCount++
			continue
		}
		switch r.Threat.Risk {
		case Unknown:
			stats.RiskUnknownCount++
		case Safe:
			stats.RiskSafeCount++
		case Normal:
			stats.RiskNormalCount++
		case Suspicious:
			stats.RiskSuspiciousCount++
			suspiciousAction(*r.Netstat)
		case Malicious:
			stats.RiskMaliciousCount++
			maliciousAction(*r.Netstat)

		}
		switch r.Threat.Credibility {
		case Low:
			stats.CredibilityLowCount++
		case Medium:
			stats.CredibilityMediumCount++
		case High:
			stats.CredibilityHighCount++
		}
	}

	stats.Update()
}

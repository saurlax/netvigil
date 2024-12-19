package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gen2brain/beeep"
)

func AddFireWall(ip string) {
	in := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name=netvigil_block_in_"+ip, "dir=in", "action=block", "remoteip="+ip)
	out := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name=netvigil_block_out_"+ip, "dir=out", "action=block", "remoteip="+ip)
	in.Stdout = os.Stdout
	in.Stderr = os.Stderr
	out.Stdout = os.Stdout
	out.Stderr = os.Stderr
	in.Run()
	out.Run()
}

func DelFireWall(ip string) {
	in := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		"name=netvigil_block_in_"+ip)
	out := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		"name=netvigil_block_out_"+ip)
	in.Stdout = os.Stdout
	in.Stderr = os.Stderr
	out.Stdout = os.Stdout
	out.Stderr = os.Stderr
	in.Run()
	out.Run()
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

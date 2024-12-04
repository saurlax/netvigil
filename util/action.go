package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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

func (t Threat) Action(n Netstat) {
	switch t.Risk {
	case Suspicious:
		suspiciousAction(n)
	case Malicious:
		maliciousAction(n)
	}
}

package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gen2brain/beeep"
)

func addFireWall(ip string) {
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

// func delFireWall(ip string) {
// 	in := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
// 		"name=netvigil_block_in_"+ip)
// 	out := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
// 		"name=netvigil_block_out_"+ip)
// 	in.Stdout = os.Stdout
// 	in.Stderr = os.Stderr
// 	out.Stdout = os.Stdout
// 	out.Stderr = os.Stderr
// 	in.Run()
// 	out.Run()
// }

func suspiciousAction(r *Record) {
	addFireWall(r.RemoteIP)
	fmt.Printf("\x1B[33mSuspicious threat detected: %s → %s\x1B[0m\n", r.Executable, r.RemoteIP)
	beeep.Notify("Suspicious threat detected!", fmt.Sprintf("%s → %s", r.Executable, r.RemoteIP), "")
}

func maliciousAction(r *Record) {
	addFireWall(r.RemoteIP)
	fmt.Printf("\x1B[31mMalicious threat detected: %s → %s\x1B[0m\n", r.Executable, r.RemoteIP)
	beeep.Notify("Malicious threat detected!", fmt.Sprintf("%s → %s", r.Executable, r.RemoteIP), "")
}

func (r Record) Action() {
	switch r.Risk {
	case Suspicious:
		suspiciousAction(&r)
	case Malicious:
		maliciousAction(&r)
	}
}

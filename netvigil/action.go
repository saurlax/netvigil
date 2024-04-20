package netvigil

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	addFireWall(strings.Split(r.RemoteAddr, ":")[0])
	fmt.Printf("\x1B[33mSuspicious threat detected: %s → %s\x1B[0m\n", r.Executable, r.RemoteAddr)
	beeep.Notify("Suspicious threat detected!", fmt.Sprintf("%s → %s", r.Executable, r.RemoteAddr), "")
}

func maliciousAction(r *Record) {
	addFireWall(strings.Split(r.RemoteAddr, ":")[0])
	fmt.Printf("\x1B[31mMalicious threat detected: %s → %s\x1B[0m\n", r.Executable, r.RemoteAddr)
	beeep.Notify("Malicious threat detected!", fmt.Sprintf("%s → %s", r.Executable, r.RemoteAddr), "")
}

func (r Record) Action() {
	switch r.Risk {
	case Suspicious:
		suspiciousAction(&r)
	case Malicious:
		maliciousAction(&r)
	}
}

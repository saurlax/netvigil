package netvigil

import "fmt"

func suspiciousAction(r *Record) {
	fmt.Printf("\x1B[33mSuspicious threat detected: %s —▸ %s\x1B[0m\n", r.Executable, r.RemoteAddr)
}

func maliciousAction(r *Record) {
	fmt.Printf("\x1B[31mMalicious threat detected: %s —▸ %s\x1B[0m\n", r.Executable, r.RemoteAddr)
}

func (r Record) Action() {
	switch r.Risk {
	case Suspicious:
		suspiciousAction(&r)
	case Malicious:
		maliciousAction(&r)
	}
}

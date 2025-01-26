package util

// Result represents the result of a check
type Result struct {
	Time    int64
	IP      string
	Netstat *Netstat
	Threat  *Threat
}

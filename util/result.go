package util

type Result struct {
	Time    int64
	IP      string
	Netstat *Netstat
	Threat  *Threat
}

func (r *Result) Save() error {
	_, err := DB.Exec("INSERT INTO results (time, ip) VALUES (?, ?)", r.Time, r.IP)
	return err
}

func init() {
	DB.Exec("CREATE TABLE IF NOT EXISTS results (time INTEGER, ip TEXT)")
}

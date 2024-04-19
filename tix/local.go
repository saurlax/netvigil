package tix

import (
	"net"

	"github.com/saurlax/net-vigil/util"
)

type Local struct {
	Blacklist []net.IP `yaml:"blacklist"`
	Whitelist []net.IP `yaml:"whitelist"`
}

func (t *Local) Check(ips []net.IP) []util.Record {
	return nil
}

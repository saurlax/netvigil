package tix

import (
	"net"
	"time"

	"github.com/saurlax/net-vigil/util"
	"github.com/spf13/viper"
)

// Threat Intelligence Center
type TIX interface {
	Check(ips []net.IP) []util.Record
}

var tixs = make([]TIX, 0)

func Create(options map[string]any) TIX {
	switch options["type"] {
	case "local":
		return &Local{
			Blacklist: options["blacklist"].([]net.IP),
			Whitelist: options["whitelist"].([]net.IP),
		}
	case "threatbook":
		// TODO
		return nil
	default:
		return nil
	}
}

func cheak() {
	for {
		select {
		case e := <-util.NetstatCh:
			println(e.RemoteAddr)
		default:
			return
		}
	}
}

func init() {
	config := viper.Get("tix").([]map[string]any)
	for _, v := range config {
		tix := Create(v)
		if tix != nil {
			tixs = append(tixs, tix)
		}
	}
	go func() {
		for {
			time.Sleep(viper.GetDuration("check_interval"))
			cheak()
		}
	}()
}

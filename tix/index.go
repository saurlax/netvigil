package tix

import (
	"fmt"
	"net"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/saurlax/netvigil/netvigil"
	"github.com/spf13/viper"
)

// Threat Intelligence Center
type TIX interface {
	Check(netstats []netstat.SockTabEntry) []netvigil.Record
}

var tixs = make([]TIX, 0)

func Create(m map[string]any) TIX {
	switch m["type"] {
	case "local":
		blacklist := make([]net.IP, 0)
		for _, v := range m["blacklist"].([]any) {
			blacklist = append(blacklist, net.ParseIP(v.(string)))
		}
		return &Local{
			Blacklist: blacklist,
		}
	case "threatbook":
		// TODO
		return nil
	default:
		return nil
	}
}

func cheak() {
	entries := make([]netstat.SockTabEntry, 0)
Loop:
	for {
		select {
		case e := <-netvigil.NetstatCh:
			entries = append(entries, e)
		default:
			break Loop
		}
	}
	for _, tix := range tixs {
		records := tix.Check(entries)
		for _, record := range records {
			record.Action()
			err := record.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func init() {
	config := viper.Get("tix").([]any)
	for _, v := range config {
		tix := Create(v.(map[string]any))
		if tix != nil {
			tixs = append(tixs, tix)
		}
	}
}

func Run() {
	for {
		time.Sleep(viper.GetDuration("check_interval"))
		cheak()
	}
}

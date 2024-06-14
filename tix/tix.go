package tix

import (
	"fmt"
	"net"
	"time"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/saurlax/netvigil/util"
	"github.com/spf13/viper"
)

// Threat Intelligence Center
type TIX interface {
	Check(netstats []netstat.SockTabEntry) []util.Record
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
		return &ThreatBook{
			APIKey: m["apikey"].(string),
		}
	case "netvigil":
		return &Netvigil{
			Server: m["server"].(string),
			Token:  m["token"].(string),
		}
	default:
		return nil
	}
}

func cheak() {
	entries := make([]netstat.SockTabEntry, 0)
Loop:
	for {
		select {
		case e := <-util.NetstatCh:
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
		m, ok := v.(map[string]any)
		if !ok {
			break
		}
		tix := Create(m)
		if tix != nil {
			fmt.Printf("[TIX] %s created\n", m["type"])
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

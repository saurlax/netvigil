package tic

import (
	"log"
	"time"

	"github.com/saurlax/netvigil/util"
	"github.com/spf13/viper"
)

// Threat Intelligence Center
type TIC interface {
	Check(netstats []*util.Netstat) []util.Result
}

var tics = make([]TIC, 0)

// create a TIC instance with config
func create(m map[string]any) TIC {
	switch m["type"] {
	case "local":
		return &Local{}
	case "threatbook":
		return &Threatbook{
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

// check netstats via all TICs
func checkAll() []*util.Result {
	var netstats []*util.Netstat
	var results []*util.Result
	for len(util.Netstats) > 0 {
		ns := <-util.Netstats
		netstats = append(netstats, &ns)
	}

	for _, tic := range tics {
		for _, res := range tic.Check(netstats) {
			res.Save()
			results = append(results, &res)
			// remove the netstat that has been checked
			filtered := make([]*util.Netstat, 0)
			for _, ns := range netstats {
				if ns.DstIP != res.IP {
					filtered = append(filtered, ns)
				}
			}
			netstats = filtered
		}
	}
	return results
}

func init() {
	config := viper.Get("tic").([]any)
	for _, v := range config {
		m, ok := v.(map[string]any)
		if !ok {
			break
		}
		tic := create(m)
		if tic != nil {
			log.Printf("[TIC] %s created\n", m["type"])
			tics = append(tics, tic)
		}
	}

	if viper.GetDuration("check_period") > 0 {
		go func() {
			for {
				time.Sleep(viper.GetDuration("check_period"))
				util.Action(checkAll())
			}
		}()
	}
}

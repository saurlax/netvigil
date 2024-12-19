package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saurlax/netvigil/tic"
	"github.com/saurlax/netvigil/util"
	"github.com/spf13/viper"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
}

func init() {
	if viper.GetDuration("check_period") > 0 {
		go func() {
			for {
				time.Sleep(viper.GetDuration("check_period"))
				util.Action(tic.CheckAll())
			}
		}()
	} else {
		log.Println("[TIC] check period is not set, TIC will not run automatically")
	}
}

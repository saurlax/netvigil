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
	util.DB.Close()
}

func init() {
	if viper.GetDuration("check_interval") > 0 {
		go func() {
			for {
				tic.Check()
				time.Sleep(viper.GetDuration("check_interval"))
			}
		}()
	}
}

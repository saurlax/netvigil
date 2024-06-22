package main

import (
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
	util.DB.Close()
}

func init() {
	captureInterval := viper.GetDuration("capture_interval")
	checkInterval := viper.GetDuration("check_interval")
	if captureInterval > 0 {
		go func() {
			for {
				util.Capture()
				time.Sleep(captureInterval)
			}
		}()
	}
	if checkInterval > 0 {
		go func() {
			for {
				tic.Check()
				time.Sleep(checkInterval)
			}
		}()
	}
}

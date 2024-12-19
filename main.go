package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/saurlax/netvigil/util"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
	util.DB.Close()
}

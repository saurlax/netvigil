package main

import (
	"github.com/saurlax/netvigil/netvigil"
	"github.com/saurlax/netvigil/tix"
)

func main() {
	defer netvigil.DB.Close()
	tix.Run()
}

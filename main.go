package main

import (
	"github.com/saurlax/netvigil/tix"
	"github.com/saurlax/netvigil/util"
)

func main() {
	defer util.DB.Close()
	tix.Run()
}

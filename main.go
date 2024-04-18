package main

import (
	"github.com/saurlax/net-vigil/util"
)

func main() {
	util.Run()
	defer util.DB.Close()
}

package main

import (
	"github.com/saurlax/net-vigil/util"
)

func main() {
	defer util.DB.Close()
}

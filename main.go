package main

import (
	"github.com/jibuene/true-rand/webserver"
)

func main() {
	go webserver.StartServer()

	select {}
}

package main

import (
	"flag"
	"log"

	"codeberg.org/snonux/gos/internal/client/tui"
	config "codeberg.org/snonux/gos/internal/config/client"
)

func main() {
	configFile := flag.String("cfg", "/etc/gos.json", "The configuration file")

	conf, err := config.New(*configFile)
	if err != nil {
		log.Fatal("error building config:", err)
	}

	tui.Run(conf)
}

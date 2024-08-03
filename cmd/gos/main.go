package main

import (
	"flag"
	"log"
	"os"

	"codeberg.org/snonux/gos/internal/client/tui"
	config "codeberg.org/snonux/gos/internal/config/client"
)

func main() {
	configFile := flag.String("cfg", "/etc/gos.json", "The configuration file")

	conf, err := config.New(*configFile)
	if err != nil {
		log.Fatal("error building config:", err)
	}

	var logFD *os.File
	if conf.LogFile != "" {
		var err error
		logFD, err = os.OpenFile(conf.LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		log.SetOutput(logFD)
	}
	defer logFD.Close()

	if err := tui.Run(conf); err != nil {
		log.Fatal("error running TUI:", err)
	}
}

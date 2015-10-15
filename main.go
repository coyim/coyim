package main

import (
	"flag"
	"log"
	"os"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
)

func init() {
	if !*config.DebugFlag {
		return
	}

	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
	log.SetFlags(flags)

	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0777)
	log.SetOutput(devNull)
}

func main() {
	flag.Parse()

	u := client.NewClient()
	defer u.Close()

	if err := u.LoadConfig(*config.ConfigFile); err != nil {
		return
	}

	u.Loop()
	os.Stdout.Write([]byte("\n"))
}

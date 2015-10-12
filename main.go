package main

import (
	"flag"
	"os"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
)

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

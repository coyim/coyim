// +build !nocli

package main

import (
	"flag"
	"os"

	"github.com/twstrike/coyim/cli"
)

func main() {
	flag.Parse()

	u := cli.NewCLI()
	defer u.Close()

	if err := u.ParseConfig(); err != nil {
		u.Alert(err.Error())
		return
	}

	u.Loop()
	os.Stdout.Write([]byte("\n"))
}

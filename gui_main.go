// +build nocli

package main

import (
	"flag"
	"os"

	"github.com/twstrike/coyim/gui"
)

func main() {
	flag.Parse()

	ui := gui.NewGTK()
	ui.LoadConfig(*configFile)

	ui.Loop()
	os.Stdout.Write([]byte("\n"))
}

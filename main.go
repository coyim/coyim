package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/twstrike/coyim/config"
)

func initLog() {
	if !*config.DebugFlag {
		log.SetOutput(ioutil.Discard)
		return
	}

	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
	log.SetFlags(flags)
}

func main() {
	flag.Parse()
	initLog()

	runClient()

	os.Stdout.Write([]byte("\n"))
}

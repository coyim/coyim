package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/twstrike/coyim/config"
)

func init() {
	if !*config.DebugFlag {
		log.SetOutput(ioutil.Discard)
		return
	}

	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
	log.SetFlags(flags)
}

func main() {
	flag.Parse()

	u := newClient()
	defer u.Close()

	if err := u.LoadConfig(*config.ConfigFile); err != nil {
		return
	}

	u.Loop()
	os.Stdout.Write([]byte("\n"))
}

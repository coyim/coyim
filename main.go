package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/tor"
)

var coyimVersion = "<UNSET>"

func initLog() {
	if !*config.DebugFlag {
		log.SetOutput(ioutil.Discard)
		return
	}

	flags := log.Ldate | log.Ltime | log.Llongfile
	log.SetFlags(flags)
	log.SetPrefix("[CoyIM] ")
}

func main() {
	useBundledTor := flag.Bool("tor", true, "use bundled tor")

	flag.Parse()

	if *config.VersionFlag {
		var versionMessage = "CoyIM version " + coyimVersion + "\n"
		os.Stdout.WriteString(versionMessage)
		return
	}

	if *useBundledTor {
		or := tor.Exec()
		defer or.Process.Release()
	}

	initLog()
	runClient()
	os.Stdout.Write([]byte("\n"))
}

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"github.com/twstrike/coyim/config"
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
	flag.Parse()

	if *config.VersionFlag {
		var versionMessage = "CoyIM version " + coyimVersion + "\n"
		os.Stdout.WriteString(versionMessage)
		return
	}

	if *config.CpuProfile != "" {
		f, err := os.Create(*config.CpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	initLog()
	runClient()
	os.Stdout.Write([]byte("\n"))
}

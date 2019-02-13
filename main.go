package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/gui"
	"github.com/coyim/coyim/gui/settings"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtka"
	"github.com/coyim/gotk3adapter/pangoa"
)

var coyimVersion = "&lt;UNSET&gt;"

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

	if *config.CPUProfile != "" {
		f, err := os.Create(*config.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	_, enableMUC := os.LookupEnv("COYIM_ENABLE_MUC")
	config.MUCEnabled = enableMUC

	_, enableEFT := os.LookupEnv("COYIM_ENABLE_ENCRYPTED_FILE_TRANSFER")
	config.EncryptedFileTransferEnabled = enableEFT

	initLog()
	runClient()
	os.Stdout.Write([]byte("\n"))
}

func runClient() {
	g := gui.CreateGraphics(
		gtka.Real,
		gliba.Real,
		gdka.Real,
		pangoa.Real,
	)

	i18n.InitLocalization(gliba.Real)
	settings.InitSettings(gliba.Real)

	gui.NewGTK(coyimVersion, session.Factory, xmpp.DialerFactory, g).Loop()
}

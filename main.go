package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	log "github.com/sirupsen/logrus"

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

// BuildCommit contains which commit the build was based on
var BuildCommit = "UNKNOWN"

// BuildShortCommit contains which commit in short format the build was based on
var BuildShortCommit = "UNKNOWN"

// BuildTag contains which tag - if any - the build was based on
var BuildTag = "(no tag)"

// BuildTimestamp contains the timestamp in GMT when the build was made
var BuildTimestamp = "UNKNOWN"

func init() {
	if BuildTag == "(no tag)" || BuildTag == "" {
		coyimVersion = BuildCommit
	} else {
		coyimVersion = BuildTag
	}
}

func initLog() {
	log.SetLevel(log.InfoLevel)
	if *config.DebugFlag {
		log.SetLevel(log.DebugLevel)
	}
	if *config.TraceFlag {
		log.SetLevel(log.TraceLevel)
	}
	log.SetReportCaller(*config.DebugFunctionCalls)
}

func main() {
	flag.Parse()

	if *config.VersionFlag {
		fmt.Printf("CoyIM version %s\n", coyimVersion)
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
	_, _ = os.Stdout.Write([]byte("\n"))
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

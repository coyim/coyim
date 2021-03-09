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
	mainInit()
}

func mainInit() {
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

func printVersion() {
	fmt.Printf("CoyIM version %s (commit: %v built at: %v)\n", coyimVersion, BuildShortCommit, BuildTimestamp)
}

var osCreate = os.Create

func startProfileIfNecessary() {
	if *config.CPUProfile != "" {
		f, err := osCreate(*config.CPUProfile)
		if err != nil {
			log.WithError(err).Fatal("could not create CPU profile")
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.WithError(err).Fatal("could not start CPU profile")
		}
	}
}

func stopProfileIfNecessary() {
	if *config.CPUProfile != "" {
		pprof.StopCPUProfile()
	}
}

func printFinalNewline() {
	_, _ = os.Stdout.Write([]byte("\n"))
}

func main() {
	flag.Parse()

	if *config.VersionFlag {
		printVersion()
		return
	}

	startProfileIfNecessary()
	defer stopProfileIfNecessary()

	initLog()
	runClient()
	printFinalNewline()
}

type looper interface {
	Loop()
}

var createGTK = func(g gui.Graphics) looper {
	return gui.NewGTK(coyimVersion, session.Factory, xmpp.DialerFactory, g, hooks())
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

	createGTK(g).Loop()
}

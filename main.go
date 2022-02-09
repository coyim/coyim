package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/apenwarr/fixconsole"
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/gui"
	"github.com/coyim/coyim/gui/settings"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/sasl"
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

var logFile *os.File

func initLogFile(name string) {
	var err error
	logFile, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.WithError(err).WithField("file", name).Error("Couldn't open file for logging")
		return
	}

	log.SetOutput(logFile)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logFile.Close()
		os.Exit(1)
	}()
}

func initLog() {
	log.SetLevel(log.InfoLevel)
	if *config.DebugFlag {
		log.SetLevel(log.DebugLevel)
	}
	if *config.TraceFlag {
		log.SetLevel(log.TraceLevel)
	}
	if *config.LogFile != "" {
		initLogFile(*config.LogFile)
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

	sasl.Debug = *config.DebugFlag

	startProfileIfNecessary()
	defer stopProfileIfNecessary()

	initLog()

	log.WithField("version", coyimVersion).Info("Welcome to CoyIM!")

	initTranslations()
	runClient()
	printFinalNewline()
	if logFile != nil {
		logFile.Close()
	}
}

type looper interface {
	Loop()
}

var createGTK = func(g gui.Graphics) looper {
	return gui.NewGTK(coyimVersion, session.Factory, xmpp.DialerFactory, g, hooks(), translationsDirectory())
}

func runClient() {
	g := gui.CreateGraphics(
		gtka.Real,
		gliba.Real,
		gdka.Real,
		pangoa.Real,
		extraGraphics,
	)

	i18n.InitLocalization(gliba.Real)
	settings.InitSettings(gliba.Real)

	createGTK(g).Loop()
}

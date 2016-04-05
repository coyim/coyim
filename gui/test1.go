package gui

import (
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/gui/settings"
)

func potentiallyRunStuff() {
	if *config.RunTestsFlag {
		settings.RunTest()
	}
}

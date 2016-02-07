// +build !cli

package main

import (
	"github.com/twstrike/coyim/gui"
	"github.com/twstrike/coyim/session"
)

func runClient() {
	gui.NewGTK(coyimVersion, session.Factory).Loop()
}

// +build !cli

package main

import (
	"github.com/twstrike/coyim/gui"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/gotk3adapter/gdka"
	"github.com/twstrike/gotk3adapter/gliba"
	"github.com/twstrike/gotk3adapter/gtka"
	"github.com/twstrike/gotk3adapter/pangoa"
)

func runClient() {
	g := gui.CreateGraphics(
		gtka.Real,
		gliba.Real,
		gdka.Real,
		pangoa.Real,
	)
	gui.NewGTK(coyimVersion, session.Factory, xmpp.DialerFactory, g).Loop()
}

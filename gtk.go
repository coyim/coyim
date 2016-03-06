// +build !cli

package main

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdka"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gliba"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtka"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/pangoa"
	"github.com/twstrike/coyim/gui"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/xmpp"
)

func runClient() {
	g := gui.CreateGraphics(
		gtka.Real,
		gliba.Real,
		gdka.Real,
		pangoa.Real,
	)
	i18n.InitLocalization(gliba.Real)
	gui.NewGTK(coyimVersion, session.Factory, xmpp.DialerFactory, g).Loop()
}

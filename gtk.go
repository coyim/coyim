package main

import (
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

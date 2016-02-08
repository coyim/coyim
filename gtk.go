// +build !cli

package main

import (
	"github.com/twstrike/coyim/gui"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/xmpp"
)

func runClient() {
	gui.NewGTK(coyimVersion, session.Factory, xmpp.DialerFactory).Loop()
}

// +build cli

package main

import (
	"github.com/twstrike/coyim/cli"
	"github.com/twstrike/coyim/cli/terminal/real"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/xmpp"
)

func runClient() {
	i18n.InitLocalization(i18n.NoLocal)
	cli.NewCLI(coyimVersion, real.Factory, session.Factory, xmpp.DialerFactory).Loop()
}

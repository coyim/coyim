// +build cli

package main

import (
	"github.com/coyim/coyim/cli"
	"github.com/coyim/coyim/cli/terminal/real"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/xmpp"
)

func runClient() {
	i18n.InitLocalization(i18n.NoLocal)
	cli.NewCLI(coyimVersion, real.Factory, session.Factory, xmpp.DialerFactory).Loop()
}

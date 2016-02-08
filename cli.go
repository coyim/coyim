// +build cli

package main

import (
	"github.com/twstrike/coyim/cli"
	"github.com/twstrike/coyim/cli/terminal/real"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/xmpp"
)

func runClient() {
	cli.NewCLI(coyimVersion, real.Factory, session.Factory, xmpp.DialerFactory).Loop()
}

// +build cli

package main

import (
	"github.com/twstrike/coyim/cli"
	"github.com/twstrike/coyim/cli/terminal/real"
)

func runClient() {
	cli.NewCLI(coyimVersion, real.Factory).Loop()
}

// +build cli

package main

import "github.com/twstrike/coyim/cli"

func runClient() {
	cli.NewCLI(coyimVersion).Loop()
}

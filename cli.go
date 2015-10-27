// +build !nocli

package main

import "github.com/twstrike/coyim/cli"

func runClient() {
	cli.NewCLI().Loop()
}

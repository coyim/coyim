// +build !nocli

package main

import (
	"github.com/twstrike/coyim/cli"
	"github.com/twstrike/coyim/client"
)

func newClient() client.Client {
	return cli.NewCLI()
}

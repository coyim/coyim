// +build !nocli

package client

import "github.com/twstrike/coyim/cli"

func NewClient() Client {
	return cli.NewCLI()
}

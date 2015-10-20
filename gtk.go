// +build nocli

package main

import (
	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/gui"
)

func newClient() client.Client {
	return gui.NewGTK()
}

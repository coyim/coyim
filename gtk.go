// +build !cli

package main

import "github.com/twstrike/coyim/gui"

func runClient() {
	gui.NewGTK(coyimVersion).Loop()
}

// +build darwin

package main

import (
	"github.com/coyim/coyim/gui"
)

var hooks = noHooks

func hooks() gui.OSHooks {
	return gui.CreateOSX()
}

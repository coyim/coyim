// +build !darwin

package main

import (
	"github.com/coyim/coyim/gui"
)

var hooks = noHooks

func noHooks() gui.OSHooks {
	return &gui.NoHooks{}
}

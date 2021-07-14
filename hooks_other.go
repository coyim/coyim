// +build !darwin

package main

import (
	"github.com/coyim/coyim/gui"
)

var hooks = noHooks
var extraGraphics interface{} = nil

func noHooks() gui.OSHooks {
	return &gui.NoHooks{}
}

// +build darwin

package main

import (
	"github.com/coyim/coyim/gui"
	"github.com/coyim/gotk3osx"
)

var hooks = gui.CreateOSX
var extraGraphics interface{} = gotk3osx.Real

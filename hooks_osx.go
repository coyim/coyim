// +build darwin

package main

import (
	"github.com/coyim/coyim/gui"
	"github.com/coyim/coyim/gui/osx"
)

func hooks() gui.OSHooks {
	return osx.Create()
}

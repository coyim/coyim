// +build !darwin

package main

import (
	"github.com/coyim/coyim/gui"
)

func hooks() gui.OSHooks {
	return &gui.NoHooks{}
}

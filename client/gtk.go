// +build nocli

package client

import "github.com/twstrike/coyim/gui"

func NewClient() Client {
	return gui.NewGTK()
}

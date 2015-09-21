// +build nocli

package main

import (
	"os"

	"github.com/twstrike/go-gtk/gdk"
	"github.com/twstrike/go-gtk/gtk"
)

type UI interface {
	Loop()
}

type gtkUI struct {
	window *gtk.Window
}

func (u *gtkUI) Loop() {
	u.window.ShowAll()

	gdk.ThreadsEnter()
	gtk.Main()
	gdk.ThreadsLeave()
}

func NewGTK() UI {
	gtk.Init(&os.Args)
	gdk.ThreadsInit()

	ui := &gtkUI{
		window: gtk.NewWindow(gtk.WINDOW_TOPLEVEL),
	}

	ui.window.SetTitle("Coy")
	ui.window.Connect("destroy", gtk.MainQuit)
	ui.window.SetSizeRequest(200, 600)

	return ui
}

func main() {
	ui := NewGTK()
	ui.Loop()
}

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

	scrolledwin := gtk.NewScrolledWindow(nil, nil)
	scrolledwin.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	//Add a list widget
	rosterList := gtk.NewTreeView()
	rosterList.SetHeadersVisible(false)

	rosterList.AppendColumn(
		gtk.NewTreeViewColumnWithAttributes("user",
			gtk.NewCellRendererText(), "text", 0),
	)

	rosterModel := gtk.NewListStore(
		gtk.TYPE_STRING, // user
		gtk.TYPE_INT,    // id
	)

	iter := &gtk.TreeIter{}
	rosterModel.Append(iter)
	rosterModel.Set(iter,
		0, "alice@riseup.net",
		1, 111,
	)

	rosterModel.Append(iter)
	rosterModel.Set(iter,
		0, "bob@riseup.net",
		1, 222,
	)

	rosterList.SetModel(rosterModel)

	scrolledwin.Add(rosterList)
	ui.window.Add(scrolledwin)

	ui.window.SetTitle("Coy")
	ui.window.Connect("destroy", gtk.MainQuit)
	ui.window.SetSizeRequest(200, 600)

	return ui
}

func main() {
	ui := NewGTK()
	ui.Loop()
}

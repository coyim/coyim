// +build nocli

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

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
	menubar := initMenuBar()
	roster := initRoster()
	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(menubar, false, false, 0)
	vbox.Add(roster)
	ui.window.Add(vbox)

	ui.window.SetTitle("Coy")
	ui.window.Connect("destroy", gtk.MainQuit)
	ui.window.SetSizeRequest(200, 600)
	return ui
}

func authors() []string {
	if b, err := exec.Command("git", "log").Output(); err == nil {
		lines := strings.Split(string(b), "\n")

		var a []string
		r := regexp.MustCompile(`^Author:\s*([^ <]+).*$`)
		for _, e := range lines {
			ms := r.FindStringSubmatch(e)
			if ms == nil {
				continue
			}
			a = append(a, ms[1])
		}
		sort.Strings(a)
		var p string
		lines = []string{}
		for _, e := range a {
			if p == e {
				continue
			}
			lines = append(lines, e)
			p = e
		}
		lines = append(lines, "STRIKE Team <strike-public(AT)thoughtworks.com>")
		return lines
	}
	return []string{"STRIKE Team <strike-public@thoughtworks.com>"}
}

func aboutDialog() {
	dialog := gtk.NewAboutDialog()
	dialog.SetName("Coy IM!")
	dialog.SetProgramName("Coyim")
	dialog.SetAuthors(authors())
	// dir, _ := path.Split(os.Args[0])
	// imagefile := path.Join(dir, "../../data/coyim-logo.png")
	// pixbuf, _ := gdkpixbuf.NewFromFile(imagefile)
	// dialog.SetLogo(pixbuf)
	dialog.SetLicense(`Copyright (c) 2012 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`)
	dialog.SetWrapLicense(true)
	dialog.Run()
	dialog.Destroy()
}

func accountDialog() {
	config := loadConfig()
	dialog := gtk.NewDialog()
	dialog.SetTitle("Account Details")
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox := dialog.GetVBox()

	accountLabel := gtk.NewLabel("Account:")
	vbox.Add(accountLabel)

	accountInput := gtk.NewEntry()
	accountInput.SetText(config.Account)
	accountInput.SetEditable(true)
	vbox.Add(accountInput)

	button := gtk.NewButtonWithLabel("OK")
	button.Connect("clicked", func() {
		fmt.Println(accountInput.GetText())
		dialog.Destroy()
	})
	vbox.Add(button)

	dialog.ShowAll()
	dialog.Add(vbox)
	dialog.ShowAll()
}

func initMenuBar() *gtk.MenuBar {
	menubar := gtk.NewMenuBar()

	//Config -> Account
	cascademenu := gtk.NewMenuItemWithMnemonic("_Preference")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)
	menuitem := gtk.NewMenuItemWithMnemonic("_Account")
	menuitem.Connect("activate", accountDialog)
	submenu.Append(menuitem)

	//Help -> About
	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)
	menuitem = gtk.NewMenuItemWithMnemonic("_About")
	menuitem.Connect("activate", aboutDialog)
	submenu.Append(menuitem)
	return menubar
}

func initRoster() *gtk.ScrolledWindow {
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
	return scrolledwin
}

func main() {
	ui := NewGTK()
	ui.Loop()
}

func loadConfig() *Config {
	flag.Parse()

	if len(*configFile) == 0 {
		configFile, _ = findConfigFile(os.Getenv("HOME"))
	}

	config, err := ParseConfig(*configFile)
	if err != nil {
		log.Fatalln("Load Configfile failed:", err)
	}
	return config
}

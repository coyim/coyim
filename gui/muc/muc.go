package muc

import (
	"fmt"
	"html"

	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoi"
)

// GUI contains methods to work with muc-gui
type GUI interface {
	ShowWindow()
}

type graphics struct {
	gtk   gtki.Gtk
	glib  glibi.Glib
	gdk   gdki.Gdk
	pango pangoi.Pango
}

type counter struct {
	total  int
	online int
}

type gtkUI struct {
	g      graphics
	window gtki.Window

	accountManager *accountManager
	roster         *roster
	roomsServer    *roomsFakeServer
	roomUI         map[string]*roomUI

	builder *builder
}

var g graphics

// InitGUI initialized MUC-GUI
func InitGUI(gtkVal gtki.Gtk, glibVal glibi.Glib, gdkVal gdki.Gdk, pangoVal pangoi.Pango) GUI {
	g = graphics{
		gtk:   gtkVal,
		glib:  glibVal,
		gdk:   gdkVal,
		pango: pangoVal,
	}

	builder := newBuilder("muc")

	u := &gtkUI{
		g:       g,
		builder: builder,
		roomUI:  map[string]*roomUI{},
	}

	u.init()

	return u
}

const (
	indexJid               = 0
	indexAccountID         = 2
	indexBackgroundColor   = 4
	indexWeight            = 5
	indexParentJid         = 0
	indexParentDisplayName = 1
	indexStatusIcon        = 7
	indexRowType           = 8
)

func (u *gtkUI) ShowWindow() {
	u.addAccountsToRoster()
	win := u.builder.get("mainWindow").(gtki.Window)
	u.doInUIThread(win.Show)
	u.window = win
}

type peerStatus string

var (
	statusConnecting peerStatus = "connecting"
	statusOnline     peerStatus = "online"
	statusOffline    peerStatus = "offline"
)

func (u *gtkUI) init() {
	u.initRooms()

	u.initRoster()

	u.initDemoAccounts()

	u.builder.ConnectSignals(map[string]interface{}{
		"on_activate_buddy": u.onActivateRosterRow,
		"on_button_press":   u.onButtonPress,
	})
}

func (u *gtkUI) onButtonPress(view gtki.TreeView, ev gdki.Event) bool {
	return false
}

func setValues(v gtki.ListStore, iter gtki.TreeIter, values ...interface{}) {
	for i, val := range values {
		if val != nil {
			_ = v.SetValue(iter, i, val)
		}
	}
}

func decideColorForPeer(cs colorSet, i *rosterItem) string {
	if !i.isOnline() {
		return cs.rosterPeerOfflineForeground
	}
	return cs.rosterPeerOnlineForeground
}

func createTooltipForPeer(i *rosterItem) string {
	pname := html.EscapeString(i.displayName())
	jid := html.EscapeString(i.id)
	if pname != jid {
		return fmt.Sprintf("%s (%s)", pname, jid)
	}
	return jid
}

func getFromModelIterMUC(m gtki.ListStore, iter gtki.TreeIter, index int) string {
	val, _ := m.GetValue(iter, index)
	v, _ := val.GetString()
	return v
}

func createGroupDisplayName(parentName string, counter *counter, isExpanded bool) string {
	name := parentName
	if !isExpanded {
		name = fmt.Sprintf("[%s]", name)
	}
	return fmt.Sprintf("%s (%d/%d)", name, counter.online, counter.total)
}

func (c *counter) inc(total, online bool) {
	if total {
		c.total++
	}
	if online {
		c.online++
	}
}

func (u *gtkUI) doInUIThread(f func()) {
	_, _ = g.glib.IdleAdd(f)
}

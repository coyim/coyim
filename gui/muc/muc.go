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

type mucUI struct {
	g      graphics
	window gtki.Window

	accountManager *mucAccountManager
	roster         *mucRoster
	roomsServer    *mucRoomsFakeServer

	main gtki.Box `gtk-widget:"main"`

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

	m := &mucUI{
		g:       g,
		builder: builder,
	}

	m.init()

	m.addAccountsToRoster()

	return m
}

func (m *mucUI) ShowWindow() {
	m.showWindow()
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

func (m *mucUI) showWindow() {
	win := m.builder.get("mainWindow").(gtki.Window)
	win.Show()
	m.window = win
}

type mucPeerStatus string

var (
	mucStatusConnecting mucPeerStatus = "connecting"
	mucStatusOnline     mucPeerStatus = "online"
	mucStatusOffline    mucPeerStatus = "offline"
)

func (m *mucUI) init() {
	m.initRooms()

	m.initRoster()

	m.initDemoAccounts()

	panicOnDevError(m.builder.bindObjects(m))

	m.builder.ConnectSignals(map[string]interface{}{
		"on_activate_buddy": m.onActivateRosterRow,
		"on_button_press":   m.onButtonPress,
	})
}

func (m *mucUI) onButtonPress(view gtki.TreeView, ev gdki.Event) bool {
	return false
}

func setValues(v gtki.ListStore, iter gtki.TreeIter, values ...interface{}) {
	for i, val := range values {
		if val != nil {
			_ = v.SetValue(iter, i, val)
		}
	}
}

func decideColorForPeer(cs colorSet, i *mucRosterItem) string {
	if !i.isOnline() {
		return cs.rosterPeerOfflineForeground
	}
	return cs.rosterPeerOnlineForeground
}

func createTooltipForPeer(i *mucRosterItem) string {
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

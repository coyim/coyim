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
	g      glibi.Glib
	window gtki.Window

	accountManager *mucAccountManager
	roster         *mucRoster
	roomsServer    *mucRoomsFakeServer

	panel       gtki.Box    `gtk-widget:"panel"`
	panelToggle gtki.Button `gtk-widget:"panel-toggle"`
	main        gtki.Box    `gtk-widget:"main"`
	room        gtki.Box    `gtk-widget:"room"`

	roomPanelOpen  bool
	roomViewActive bool

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

	builder := newBuilder("mainWindow")

	m := &mucUI{
		roomPanelOpen: false,
		builder:       builder,
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
	win := m.builder.get("mucWindow").(gtki.Window)
	win.Show()
	m.window = win
}

func (m *mucUI) togglePanel() {
	isOpen := !m.roomPanelOpen

	var toggleLabel string
	if isOpen {
		toggleLabel = "Hide panel"
	} else {
		toggleLabel = "Show panel"
	}
	m.panelToggle.SetProperty("label", toggleLabel)
	m.panel.SetVisible(isOpen)
	m.roomPanelOpen = isOpen
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

	builder := newBuilder("mainWindow")

	panicOnDevError(builder.bindObjects(m))

	builder.ConnectSignals(map[string]interface{}{
		"on_activate_buddy": m.onActivateRosterRow,
		"on_button_press":   m.onButtonPress,
		"on_toggle_panel":   m.togglePanel,
		"on_close_panel":    m.closeRoomWindow,
	})
}

func (m *mucUI) initDemoAccounts() {
	m.accountManager = &mucAccountManager{}

	accounts := []*mucAccount{
		&mucAccount{
			mucRosterItem: &mucRosterItem{
				id:     "sandy@autonomia.digital",
				status: mucStatusOnline,
			},
			contacts: []*mucAccountContact{
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "pedro@autonomia.digital",
						name:   "Pedro Enrique",
						status: mucStatusOnline,
					},
				},
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "rafael@autonomia.digital",
						status: mucStatusOnline,
					},
				},
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "cristina@autonomia.digital",
						name:   "Cristina Salcedo",
						status: mucStatusOffline,
					},
				},
			},
			rooms: []string{
				"#coyim:matrix:autonomia.digital",
				"#wahay:matrix:autonomia.digital",
			},
		},
		&mucAccount{
			mucRosterItem: &mucRosterItem{
				id:     "pedro@autonomia.digital",
				status: mucStatusOnline,
			},
			contacts: []*mucAccountContact{
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "sandy@autonomia.digital",
						name:   "Sandy Acurio",
						status: mucStatusOnline,
					},
				},
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "rafael@autonomia.digital",
						status: mucStatusOnline,
					},
				},
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "cristina@autonomia.digital",
						name:   "Cristina Salcedo",
						status: mucStatusOffline,
					},
				},
			},
			rooms: []string{
				"#main:matrix:autonomia.digital",
			},
		},
		&mucAccount{
			mucRosterItem: &mucRosterItem{
				id:     "pedro@coy.im",
				name:   "Pedro CoyIM",
				status: mucStatusOffline,
			},
		},
	}

	for _, a := range accounts {
		m.accountManager.addAccount(a)
	}
}

func (m *mucUI) closeRoomWindow() {
	if !m.roomViewActive {
		return
	}

	m.main.SetHExpand(true)
	m.room.SetVisible(false)
	m.roomViewActive = false
}

func (m *mucUI) onButtonPress(view gtki.TreeView, ev gdki.Event) bool {
	return false
}

func (m *mucUI) initRooms() {
	s := &mucRoomsFakeServer{
		rooms: map[string]*mucRoom{},
	}

	rooms := map[string]*mucRoom{
		"#coyim:matrix:autonomia.digital": &mucRoom{
			name: "CoyIM",
		},
		"#wahay:matrix:autonomia.digital": &mucRoom{
			name: "Wahay",
		},
		"#main:matrix:autonomia.digital": &mucRoom{
			name: "Main",
		},
	}

	for id, r := range rooms {
		s.addRoom(id, r)
	}

	m.roomsServer = s
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

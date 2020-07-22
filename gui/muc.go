package gui

import (
	"fmt"
	"html"
	"log"

	"github.com/coyim/gotk3adapter/gtki"
)

type mucUI struct {
	u              *gtkUI
	accountManager *mucAccountManager
	roster         *mucRoster
	roomsServer    *mucRoomsFakeServer
	builder        *builder

	panel       gtki.Box    `gtk-widget:"panel"`
	panelToggle gtki.Button `gtk-widget:"panel-toggle"`

	panelOpen bool
}

func (m *mucUI) togglePanel() {
	var toggleLabel string
	if m.panelOpen {
		toggleLabel = "Hide panel"
	} else {
		toggleLabel = "Show panel"
	}
	m.panel.SetVisible(m.panelOpen)
	m.panelToggle.SetProperty("label", toggleLabel)
	m.panelOpen = !m.panelOpen
}

type mucPeerStatus string

var (
	mucStatusConnecting mucPeerStatus = "connecting"
	mucStatusOnline     mucPeerStatus = "online"
	mucStatusOffline    mucPeerStatus = "offline"
)

type mucAccountContact struct {
	*mucRosterItem
}

type mucAccount struct {
	*mucRosterItem
	contacts []*mucAccountContact
	rooms    []string
}

type mucAccountManager struct {
	accounts []*mucAccount
}

type mucRoster struct {
	widget gtki.ScrolledWindow `gtk-widget:"roster"`
	model  gtki.ListStore      `gtk-widget:"roster-model"`
	view   gtki.TreeView       `gtk-widget:"roster-tree"`

	u     *gtkUI
	rooms *mucRoomsFakeServer
}

type mucRoomsFakeServer struct {
	rooms map[string]*mucRoom
}

type mucRoom struct {
	id      string
	name    string
	status  mucPeerStatus
	members *mucMembers
}

type mucMembers struct {
	widget gtki.ScrolledWindow `gtk-widget:"room-members"`
	model  gtki.ListStore      `gtk-widget:"room-members-model"`
	view   gtki.TreeView       `gtk-widget:"room-members-tree"`

	u *gtkUI
}

type mucRosterItem struct {
	id     string
	name   string
	status mucPeerStatus
}

func (u *gtkUI) initMUCMockups() {
	builder := newBuilder("MUC")

	m := &mucUI{
		u:       u,
		builder: builder,
	}

	m.init()

	m.addAccountsToRoster()

	m.showWindow()
}

func (m *mucUI) init() {
	m.initRooms()

	m.initRoster()

	m.initDemoAccounts()

	panicOnDevError(m.builder.bindObjects(m))

	m.builder.ConnectSignals(map[string]interface{}{
		"on_toggle_panel": m.togglePanel,
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

func (m *mucUI) initRoster() {
	r := &mucRoster{
		u:     m.u,
		rooms: m.roomsServer,
	}

	panicOnDevError(m.builder.bindObjects(r))

	m.roster = r
}

func (m *mucUI) addAccountsToRoster() {
	for _, a := range m.accountManager.accounts {
		m.roster.add(a)
	}
}

func (m *mucUI) showWindow() {
	win := m.builder.get("mucWindow").(gtki.Window)
	win.Show()
}

func (m *mucAccountManager) addAccount(a *mucAccount) {
	m.accounts = append(m.accounts, a)
}

func (r *mucRoster) add(account *mucAccount) {
	r.addAccount(account)
	r.view.ExpandAll()
}

func (r *mucRoster) addAccount(account *mucAccount) {
	cs := r.u.currentColorSet()
	parentIter := r.model.Append()

	accountCounter := &counter{}

	// Contacts for this account
	r.addAccountContacts(account.contacts, accountCounter)

	// Rooms this contact is suscribed to
	r.addAccountRooms(account.rooms)

	displayName := account.displayName()

	_ = r.model.SetValue(parentIter, indexParentJid, displayName)
	_ = r.model.SetValue(parentIter, indexAccountID, account.id)
	_ = r.model.SetValue(parentIter, indexRowType, "account")
	_ = r.model.SetValue(parentIter, indexWeight, 700)

	bgcolor := cs.rosterAccountOnlineBackground
	if account.isOffline() {
		bgcolor = cs.rosterAccountOfflineBackground
	}
	_ = r.model.SetValue(parentIter, indexBackgroundColor, bgcolor)

	_ = r.model.SetValue(parentIter, indexStatusIcon, statusIcons[account.getStatus()].GetPixbuf())
	_ = r.model.SetValue(parentIter, indexParentDisplayName, createGroupDisplayName(displayName, accountCounter, true))
}

func (r *mucRoster) addAccountContacts(contacts []*mucAccountContact, accountCounter *counter) {
	groupCounter := &counter{}

	for _, item := range contacts {
		o := item.isOnline()
		accountCounter.inc(true, o)
		groupCounter.inc(true, o)
		r.addItem(item.mucRosterItem, "")
	}
}

func (r *mucRoster) addAccountRooms(rooms []string) {
	for _, id := range rooms {
		room, err := r.rooms.byID(id)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		r.addRoom(id, room)
	}
}

func (r *mucRoster) addRoom(id string, room *mucRoom) {
	r.addItem(&mucRosterItem{
		id:     id,
		name:   room.name,
		status: room.status,
	}, "#")
}

func (r *mucRoster) addItem(item *mucRosterItem, indent string) {
	cs := r.u.currentColorSet()
	iter := r.model.Append()

	setValues(
		r.model,
		iter,
		item.id,
		fmt.Sprintf("%s%s", indent, item.displayName()),
		"BelongsTo",
		decideColorForPeer(cs, item),
		cs.rosterPeerBackground,
		nil,
		createTooltipForPeer(item),
	)
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

func (r *mucRoomsFakeServer) addRoom(id string, room *mucRoom) {
	r.rooms[id] = room
}

func (r *mucRoomsFakeServer) byID(id string) (*mucRoom, error) {
	if room, ok := r.rooms[id]; ok {
		return room, nil
	}
	return nil, fmt.Errorf("roomt %s not found", id)
}

func (a *mucAccount) isOffline() bool {
	return false
}

func (a *mucAccount) displayName() string {
	if a.name != "" {
		return a.name
	}

	return a.id
}

func setValues(v gtki.ListStore, iter gtki.TreeIter, values ...interface{}) {
	for i, val := range values {
		if val != nil {
			_ = v.SetValue(iter, i, val)
		}
	}
}

func (i *mucRosterItem) displayName() string {
	if i.name != "" {
		return i.name
	}

	return i.id
}

func (i *mucRosterItem) isOnline() bool {
	return i.status == mucStatusOnline
}

func (i *mucRosterItem) getStatus() string {
	if i.status == mucStatusConnecting {
		return "connecting"
	}

	if i.status == mucStatusOffline {
		return "offline"
	}

	return "available"
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

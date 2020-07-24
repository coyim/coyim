package muc

import (
	"fmt"
	"log"

	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoster struct {
	widget gtki.ScrolledWindow `gtk-widget:"roster"`
	model  gtki.ListStore      `gtk-widget:"roster-model"`
	view   gtki.TreeView       `gtk-widget:"roster-tree"`

	rooms *mucRoomsFakeServer

	u *mucUI
}

func (m *mucUI) initRoster() {
	r := &mucRoster{
		rooms: m.roomsServer,
		u:     m,
	}

	panicOnDevError(m.builder.bindObjects(r))

	m.roster = r
}

func (m *mucUI) addAccountsToRoster() {
	for _, a := range m.accountManager.accounts {
		m.roster.add(a)
	}
}

func (m *mucUI) onActivateRosterRow(v gtki.TreeView, path gtki.TreePath) {
	iter, err := m.roster.model.GetIter(path)
	if err != nil {
		return
	}

	peer := getFromModelIterMUC(m.roster.model, iter, indexJid)
	rowType := getFromModelIterMUC(m.roster.model, iter, indexRowType)

	switch rowType {
	case "room":
		m.openRoomView(peer)
	case "group":
		// We ignore this, since a double click on the group doesn't really have any effect
	default:
		panic(fmt.Sprintf("unknown roster row type: %s", rowType))
	}
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
		r.addItem(item.mucRosterItem, "peer", "")
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
	roomItem := &mucRosterItem{
		id:     id,
		name:   room.name,
		status: room.status,
	}

	r.addItem(roomItem, "room", "#")
}

func (r *mucRoster) addItem(item *mucRosterItem, rowType string, indent string) {
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

	_ = r.model.SetValue(iter, indexRowType, rowType)
}

type mucRosterItem struct {
	id     string
	name   string
	status mucPeerStatus
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

func (c *counter) inc(total, online bool) {
	if total {
		c.total++
	}
	if online {
		c.online++
	}
}

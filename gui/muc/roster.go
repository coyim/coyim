package muc

import (
	"fmt"
	"log"

	"github.com/coyim/gotk3adapter/gtki"
)

type roster struct {
	widget gtki.ScrolledWindow `gtk-widget:"roster"`
	model  gtki.ListStore      `gtk-widget:"roster-model"`
	view   gtki.TreeView       `gtk-widget:"roster-tree"`

	rooms *roomsFakeServer

	u *gtkUI
}

func (u *gtkUI) initRoster() {
	r := &roster{
		rooms: u.roomsServer,
		u:     u,
	}

	panicOnDevError(u.builder.bindObjects(r))

	u.roster = r
}

func (u *gtkUI) addAccountsToRoster() {
	for _, a := range u.accountManager.accounts {
		u.roster.add(a)
	}
}

func (u *gtkUI) onActivateRosterRow(v gtki.TreeView, path gtki.TreePath) {
	iter, err := u.roster.model.GetIter(path)
	if err != nil {
		return
	}

	peer := getFromModelIterMUC(u.roster.model, iter, indexJid)
	rowType := getFromModelIterMUC(u.roster.model, iter, indexRowType)

	switch rowType {
	case "room":
		u.openRoomView(peer)
	case "group":
		// We ignore this, since a double click on the group doesn't really have any effect
	default:
		panic(fmt.Sprintf("unknown roster row type: %s", rowType))
	}
}

func (r *roster) add(a *account) {
	r.addAccount(a)
	r.view.ExpandAll()
}

func (r *roster) addAccount(a *account) {
	cs := r.u.currentColorSet()
	parentIter := r.model.Append()

	accountCounter := &counter{}

	// Contacts for this account
	r.addAccountContacts(a.contacts, accountCounter)

	// Rooms this contact is suscribed to
	r.addAccountRooms(a.rooms)

	displayName := a.displayName()

	_ = r.model.SetValue(parentIter, indexParentJid, displayName)
	_ = r.model.SetValue(parentIter, indexAccountID, a.id)
	_ = r.model.SetValue(parentIter, indexRowType, "account")
	_ = r.model.SetValue(parentIter, indexWeight, 700)

	bgcolor := cs.rosterAccountOnlineBackground
	if a.isOffline() {
		bgcolor = cs.rosterAccountOfflineBackground
	}
	_ = r.model.SetValue(parentIter, indexBackgroundColor, bgcolor)

	_ = r.model.SetValue(parentIter, indexStatusIcon, statusIcons[a.getStatus()].GetPixbuf())
	_ = r.model.SetValue(parentIter, indexParentDisplayName, createGroupDisplayName(displayName, accountCounter, true))
}

func (r *roster) addAccountContacts(contacts []*contact, accountCounter *counter) {
	groupCounter := &counter{}

	for _, item := range contacts {
		o := item.isOnline()
		accountCounter.inc(true, o)
		groupCounter.inc(true, o)
		r.addItem(item.rosterItem, "peer", "")
	}
}

func (r *roster) addAccountRooms(rooms []string) {
	for _, id := range rooms {
		room, err := r.rooms.byID(id)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		r.addRoom(id, room)
	}
}

func (r *roster) addRoom(id string, r2 *room) {
	roomItem := &rosterItem{
		id:     id,
		name:   r2.name,
		status: r2.status,
	}

	r.addItem(roomItem, "room", "#")
}

func (r *roster) addItem(item *rosterItem, rowType string, indent string) {
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

type rosterItem struct {
	id     string
	name   string
	status peerStatus
}

func (i *rosterItem) displayName() string {
	if i.name != "" {
		return i.name
	}

	return i.id
}

func (i *rosterItem) isOnline() bool {
	return i.status == statusOnline
}

func (i *rosterItem) isOffline() bool {
	return i.status == statusOffline
}

func (i *rosterItem) getStatus() string {
	if i.status == statusConnecting {
		return "connecting"
	}

	if i.status == statusOffline {
		return "offline"
	}

	return "available"
}

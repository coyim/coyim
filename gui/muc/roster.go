package muc

import (
	"fmt"
	"log"

	"github.com/coyim/gotk3adapter/gtki"
)

type roster struct {
	widget gtki.ScrolledWindow `gtk-widget:"roster"`
	model  gtki.TreeStore      `gtk-widget:"roster-model"`
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
	r.addAccount(a, nil)
	r.view.ExpandAll()
}

func (r *roster) addAccount(a *account, p gtki.TreeIter) {
	cs := r.u.currentColorSet()
	parentIter := r.model.Append(p)

	accountCounter := &counter{}

	// Contacts for this account
	r.addAccountContacts(a.contacts, accountCounter, "", parentIter)

	// Rooms this contact is suscribed to
	r.addAccountRooms(a.rooms, identRoom, parentIter)

	// Groups for this account
	r.addAccountGroups(a.groups, parentIter)

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

func (r *roster) addAccountContacts(contacts []*contact, accountCounter *counter, ident string, p gtki.TreeIter) {
	groupCounter := &counter{}

	for _, item := range contacts {
		o := item.isOnline()
		accountCounter.inc(true, o)
		groupCounter.inc(true, o)
		r.addItem(item.rosterItem, "peer", ident, p)
	}
}

func (r *roster) addAccountRooms(rooms []*accountRoom, ident string, p gtki.TreeIter) {
	for _, r2 := range rooms {
		room, err := r.rooms.byID(r2.id)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		r.addRoom(r2, room, ident, p)
	}
}

func (r *roster) addAccountGroups(groups []*group, p gtki.TreeIter) {
	for _, g := range groups {
		r.addGroup(g, p)
	}
}

func (r *roster) addRoom(r1 *accountRoom, r2 *room, ident string, p gtki.TreeIter) {
	roomItem := &rosterItem{
		id:     r1.id,
		name:   r2.name,
		status: r1.status,
	}

	r.addItem(roomItem, "room", ident, p)
}

func (r *roster) addGroup(g *group, p gtki.TreeIter) {
	cs := r.u.currentColorSet()
	iter := r.model.Append(p)

	setValues(
		r.model,
		iter,
		g.id,
		g.displayName(),
		"BelongsTo",
		decideColorForPeer(cs, g.rosterItem),
		cs.rosterPeerBackground,
		400,
		createTooltipForPeer(g.rosterItem),
	)

	_ = r.model.SetValue(iter, indexRowType, "group")

	r.addAccountContacts(g.contacts, &counter{}, "", iter)
	r.addAccountRooms(g.rooms, "#", iter)
}

func (r *roster) addItem(item *rosterItem, rowType string, indent string, p gtki.TreeIter) {
	cs := r.u.currentColorSet()
	iter := r.model.Append(p)

	setValues(
		r.model,
		iter,
		item.id,
		fmt.Sprintf("%s%s", indent, item.displayName()),
		"BelongsTo",
		decideColorForPeer(cs, item),
		cs.rosterPeerBackground,
		400,
		createTooltipForPeer(item),
	)

	_ = r.model.SetValue(iter, indexRowType, rowType)

	if item.hasStatus() {
		_ = r.model.SetValue(iter, indexStatusIcon, statusIcons[decideStatusFor(item)].GetPixbuf())
	}
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

func (i *rosterItem) hasStatus() bool {
	return true
}

func decideStatusFor(r *rosterItem) string {
	if !r.isOnline() {
		return "offline"
	}

	return "available"
}

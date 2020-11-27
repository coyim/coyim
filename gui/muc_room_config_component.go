package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigComponent struct {
	u *gtkUI

	notebook            gtki.Notebook `gtk-widget:"config-room-notebook"`
	bannedList          gtki.TreeView `gtk-widget:"config-banned-list"`
	membersList         gtki.TreeView `gtk-widget:"config-members-list"`
	ownersList          gtki.TreeView `gtk-widget:"config-owners-list"`
	adminList           gtki.TreeView `gtk-widget:"config-admin-list"`
	bannedAddButton     gtki.Button   `gtk-widget:"ban-add"`
	membersAddButton    gtki.Button   `gtk-widget:"member-add"`
	ownersAddButton     gtki.Button   `gtk-widget:"owner-add"`
	adminAddButton      gtki.Button   `gtk-widget:"admin-add"`
	bannedRemoveButton  gtki.Button   `gtk-widget:"ban-remove"`
	membersRemoveButton gtki.Button   `gtk-widget:"member-remove"`
	ownersRemoveButton  gtki.Button   `gtk-widget:"owner-remove"`
	adminRemoveButton   gtki.Button   `gtk-widget:"admin-remove"`

	bannedListComponent  *mucRoomConfigListComponent
	membersListComponent *mucRoomConfigListComponent
	ownersListComponent  *mucRoomConfigListComponent
	adminListComponent   *mucRoomConfigListComponent

	bannedListAddComponent  *mucRoomConfigListAddComponent
	membersListAddComponent *mucRoomConfigListAddComponent
	ownersListAddComponent  *mucRoomConfigListAddComponent
	adminListAddComponent   *mucRoomConfigListAddComponent
}

func (u *gtkUI) newMUCRoomConfigComponent() *mucRoomConfigComponent {
	c := &mucRoomConfigComponent{u: u}

	c.initBuilder()
	c.initMembersListComponent()

	return c
}

func (c *mucRoomConfigComponent) initBuilder() {
	b := newBuilder("MUCRoomConfig")
	panicOnDevError(b.bindObjects(c))
}

func (c *mucRoomConfigComponent) initMembersListComponent() {
	membersListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
		// nickname
		glibi.TYPE_STRING,
		// role
		glibi.TYPE_STRING,
	}

	c.membersListComponent = c.u.newMUCRoomConfigListComponent(
		c.membersList,
		membersListColumns,
		c.membersAddButton,
		c.membersRemoveButton,
		c.onAddMembersToList,
	)

	c.membersListAddComponent = c.u.newMUCRoomConfigListAddComponent(
		i18n.Local("Adding member..."),
		i18n.Local("Whom do you want to make a member?"),
		newMUCRoomConfigListMembersForm(nil),
		c.membersListComponent.addListItem,
	)
}

func (c *mucRoomConfigComponent) onAddMembersToList() {
	c.membersListAddComponent.show()
}

func (c *mucRoomConfigComponent) configurationView() gtki.Widget {
	return c.notebook
}

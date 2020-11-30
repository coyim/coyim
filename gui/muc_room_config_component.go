package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigComponent struct {
	u *gtkUI

	notebook            gtki.Notebook `gtk-widget:"config-room-notebook"`
	kickedList          gtki.TreeView `gtk-widget:"config-banned-list"`
	membersList         gtki.TreeView `gtk-widget:"config-members-list"`
	ownersList          gtki.TreeView `gtk-widget:"config-owners-list"`
	adminList           gtki.TreeView `gtk-widget:"config-admin-list"`
	kickedAddButton     gtki.Button   `gtk-widget:"ban-add"`
	membersAddButton    gtki.Button   `gtk-widget:"member-add"`
	ownersAddButton     gtki.Button   `gtk-widget:"owner-add"`
	adminAddButton      gtki.Button   `gtk-widget:"admin-add"`
	kickedRemoveButton  gtki.Button   `gtk-widget:"ban-remove"`
	membersRemoveButton gtki.Button   `gtk-widget:"member-remove"`
	ownersRemoveButton  gtki.Button   `gtk-widget:"owner-remove"`
	adminRemoveButton   gtki.Button   `gtk-widget:"admin-remove"`

	kickedListController  *mucRoomConfigListController
	membersListController *mucRoomConfigListController
	ownersListController  *mucRoomConfigListController
	adminsListController  *mucRoomConfigListController
}

func (u *gtkUI) newMUCRoomConfigComponent() *mucRoomConfigComponent {
	c := &mucRoomConfigComponent{u: u}

	c.initBuilder()
	c.initKickedListController()
	c.initMembersListController()
	c.initOwnersListController()
	c.initAdminsListController()

	return c
}

func (c *mucRoomConfigComponent) initBuilder() {
	b := newBuilder("MUCRoomConfig")
	panicOnDevError(b.bindObjects(c))
}

func (c *mucRoomConfigComponent) initKickedListController() {
	kickedListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
		// reason
		glibi.TYPE_STRING,
	}

	c.kickedListController = c.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        c.kickedAddButton,
		removeOccupantButton:     c.kickedRemoveButton,
		occupantsTreeView:        c.kickedList,
		occupantsTreeViewColumns: kickedListColumns,
		addOccupantDialogTitle:   i18n.Local("Adding banned member..."),
		addOccupantDescription:   i18n.Local("Whom do you want to ban a member?"),
		addOccupantForm:          newMUCRoomConfigListKickedForm,
	})
}

func (c *mucRoomConfigComponent) initMembersListController() {
	membersListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
		// nickname
		glibi.TYPE_STRING,
		// role
		glibi.TYPE_STRING,
	}

	c.membersListController = c.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        c.membersAddButton,
		removeOccupantButton:     c.membersRemoveButton,
		occupantsTreeView:        c.membersList,
		occupantsTreeViewColumns: membersListColumns,
		addOccupantDialogTitle:   i18n.Local("Adding member..."),
		addOccupantDescription:   i18n.Local("Whom do you want to make a member?"),
		addOccupantForm:          newMUCRoomConfigListMembersForm,
	})
}

func (c *mucRoomConfigComponent) initOwnersListController() {
	ownersListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
	}

	c.ownersListController = c.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        c.ownersAddButton,
		removeOccupantButton:     c.ownersRemoveButton,
		occupantsTreeView:        c.ownersList,
		occupantsTreeViewColumns: ownersListColumns,
		addOccupantDialogTitle:   i18n.Local("Adding owner member..."),
		addOccupantDescription:   i18n.Local("Whom do you want to make an owner member?"),
		addOccupantForm:          newMUCRoomConfigListOwnersForm,
	})
}

func (c *mucRoomConfigComponent) initAdminsListController() {
	adminsListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
	}

	c.adminsListController = c.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        c.adminAddButton,
		removeOccupantButton:     c.adminRemoveButton,
		occupantsTreeView:        c.adminList,
		occupantsTreeViewColumns: adminsListColumns,
		addOccupantDialogTitle:   i18n.Local("Adding admin member..."),
		addOccupantDescription:   i18n.Local("Whom do you want to make an admin member?"),
		addOccupantForm:          newMUCRoomConfigListAdminsForm,
	})
}

func (c *mucRoomConfigComponent) configurationView() gtki.Widget {
	return c.notebook
}

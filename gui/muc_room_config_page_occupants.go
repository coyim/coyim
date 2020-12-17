package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigOccupantsPage struct {
	*roomConfigPageBase

	configOccupantsBox gtki.Box `gtk-widget:"room-config-occupants-page"`

	// Ban list-related UI elements
	banList           gtki.TreeView `gtk-widget:"room-config-ban-list"`
	banAddButton      gtki.Button   `gtk-widget:"room-ban-add"`
	banRemoveButton   gtki.Button   `gtk-widget:"room-ban-remove"`
	banListController *mucRoomConfigListController

	// Members list-related UI elements
	membersList           gtki.TreeView `gtk-widget:"room-config-members-list"`
	membersAddButton      gtki.Button   `gtk-widget:"room-member-add"`
	membersRemoveButton   gtki.Button   `gtk-widget:"room-member-remove"`
	membersListController *mucRoomConfigListController

	// Owners list-related UI elements
	ownersList           gtki.TreeView `gtk-widget:"room-config-owners-list"`
	ownersAddButton      gtki.Button   `gtk-widget:"room-owner-add"`
	ownersRemoveButton   gtki.Button   `gtk-widget:"room-owner-remove"`
	ownersListController *mucRoomConfigListController

	// Admin list-related UI elements
	adminList            gtki.TreeView `gtk-widget:"room-config-admin-list"`
	adminAddButton       gtki.Button   `gtk-widget:"room-admin-add"`
	adminRemoveButton    gtki.Button   `gtk-widget:"room-admin-remove"`
	adminsListController *mucRoomConfigListController
}

func (c *mucRoomConfigComponent) newRoomConfigOccupantsPage() mucRoomConfigPage {
	p := &roomConfigOccupantsPage{}

	builder := newBuilder("MUCRoomConfigPageOccupants")
	panicOnDevError(builder.bindObjects(p))

	p.roomConfigPageBase = c.newConfigPage(p.configOccupantsBox)

	p.initOccupantsLists()
	p.initDefaultValues()

	return p
}

func (p *roomConfigOccupantsPage) initOccupantsLists() {
	p.initBanListController()
	p.initMembersListController()
	p.initOwnersListController()
	p.initAdminsListController()
}

func (p *roomConfigOccupantsPage) initBanListController() {
	banListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
		// reason
		glibi.TYPE_STRING,
	}

	p.banListController = p.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        p.banAddButton,
		removeOccupantButton:     p.banRemoveButton,
		occupantsTreeView:        p.banList,
		occupantsTreeViewColumns: banListColumns,
		addOccupantDialogTitle:   i18n.Local("Add a blocked (banned) occupant"),
		addOccupantDescription:   i18n.Local("Please specify the information of the room occupant do you want to block or ban"),
		addOccupantForm:          newMUCRoomConfigListKickedForm,
	})
}

func (p *roomConfigOccupantsPage) initMembersListController() {
	membersListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
		// nickname
		glibi.TYPE_STRING,
		// role
		glibi.TYPE_STRING,
	}

	p.membersListController = p.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        p.membersAddButton,
		removeOccupantButton:     p.membersRemoveButton,
		occupantsTreeView:        p.membersList,
		occupantsTreeViewColumns: membersListColumns,
		addOccupantDialogTitle:   i18n.Local("Add a room member"),
		addOccupantDescription:   i18n.Local("Please specify the information of the room occupant you want to add as a permanent member"),
		addOccupantForm:          newMUCRoomConfigListMembersForm,
	})
}

func (p *roomConfigOccupantsPage) initOwnersListController() {
	ownersListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
	}

	p.ownersListController = p.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        p.ownersAddButton,
		removeOccupantButton:     p.ownersRemoveButton,
		occupantsTreeView:        p.ownersList,
		occupantsTreeViewColumns: ownersListColumns,
		addOccupantDialogTitle:   i18n.Local("Add a room owner"),
		addOccupantDescription:   i18n.Local("Please specify the information of the occupant you want to add as a room owner"),
		addOccupantForm:          newMUCRoomConfigListOwnersForm,
	})
}

func (p *roomConfigOccupantsPage) initAdminsListController() {
	adminsListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
	}

	p.adminsListController = p.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        p.adminAddButton,
		removeOccupantButton:     p.adminRemoveButton,
		occupantsTreeView:        p.adminList,
		occupantsTreeViewColumns: adminsListColumns,
		addOccupantDialogTitle:   i18n.Local("Add room administrator"),
		addOccupantDescription:   i18n.Local("Please specify the information of the occupant you want to add as a room administrator"),
		addOccupantForm:          newMUCRoomConfigListAdminsForm,
	})
}

func (p *roomConfigOccupantsPage) initDefaultValues() {
	// TODO
}

func (p *roomConfigOccupantsPage) collectData() {
	// TODO
}

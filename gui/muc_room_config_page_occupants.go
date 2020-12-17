package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigOccupantsPage struct {
	*roomConfigPageBase

	configOccupantsBox gtki.Box      `gtk-widget:"room-config-occupants-page"`
	ownersList         gtki.TreeView `gtk-widget:"room-config-owners-list"`
	ownersAddButton    gtki.Button   `gtk-widget:"room-owner-add"`
	ownersRemoveButton gtki.Button   `gtk-widget:"room-owner-remove"`
	adminList          gtki.TreeView `gtk-widget:"room-config-admin-list"`
	adminAddButton     gtki.Button   `gtk-widget:"room-admin-add"`
	adminRemoveButton  gtki.Button   `gtk-widget:"room-admin-remove"`

	ownersListController *mucRoomConfigListController
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
	p.initOwnersListController()
	p.initAdminsListController()
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
	p.form.Owners = jidListFromConfigListController(p.ownersListController)
	p.form.Admins = jidListFromConfigListController(p.adminsListController)
}

func jidListFromConfigListController(l *mucRoomConfigListController) (result []jid.Any) {
	for _, i := range l.listItems() {
		result = append(result, jid.Parse(i[0]))
	}
	return result
}

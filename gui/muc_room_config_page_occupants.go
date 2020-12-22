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
	notificationBox    gtki.Box      `gtk-widget:"notification-box"`
	ownersList         gtki.TreeView `gtk-widget:"room-config-owners-list"`
	ownersAddButton    gtki.Button   `gtk-widget:"room-owner-add"`
	ownersRemoveButton gtki.Button   `gtk-widget:"room-owner-remove"`
	adminList          gtki.TreeView `gtk-widget:"room-config-admin-list"`
	adminAddButton     gtki.Button   `gtk-widget:"room-admin-add"`
	adminRemoveButton  gtki.Button   `gtk-widget:"room-admin-remove"`

	ownersListController *mucRoomConfigListController
	adminsListController *mucRoomConfigListController
}

func (c *mucRoomConfigComponent) newRoomConfigOccupantsPage(parent gtki.Window) mucRoomConfigPage {
	p := &roomConfigOccupantsPage{}

	builder := newBuilder("MUCRoomConfigPageOccupants")
	panicOnDevError(builder.bindObjects(p))

	p.roomConfigPageBase = c.newConfigPage(p.configOccupantsBox, p.notificationBox)

	p.initOccupantsLists(parent)
	p.initDefaultValues()

	return p
}

func (p *roomConfigOccupantsPage) initOccupantsLists(parent gtki.Window) {
	p.initOwnersListController(parent)
	p.initAdminsListController(parent)
}

func (p *roomConfigOccupantsPage) initOwnersListController(parent gtki.Window) {
	ownersListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
	}

	p.ownersListController = p.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        p.ownersAddButton,
		removeOccupantButton:     p.ownersRemoveButton,
		occupantsTreeView:        p.ownersList,
		occupantsTreeViewColumns: ownersListColumns,
		parentWindow:             parent,
		addOccupantDialogTitle:   i18n.Local("Add a room owner"),
		addOccupantDescription:   i18n.Local("Please specify the information of the occupant you want to add as a room owner"),
		addOccupantForm:          newMUCRoomConfigListOwnersForm,
	})
}

func (p *roomConfigOccupantsPage) initAdminsListController(parent gtki.Window) {
	adminsListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
	}

	p.adminsListController = p.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        p.adminAddButton,
		removeOccupantButton:     p.adminRemoveButton,
		occupantsTreeView:        p.adminList,
		occupantsTreeViewColumns: adminsListColumns,
		parentWindow:             parent,
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

package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

const (
	ownersListColumnJidIndex = 0
	adminsListColumnJidIndex = 0
)

type roomConfigOccupantsPage struct {
	*roomConfigPageBase

	ownersList            gtki.TreeView         `gtk-widget:"room-config-owners-list"`
	ownersAddButton       gtki.Button           `gtk-widget:"room-owner-add"`
	ownersRemoveButton    gtki.Button           `gtk-widget:"room-owner-remove"`
	ownersListJidRenderer gtki.CellRendererText `gtk-widget:"owner-jid-text-renderer"`
	adminList             gtki.TreeView         `gtk-widget:"room-config-admin-list"`
	adminAddButton        gtki.Button           `gtk-widget:"room-admin-add"`
	adminRemoveButton     gtki.Button           `gtk-widget:"room-admin-remove"`
	adminListJidRenderer  gtki.CellRendererText `gtk-widget:"admin-jid-text-renderer"`

	ownersListController *mucRoomConfigListController
	adminsListController *mucRoomConfigListController
}

func (c *mucRoomConfigComponent) newRoomConfigOccupantsPage(parent gtki.Window) mucRoomConfigPage {
	p := &roomConfigOccupantsPage{}
	p.roomConfigPageBase = c.newConfigPage("occupants", "MUCRoomConfigPageOccupants", p, map[string]interface{}{
		"on_owner_jid_edited": p.onOwnerJidEdited,
		"on_admin_jid_edited": p.onAdminJidEdited,
	})

	p.initOccupantsEditableCells()
	p.initOccupantsLists(parent)
	p.initDefaultValues()

	return p
}

func (p *roomConfigOccupantsPage) onOwnerJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	p.updateOccupantListCellForString("owners", p.ownersListController, ownersListColumnJidIndex, path, newValue)
}

func (p *roomConfigOccupantsPage) onAdminJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	p.updateOccupantListCellForString("admins", p.adminsListController, adminsListColumnJidIndex, path, newValue)
}

func (p *roomConfigOccupantsPage) updateOccupantListCellForString(list string, controller *mucRoomConfigListController, column int, path string, newValue string) {
	err := controller.updateCellForString(column, path, newValue)
	if err != nil {
		p.log.WithError(err).WithFields(log.Fields{
			"path":    path,
			"newText": newValue,
			"list":    list,
		}).Error("The occupant's jid can't be updated")
	}
}

func (p *roomConfigOccupantsPage) initOccupantsEditableCells() {
	p.ownersListJidRenderer.SetProperty("editable", true)
	p.adminListJidRenderer.SetProperty("editable", true)
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
		addOccupantDialogTitle:   i18n.Local("Add owners"),
		addOccupantDescription:   i18n.Local("Here you can add one or more new owners to the room. You will have to use the account address of the user in order to make them an owner. This address can either be a simple one, such as user@example.org or a full one, such as user@example.org/abcdef."),
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
		addOccupantDialogTitle:   i18n.Local("Add administrators"),
		addOccupantDescription:   i18n.Local("Here you can add one or more new administrators to the room. You will have to use the account address of the user in order to make them an administrator. This address can either be a simple one, such as user@example.org or a full one, such as user@example.org/abcdef."),
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

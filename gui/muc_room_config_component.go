package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomConfigInformationPage = iota
	roomConfigAccessPage
	roomConfigPermissionsPage
	roomConfigOccupantsPage
	roomConfigOthersPage
	roomConfigSummaryPage
)

type mucRoomConfigComponent struct {
	u    *gtkUI
	form *muc.RoomConfigForm

	// Information cofiguration page
	configInformationPage *mucRoomConfigPage
	configInformationBox  gtki.Box      `gtk-widget:"room-config-info-page"`
	roomTitle             gtki.Entry    `gtk-widget:"room-title"`
	roomDescription       gtki.TextView `gtk-widget:"room-description"`
	roomLanguage          gtki.Entry    `gtk-widget:"room-language"`
	roomPersistent        gtki.Switch   `gtk-widget:"room-persistent"`
	roomPublic            gtki.Switch   `gtk-widget:"room-public"`

	// Access cofiguration page
	configAccessPage *mucRoomConfigPage
	configAccessBox  gtki.Box    `gtk-widget:"room-config-access-page"`
	roomPassword     gtki.Entry  `gtk-widget:"room-password"`
	roomMembersOnly  gtki.Switch `gtk-widget:"room-membersonly"`
	roomAllowInvites gtki.Switch `gtk-widget:"room-allowinvites"`

	// Permissions cofiguration page
	configPermissionsPage *mucRoomConfigPage
	configPermissionsBox  gtki.Box      `gtk-widget:"room-config-permissions-page"`
	roomChangeSubject     gtki.Switch   `gtk-widget:"room-changesubject"`
	roomModerated         gtki.Switch   `gtk-widget:"room-moderated"`
	roomWhois             gtki.ComboBox `gtk-widget:"room-whois"`

	// Members cofiguration page
	configOccupantsPage *mucRoomConfigPage
	configOccupantsBox  gtki.Box `gtk-widget:"room-config-occupants-page"`

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

	// Others cofiguration page
	configOthersPage    *mucRoomConfigPage
	configOthersBox     gtki.Box        `gtk-widget:"room-config-others-page"`
	roomMaxHistoryFetch gtki.SpinButton `gtk-widget:"room-maxhistoryfetch"`
	roomMaxOccupants    gtki.SpinButton `gtk-widget:"room-maxoccupants"`
	roomEnableLoggin    gtki.Switch     `gtk-widget:"room-enablelogging"`

	// Summary cofiguration page
	configSummaryPage *mucRoomConfigPage
	configSummaryBox  gtki.Box `gtk-widget:"room-config-summary-page"`
}

func (u *gtkUI) newMUCRoomConfigComponent(f *muc.RoomConfigForm) *mucRoomConfigComponent {
	c := &mucRoomConfigComponent{u: u, form: f}

	c.initBuilder()
	c.initConfigPages()

	return c
}

func (c *mucRoomConfigComponent) initBuilder() {
	b := newBuilder("MUCRoomConfig")
	panicOnDevError(b.bindObjects(c))
}

func (c *mucRoomConfigComponent) initConfigPages() {
	c.initConfigInformationPage()
	c.initConfigAccessPage()
	c.initConfigPermissionsPage()
	c.initConfigOccupantsPage()
	c.initConfigOthersPage()
	c.initConfigSummaryPage()
}

func (c *mucRoomConfigComponent) initConfigInformationPage() {
	c.configInformationPage = newMUCRoomConfigPage(c.configInformationBox)
}

func (c *mucRoomConfigComponent) initConfigAccessPage() {
	c.configAccessPage = newMUCRoomConfigPage(c.configAccessBox)
}

func (c *mucRoomConfigComponent) initConfigPermissionsPage() {
	c.configPermissionsPage = newMUCRoomConfigPage(c.configPermissionsBox)
}

func (c *mucRoomConfigComponent) initConfigOccupantsPage() {
	c.configOccupantsPage = newMUCRoomConfigPage(c.configOccupantsBox)
	c.initConfigOccupantsLists()
}

func (c *mucRoomConfigComponent) initConfigOthersPage() {
	c.configOthersPage = newMUCRoomConfigPage(c.configOthersBox)
}

func (c *mucRoomConfigComponent) initConfigSummaryPage() {
	c.configSummaryPage = newMUCRoomConfigPage(c.configSummaryBox)
}

func (c *mucRoomConfigComponent) initConfigOccupantsLists() {
	c.initBanListController()
	c.initMembersListController()
	c.initOwnersListController()
	c.initAdminsListController()
}

func (c *mucRoomConfigComponent) initBanListController() {
	banListColumns := []glibi.Type{
		// jid
		glibi.TYPE_STRING,
		// reason
		glibi.TYPE_STRING,
	}

	c.banListController = c.u.newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:        c.banAddButton,
		removeOccupantButton:     c.banRemoveButton,
		occupantsTreeView:        c.banList,
		occupantsTreeViewColumns: banListColumns,
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

func (c *mucRoomConfigComponent) getConfigPage(p string) *mucRoomConfigPage {
	switch p {
	case "information":
		return c.configInformationPage
	case "access":
		return c.configAccessPage
	case "permissions":
		return c.configPermissionsPage
	case "occupants":
		return c.configOccupantsPage
	case "others":
		return c.configOthersPage
	case "summary":
		return c.configSummaryPage
	default:
		return nil
	}
}

type mucRoomConfigPage struct {
	box gtki.Box
}

func newMUCRoomConfigPage(b gtki.Box) *mucRoomConfigPage {
	return &mucRoomConfigPage{b}
}

func (p *mucRoomConfigPage) getPageView() gtki.Box {
	return p.box
}

package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryPage struct {
	*roomConfigPageBase
	autoJoin bool

	basicInformation        gtki.LinkButton     `gtk-widget:"room-config-summary-basic-information"`
	access                  gtki.LinkButton     `gtk-widget:"room-config-summary-access"`
	permissions             gtki.LinkButton     `gtk-widget:"room-config-summary-permissions"`
	occupants               gtki.LinkButton     `gtk-widget:"room-config-summary-occupants"`
	others                  gtki.LinkButton     `gtk-widget:"room-config-summary-others"`
	title                   gtki.Label          `gtk-widget:"room-config-summary-title"`
	descriptionNotAssigned  gtki.Label          `gtk-widget:"room-config-summary-description-not-assigned"`
	descriptionScrollWindow gtki.ScrolledWindow `gtk-widget:"room-config-summary-description-scrolled"`
	description             gtki.TextView       `gtk-widget:"room-config-summary-description"`
	language                gtki.Label          `gtk-widget:"room-config-summary-language"`
	includePublicList       gtki.Label          `gtk-widget:"room-config-summary-public-label"`
	persistent              gtki.Label          `gtk-widget:"room-config-summary-persistent-label"`
	password                gtki.Label          `gtk-widget:"room-config-summary-password-label"`
	allowInviteUsers        gtki.Label          `gtk-widget:"room-config-summary-invite-label"`
	onlyMembers             gtki.Label          `gtk-widget:"room-config-summary-onlymembers-label"`
	allowSetRoomSubject     gtki.Label          `gtk-widget:"room-config-summary-changesubject-label"`
	moderatedRoom           gtki.Label          `gtk-widget:"room-config-summary-moderated-label"`
	whoIs                   gtki.Label          `gtk-widget:"room-config-summary-whois"`
	ownersTreeView          gtki.TreeView       `gtk-widget:"room-config-summary-owners-tree"`
	adminsTreeView          gtki.TreeView       `gtk-widget:"room-config-summary-admins-tree"`
	maxHistoryFetch         gtki.Label          `gtk-widget:"room-config-summary-maxhistoryfetch"`
	maxOccupants            gtki.Label          `gtk-widget:"room-config-summary-maxoccupants"`
	enableArchiving         gtki.Label          `gtk-widget:"room-config-summary-archive-label"`
	autojoinCheckButton     gtki.CheckButton    `gtk-widget:"room-config-autojoin"`

	ownersTreeModel gtki.ListStore
	adminsTreeModel gtki.ListStore
}

func (c *mucRoomConfigComponent) newRoomConfigSummaryPage() mucRoomConfigPage {
	p := &roomConfigSummaryPage{autoJoin: c.autoJoin}
	p.roomConfigPageBase = c.newConfigPage("summary", "MUCRoomConfigPageSummary", p, map[string]interface{}{
		"on_autojoin_toggled": func() {
			c.updateAutoJoin(p.autojoinCheckButton.GetActive())
		},
		"go_basic_information_page": func() {
			c.setCurrentPage(roomConfigInformationPageIndex)
		},
		"go_access_page": func() {
			c.setCurrentPage(roomConfigAccessPageIndex)
		},
		"go_permissions_page": func() {
			c.setCurrentPage(roomConfigPermissionsPageIndex)
		},
		"go_occupants_page": func() {
			c.setCurrentPage(roomConfigOccupantsPageIndex)
		},
		"go_other_page": func() {
			c.setCurrentPage(roomConfigOthersPageIndex)
		},
	})

	p.onRefresh.add(p.onSummaryPageRefresh)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.basicInformation)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.access)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.permissions)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.occupants)
	mucStyles.setRoomConfigSummarySectionLinkButtonStyle(p.others)

	p.ownersTreeModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)
	p.adminsTreeModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)

	p.ownersTreeView.SetModel(p.ownersTreeModel)
	p.adminsTreeView.SetModel(p.adminsTreeModel)

	return p
}

func (p *roomConfigSummaryPage) handleDescriptionField() {
	if p.form.Description == "" {
		p.descriptionScrollWindow.Hide()
		p.descriptionNotAssigned.SetVisible(true)
		return
	}

	p.descriptionScrollWindow.Show()
	p.descriptionNotAssigned.SetVisible(false)
}

func (p *roomConfigSummaryPage) onSummaryPageRefresh() {
	p.autojoinCheckButton.SetActive(p.autoJoin)

	// Basic information
	setLabelText(p.title, setDefaultLabelText(p.form.Title))
	p.handleDescriptionField()
	setTextViewText(p.description, setDefaultLabelText(p.form.Description))
	setLabelText(p.descriptionNotAssigned, setDefaultLabelText(p.form.Description))
	setLabelText(p.language, supportedLanguageDescription(p.form.Language))
	setLabelText(p.includePublicList, getStringFromActiveValue(p.form.Public))
	setLabelText(p.persistent, getStringFromActiveValue(p.form.Persistent))
	// Access
	setLabelText(p.password, passwordMaskBasedOn(p.form.PasswordProtected))
	setLabelText(p.allowInviteUsers, getStringFromActiveValue(p.form.OccupantsCanInvite))
	setLabelText(p.onlyMembers, getStringFromActiveValue(p.form.MembersOnly))
	// Permissions
	setLabelText(p.allowSetRoomSubject, getStringFromActiveValue(p.form.OccupantsCanChangeSubject))
	setLabelText(p.moderatedRoom, getStringFromActiveValue(p.form.Moderated))
	setLabelText(p.whoIs, configOptionToFriendlyMessage(p.form.Whois.CurrentValue()))

	// Occupants
	summaryValueOfOccupantList(p.ownersTreeModel, p.form.Owners)
	summaryValueOfOccupantList(p.adminsTreeModel, p.form.Admins)

	// Other settings
	setLabelText(p.maxHistoryFetch, summaryValueForConfigOption(p.form.MaxHistoryFetch.CurrentValue()))
	setLabelText(p.maxOccupants, summaryValueForConfigOption(p.form.MaxOccupantsNumber.CurrentValue()))
	setLabelText(p.enableArchiving, getStringFromActiveValue(p.form.Logged))
}

func summaryValueOfOccupantList(model gtki.ListStore, items []jid.Any) {
	model.Clear()
	for _, j := range items {
		iter := model.Append()
		model.SetValue(iter, 0, j.String())
	}
}

func summaryValueForConfigOption(v string) string {
	if v == "" {
		v = muc.RoomConfigOptionNone
	}
	return configOptionToFriendlyMessage(v)
}

func passwordMaskBasedOn(value bool) string {
	if value {
		return "**********"
	}
	return i18n.Local("Not assigned")
}

func getStringFromActiveValue(value bool) string {
	if value {
		return i18n.Local("Yes")
	}
	return i18n.Local("No")
}

func setDefaultLabelText(label string) string {
	if label != "" {
		return label
	}
	return i18n.Local("Not assigned")
}

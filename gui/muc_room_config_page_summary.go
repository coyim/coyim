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

	overlay                 gtki.Overlay     `gtk-widget:"room-config-overlay"`
	infoSectionLabel        gtki.Label       `gtk-widget:"room-config-information-title"`
	title                   gtki.Label       `gtk-widget:"room-config-summary-title"`
	description             gtki.Label       `gtk-widget:"room-config-summary-description"`
	language                gtki.Label       `gtk-widget:"room-config-summary-language"`
	includePublicList       gtki.CheckButton `gtk-widget:"room-config-summary-public"`
	persistent              gtki.CheckButton `gtk-widget:"room-config-summary-persistent"`
	accessSectionLabel      gtki.Label       `gtk-widget:"room-config-access-title"`
	password                gtki.CheckButton `gtk-widget:"room-config-summary-password"`
	allowInviteUsers        gtki.CheckButton `gtk-widget:"room-config-summary-invite"`
	onlyMembers             gtki.CheckButton `gtk-widget:"room-config-summary-onlymembers"`
	permsisionsSectionLabel gtki.Label       `gtk-widget:"room-config-permissions-title"`
	allowSetRoomSubject     gtki.CheckButton `gtk-widget:"room-config-summary-changesubject"`
	moderatedRoom           gtki.CheckButton `gtk-widget:"room-config-summary-moderated"`
	whoIs                   gtki.Label       `gtk-widget:"room-config-summary-whois"`
	ownersSectionLabel      gtki.Label       `gtk-widget:"room-config-owners-title"`
	ownersTreeView          gtki.TreeView    `gtk-widget:"room-config-summary-owners-tree"`
	adminsSectionLabel      gtki.Label       `gtk-widget:"room-config-admins-title"`
	adminsTreeView          gtki.TreeView    `gtk-widget:"room-config-summary-admins-tree"`
	othersSectionLabel      gtki.Label       `gtk-widget:"room-config-others-title"`
	maxHistoryFetch         gtki.Label       `gtk-widget:"room-config-summary-maxhistoryfetch"`
	maxOccupants            gtki.Label       `gtk-widget:"room-config-summary-maxoccupants"`
	enableArchiving         gtki.CheckButton `gtk-widget:"room-config-summary-archive"`
	autojoinCheckButton     gtki.CheckButton `gtk-widget:"room-config-autojoin"`

	ownersTreeModel gtki.ListStore
	adminsTreeModel gtki.ListStore
}

func (c *mucRoomConfigComponent) newRoomConfigSummaryPage() mucRoomConfigPage {
	p := &roomConfigSummaryPage{autoJoin: c.autoJoin}
	p.roomConfigPageBase = c.newConfigPage("summary", "MUCRoomConfigPageSummary", p, map[string]interface{}{
		"on_autojoin_toggled": func() {
			c.updateAutoJoin(p.autojoinCheckButton.GetActive())
		},
	})

	p.onRefresh.add(p.onSummaryPageRefresh)

	mucStyles.setRoomConfigSummarySectionLabelStyle(p.infoSectionLabel)
	mucStyles.setRoomConfigSummarySectionLabelStyle(p.accessSectionLabel)
	mucStyles.setRoomConfigSummarySectionLabelStyle(p.permsisionsSectionLabel)
	mucStyles.setRoomConfigSummarySectionLabelStyle(p.adminsSectionLabel)
	mucStyles.setRoomConfigSummarySectionLabelStyle(p.ownersSectionLabel)
	mucStyles.setRoomConfigSummarySectionLabelStyle(p.othersSectionLabel)

	mucStyles.setRoomConfigSummaryRoomTitleLabelStyle(p.title)
	mucStyles.setRoomConfigSummaryRoomDescriptionLabelStyle(p.description)

	p.ownersTreeModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)
	p.adminsTreeModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)

	p.ownersTreeView.SetModel(p.ownersTreeModel)
	p.adminsTreeView.SetModel(p.adminsTreeModel)

	p.overlay.AddOverlay(p.loadingOverlay.overlay)

	return p
}

func (p *roomConfigSummaryPage) onSummaryPageRefresh() {
	p.autojoinCheckButton.SetActive(p.autoJoin)

	// Basic information
	setLabelText(p.title, p.form.Title)
	setLabelText(p.description, p.form.Description)
	setLabelText(p.language, i18n.Localf("%s (%s)", p.form.Language, getLanguage(p.form.Language)))
	p.includePublicList.SetActive(p.form.Public)
	p.includePublicList.SetSensitive(false)
	p.persistent.SetActive(p.form.Persistent)
	p.persistent.SetSensitive(false)

	// Access
	p.password.SetActive(p.form.PasswordProtected)
	p.password.SetSensitive(false)
	p.allowInviteUsers.SetActive(p.form.OccupantsCanInvite)
	p.allowInviteUsers.SetSensitive(false)
	p.onlyMembers.SetActive(p.form.MembersOnly)
	p.onlyMembers.SetSensitive(false)

	// Permissions
	p.allowSetRoomSubject.SetActive(p.form.OccupantsCanChangeSubject)
	p.allowSetRoomSubject.SetSensitive(false)
	p.moderatedRoom.SetActive(p.form.Moderated)
	p.moderatedRoom.SetSensitive(false)
	setLabelText(p.whoIs, configOptionToFriendlyMessage(p.form.Whois.CurrentValue()))

	// Occupants
	summaryValueOfOccupantList(p.ownersTreeModel, p.form.Owners)
	summaryValueOfOccupantList(p.adminsTreeModel, p.form.Admins)

	// Other settings
	setLabelText(p.maxHistoryFetch, summaryValueForConfigOption(p.form.MaxHistoryFetch.CurrentValue()))
	setLabelText(p.maxOccupants, summaryValueForConfigOption(p.form.MaxOccupantsNumber.CurrentValue()))
	p.enableArchiving.SetActive(p.form.Logged)
	p.enableArchiving.SetSensitive(false)
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

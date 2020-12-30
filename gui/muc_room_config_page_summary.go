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
	box                     gtki.Box         `gtk-widget:"room-config-summary-page"`
	infoSectionLabel        gtki.Label       `gtk-widget:"room-config-information-title"`
	title                   gtki.Label       `gtk-widget:"room-config-summary-title"`
	description             gtki.Label       `gtk-widget:"room-config-summary-description"`
	language                gtki.Label       `gtk-widget:"room-config-summary-language"`
	includePublicList       gtki.Image       `gtk-widget:"room-config-summary-public"`
	persistent              gtki.Image       `gtk-widget:"room-config-summary-persistent"`
	accessSectionLabel      gtki.Label       `gtk-widget:"room-config-access-title"`
	password                gtki.Image       `gtk-widget:"room-config-summary-password"`
	allowInviteUsers        gtki.Image       `gtk-widget:"room-config-summary-invite"`
	onlyMembers             gtki.Image       `gtk-widget:"room-config-summary-onlymembers"`
	permsisionsSectionLabel gtki.Label       `gtk-widget:"room-config-permissions-title"`
	allowSetRoomSubject     gtki.Image       `gtk-widget:"room-config-summary-changesubject"`
	moderatedRoom           gtki.Image       `gtk-widget:"room-config-summary-moderated"`
	whoIs                   gtki.Label       `gtk-widget:"room-config-summary-whois"`
	occupantsSectionLabel   gtki.Label       `gtk-widget:"room-config-occupants-title"`
	ownersTreeView          gtki.TreeView    `gtk-widget:"room-config-summary-owners-tree"`
	adminsTreeView          gtki.TreeView    `gtk-widget:"room-config-summary-admins-tree"`
	othersSectionLabel      gtki.Label       `gtk-widget:"room-config-others-title"`
	maxHistoryFetch         gtki.Label       `gtk-widget:"room-config-summary-maxhistoryfetch"`
	maxOccupants            gtki.Label       `gtk-widget:"room-config-summary-maxoccupants"`
	enableArchiving         gtki.Image       `gtk-widget:"room-config-summary-archive"`
	autojoinCheckButton     gtki.CheckButton `gtk-widget:"room-config-autojoin"`
	notificationBox         gtki.Box         `gtk-widget:"notification-box"`

	ownersTreeModel gtki.ListStore
	adminsTreeModel gtki.ListStore
}

func (c *mucRoomConfigComponent) newRoomConfigSummaryPage() mucRoomConfigPage {
	p := &roomConfigSummaryPage{
		autoJoin: c.autoJoin,
	}

	builder := newBuilder("MUCRoomConfigPageSummary")
	panicOnDevError(builder.bindObjects(p))

	builder.ConnectSignals(map[string]interface{}{
		"on_autojoin_toggled": func() {
			c.updateAutoJoin(p.autojoinCheckButton.GetActive())
		},
	})

	p.roomConfigPageBase = c.newConfigPage(p.box, p.notificationBox)
	p.onRefresh(p.onSummaryPageRefresh)

	mucStyles.setRoomConfigSummarySectionLabelStyle(p.infoSectionLabel)
	mucStyles.setRoomConfigSummarySectionLabelStyle(p.accessSectionLabel)
	mucStyles.setRoomConfigSummarySectionLabelStyle(p.permsisionsSectionLabel)
	mucStyles.setRoomConfigSummarySectionLabelStyle(p.occupantsSectionLabel)
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
	setImageYesOrNo(p.includePublicList, p.form.Public)
	setImageYesOrNo(p.persistent, p.form.Persistent)

	// Access
	setImageYesOrNo(p.password, p.form.PasswordProtected)
	setImageYesOrNo(p.allowInviteUsers, p.form.OccupantsCanInvite)
	setImageYesOrNo(p.onlyMembers, p.form.MembersOnly)

	// Permissions
	setImageYesOrNo(p.allowSetRoomSubject, p.form.OccupantsCanChangeSubject)
	setImageYesOrNo(p.moderatedRoom, p.form.Moderated)
	setLabelText(p.whoIs, configOptionToFriendlyMessage(p.form.Whois.CurrentValue()))

	// Occupants
	summaryValueOfOccupantList(p.ownersTreeModel, p.form.Owners)
	summaryValueOfOccupantList(p.adminsTreeModel, p.form.Admins)

	// Other settings
	setLabelText(p.maxHistoryFetch, summaryValueForConfigOption(p.form.MaxHistoryFetch.CurrentValue()))
	setLabelText(p.maxOccupants, summaryValueForConfigOption(p.form.MaxOccupantsNumber.CurrentValue()))
	setImageYesOrNo(p.enableArchiving, p.form.Logged)
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

func setImageYesOrNo(img gtki.Image, v bool) {
	icon := "no"
	if v {
		icon = "yes"
	}

	img.SetFromIconName("gtk-"+icon, gtki.ICON_SIZE_BUTTON)
}

package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryPage struct {
	*roomConfigPageBase
	autoJoin bool

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
	othersSectionLabel      gtki.Label       `gtk-widget:"room-config-others-title"`
	maxHistoryFetch         gtki.Label       `gtk-widget:"room-config-summary-maxhistoryfetch"`
	maxOccupants            gtki.Label       `gtk-widget:"room-config-summary-maxoccupants"`
	enableArchiving         gtki.Image       `gtk-widget:"room-config-summary-archive"`
	autojoinCheckButton     gtki.CheckButton `gtk-widget:"room-config-autojoin"`
	notificationBox         gtki.Box         `gtk-widget:"notification-box"`
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

	return p
}

func (p *roomConfigSummaryPage) onSummaryPageRefresh() {
	p.autojoinCheckButton.SetActive(p.autoJoin)

	// Basic information
	setLabelText(p.title, p.form.Title)
	setLabelText(p.description, p.form.Description)
	setLabelText(p.language, p.form.Language)
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

	// TODO: Occupants

	// Other settings
	setLabelText(p.maxHistoryFetch, summaryValueForConfigOption(p.form.MaxHistoryFetch.CurrentValue()))
	setLabelText(p.maxOccupants, summaryValueForConfigOption(p.form.MaxOccupantsNumber.CurrentValue()))
	setImageYesOrNo(p.enableArchiving, p.form.Logged)
}

func (p *roomConfigSummaryPage) collectData() {
	// Nothing to do, just implement the interface
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

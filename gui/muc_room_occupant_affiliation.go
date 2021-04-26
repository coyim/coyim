package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
)

func (r *roomViewRosterInfo) onChangeAffiliation() {
	av := r.newOccupantAffiliationUpdateView(r.account, r.roomID, r.occupant)
	av.showDialog()
}

type occupantAffiliationUpdateView struct {
	account        *account
	roomID         jid.Bare
	occupant       *muc.Occupant
	selfOccupant   *muc.Occupant
	rosterInfoView *roomViewRosterInfo

	dialog               gtki.Dialog      `gtk-widget:"affiliation-dialog"`
	contentBox           gtki.Box         `gtk-widget:"affiliation-content-box"`
	titleLabel           gtki.Label       `gtk-widget:"affiliation-title-label"`
	optionsDisabledLabel gtki.Label       `gtk-widget:"affiliation-options-disabled-label"`
	ownerOption          gtki.RadioButton `gtk-widget:"affiliation-owner"`
	adminOption          gtki.RadioButton `gtk-widget:"affiliation-admin"`
	memberOption         gtki.RadioButton `gtk-widget:"affiliation-member"`
	noneOption           gtki.RadioButton `gtk-widget:"affiliation-none"`
	reasonEntry          gtki.TextView    `gtk-widget:"affiliation-reason-entry"`
	applyButton          gtki.Button      `gtk-widget:"affiliation-apply-button"`
}

func (r *roomViewRosterInfo) newOccupantAffiliationUpdateView(a *account, roomID jid.Bare, o *muc.Occupant) *occupantAffiliationUpdateView {
	av := &occupantAffiliationUpdateView{
		account:        a,
		roomID:         roomID,
		rosterInfoView: r,
		occupant:       o,
		selfOccupant:   r.selfOccupant,
	}

	av.initBuilder()
	av.initDefaults()
	av.initRadioButtonsValues()

	return av
}

func (av *occupantAffiliationUpdateView) initBuilder() {
	builder := newBuilder("MUCRoomAffiliationDialog")
	panicOnDevError(builder.bindObjects(av))

	builder.ConnectSignals(map[string]interface{}{
		"on_cancel":    av.closeDialog,
		"on_apply":     av.onApply,
		"on_key_press": av.onKeyPress,
		"on_toggled":   av.onRadioButtonToggled,
	})
}

func (av *occupantAffiliationUpdateView) initDefaults() {
	av.dialog.SetTransientFor(av.rosterInfoView.parentWindow())

	av.titleLabel.SetText(av.affiliationChangingLabelText())

	mucStyles.setFormSectionLabelStyle(av.titleLabel)
	mucStyles.setHelpTextStyle(av.contentBox)
}

// initRadioButtonsValues MUST be called from UI thread
func (av *occupantAffiliationUpdateView) initRadioButtonsValues() {
	if av.selfOccupant.Affiliation.IsAdmin() {
		av.adminOption.SetSensitive(false)
		av.ownerOption.SetSensitive(false)
		av.optionsDisabledLabel.SetVisible(true)
	}

	switch av.occupant.Affiliation.(type) {
	case *data.OwnerAffiliation:
		av.ownerOption.SetActive(true)
	case *data.AdminAffiliation:
		av.adminOption.SetActive(true)
	case *data.MemberAffiliation:
		av.memberOption.SetActive(true)
	case *data.NoneAffiliation:
		av.noneOption.SetActive(true)
	}
}

// onKeyPress MUST be called from the UI thread
func (av *occupantAffiliationUpdateView) onKeyPress(_ gtki.Widget, ev gdki.Event) {
	if isNormalEnter(g.gdk.EventKeyFrom(ev)) {
		av.onApply()
	}
}

// onRadioButtonToggled MUST be called from the UI thread
func (av *occupantAffiliationUpdateView) onRadioButtonToggled() {
	s := av.occupant.Affiliation.IsDifferentFrom(av.affiliationBasedOnSelectedRadio())
	av.applyButton.SetSensitive(s)
}

// onApply MUST be called from the UI thread
func (av *occupantAffiliationUpdateView) onApply() {
	affiliation := av.affiliationBasedOnSelectedRadio()
	reason := getTextViewText(av.reasonEntry)

	go av.rosterInfoView.updateOccupantAffiliation(av.occupant, affiliation, reason)

	av.closeDialog()
}

// affiliationBasedOnSelectedRadio MUST be called from the UI thread
func (av *occupantAffiliationUpdateView) affiliationBasedOnSelectedRadio() data.Affiliation {
	switch {
	case av.ownerOption.GetActive():
		return &data.OwnerAffiliation{}
	case av.adminOption.GetActive():
		return &data.AdminAffiliation{}
	case av.memberOption.GetActive():
		return &data.MemberAffiliation{}
	}
	return &data.NoneAffiliation{}
}

// closeDialog MUST be called from the UI thread
func (av *occupantAffiliationUpdateView) closeDialog() {
	av.dialog.Destroy()
}

// showDialog MUST be called from the UI thread
func (av *occupantAffiliationUpdateView) showDialog() {
	av.dialog.Show()
}

func (av *occupantAffiliationUpdateView) affiliationChangingLabelText() string {
	affiliation := av.occupant.Affiliation
	nickname := av.occupant.Nickname

	switch {
	case affiliation.IsOwner():
		return i18n.Localf("You are changing the position of %s from owner to:", nickname)
	case affiliation.IsAdmin():
		return i18n.Localf("You are changing the position of %s from administrator to:", nickname)
	case affiliation.IsMember():
		return i18n.Localf("You are changing the position of %s from member to:", nickname)
	}

	return i18n.Localf("You are changing the position of %s to:", nickname)
}

func occupantAffiliationName(a data.Affiliation) string {
	switch a.(type) {
	case *data.OwnerAffiliation:
		return i18n.Local("Owner")
	case *data.AdminAffiliation:
		return i18n.Local("Admin")
	case *data.MemberAffiliation:
		return i18n.Local("Member")
	case *data.OutcastAffiliation:
		return i18n.Local("Outcast")
	case *data.NoneAffiliation:
		return i18n.Local("None")
	}
	return ""
}

func occupantRoleName(a data.Role) string {
	switch a.(type) {
	case *data.ModeratorRole:
		return i18n.Local("Moderator")
	case *data.ParticipantRole:
		return i18n.Local("Participant")
	case *data.VisitorRole:
		return i18n.Local("Visitor")
	case *data.NoneRole:
		return i18n.Local("None")
	}
	return ""
}

func affiliationUpdateErrorMessage(err error) string {
	if err == session.ErrUpdateOccupantResponse {
		return i18n.Local("We couldn't update the occupant affiliation because, either you don't have permission to do it or the server is busy. Please try again.")
	}
	return i18n.Local("An error occurred when updating the occupant affiliation. Please try again.")
}

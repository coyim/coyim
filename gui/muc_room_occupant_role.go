package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
)

func (r *roomViewRosterInfo) onChangeRole() {
	rv := r.newOccupantRoleUpdateView(r.account, r.roomID, r.occupant)
	rv.showDialog()
}

type occupantRoleUpdateView struct {
	account        *account
	roomID         jid.Bare
	occupant       *muc.Occupant
	rosterInfoView *roomViewRosterInfo

	dialog         gtki.Dialog      `gtk-widget:"role-dialog"`
	roleLabel      gtki.Label       `gtk-widget:"role-type-label"`
	moderatorRadio gtki.RadioButton `gtk-widget:"role-moderator"`
	reasonLabel    gtki.Label       `gtk-widget:"role-reason-label"`
	reasonEntry    gtki.TextView    `gtk-widget:"role-reason-entry"`
	applyButton    gtki.Button      `gtk-widget:"role-apply-button"`
}

func (r *roomViewRosterInfo) newOccupantRoleUpdateView(a *account, roomID jid.Bare, o *muc.Occupant) *occupantRoleUpdateView {
	rv := &occupantRoleUpdateView{
		account:        a,
		roomID:         roomID,
		rosterInfoView: r,
		occupant:       o,
	}

	rv.initBuilder()
	rv.initDefaults()

	return rv
}

func (rv *occupantRoleUpdateView) initBuilder() {
	builder := newBuilder("MUCRoomRoleDialog")
	panicOnDevError(builder.bindObjects(rv))

	builder.ConnectSignals(map[string]interface{}{
		"on_cancel":    rv.closeDialog,
		"on_apply":     rv.onApply,
		"on_key_press": rv.onKeyPress,
	})
}

func (rv *occupantRoleUpdateView) onKeyPress(_ gtki.Widget, ev gdki.Event) {
	if isNormalEnter(g.gdk.EventKeyFrom(ev)) {
		rv.onApply()
	}
}

func (rv *occupantRoleUpdateView) initDefaults() {
	rv.dialog.SetTransientFor(rv.rosterInfoView.parentWindow())
	mucStyles.setFormSectionLabelStyle(rv.roleLabel)
	// TODO: This button is enabled until another radio option will be added
	rv.applyButton.SetSensitive(true)

	switch rv.occupant.Role.(type) {
	case *data.ModeratorRole:
		rv.moderatorRadio.SetActive(true)
	}
}

// onApply MUST be called from the UI thread
func (rv *occupantRoleUpdateView) onApply() {
	go rv.rosterInfoView.updateOccupantRole(rv.occupant, rv.getRoleBasedOnRadioSelected(), getTextViewText(rv.reasonEntry))
	rv.closeDialog()
}

func (rv *occupantRoleUpdateView) getRoleBasedOnRadioSelected() data.Role {
	switch {
	case rv.moderatorRadio.GetActive():
		return &data.ModeratorRole{}
	default:
		return &data.NoneRole{}
	}
}

// close MUST be called from the UI thread
func (rv *occupantRoleUpdateView) closeDialog() {
	rv.dialog.Destroy()
}

// show MUST be called from the UI thread
func (rv *occupantRoleUpdateView) showDialog() {
	rv.dialog.Show()
}

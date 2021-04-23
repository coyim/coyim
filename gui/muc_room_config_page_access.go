package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAccessPage struct {
	*roomConfigPageBase
	roomPassword *passwordConfirmationComponent

	roomPasswordBox  gtki.Box    `gtk-widget:"room-password-box"`
	roomMembersOnly  gtki.Switch `gtk-widget:"room-membersonly"`
	roomAllowInvites gtki.Switch `gtk-widget:"room-allowinvites"`
}

func (c *mucRoomConfigComponent) newRoomConfigAccessPage() mucRoomConfigPage {
	p := &roomConfigAccessPage{}
	p.roomConfigPageBase = c.newConfigPage("access", "MUCRoomConfigPageAccess", p, nil)

	p.initPasswordComponent()
	p.initDefaultValues()

	return p
}

func (p *roomConfigAccessPage) initPasswordComponent() {
	p.roomPassword = p.u.createPasswordConfirmationComponent()
	p.roomPasswordBox.Add(p.roomPassword.contentBox())
	p.doAfterRefresh.add(p.roomPassword.onShowConfirmPasswordBasedOnMatchError)
}

func (p *roomConfigAccessPage) initDefaultValues() {
	p.roomPassword.setPassword(p.form.Password)
	setSwitchActive(p.roomMembersOnly, p.form.MembersOnly)
	setSwitchActive(p.roomAllowInvites, p.form.OccupantsCanInvite)
}

func (p *roomConfigAccessPage) collectData() {
	p.form.Password = p.roomPassword.currentPassword()
	p.form.PasswordProtected = p.form.Password != ""
	p.form.MembersOnly = getSwitchActive(p.roomMembersOnly)
	p.form.OccupantsCanInvite = getSwitchActive(p.roomAllowInvites)
}

func (p *roomConfigAccessPage) isInvalid() bool {
	return !p.roomPassword.passwordsMatch()
}

func (p *roomConfigAccessPage) showValidationErrors() {
	p.roomPassword.changeConfirmPasswordEntryStyle()
	p.roomPassword.onShowConfirmPasswordBasedOnMatchError()
	p.roomPassword.focusConfirm()
}

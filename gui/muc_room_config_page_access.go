package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAccessPage struct {
	*roomConfigPageBase

	roomPasswordBox  gtki.Box    `gtk-widget:"room-password-box"`
	roomMembersOnly  gtki.Switch `gtk-widget:"room-membersonly"`
	roomAllowInvites gtki.Switch `gtk-widget:"room-allowinvites"`
}

func (c *mucRoomConfigComponent) newRoomConfigAccessPage() mucRoomConfigPage {
	p := &roomConfigAccessPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigAccess, "MUCRoomConfigPageAccess", p, nil)

	p.initDefaultValues()

	return p
}

func (p *roomConfigAccessPage) initDefaultValues() {
	setSwitchActive(p.roomMembersOnly, p.form.MembersOnly)
	setSwitchActive(p.roomAllowInvites, p.form.OccupantsCanInvite)
}

func (p *roomConfigAccessPage) collectData() {
	p.form.PasswordProtected = p.form.Password != ""
	p.form.MembersOnly = getSwitchActive(p.roomMembersOnly)
	p.form.OccupantsCanInvite = getSwitchActive(p.roomAllowInvites)
}

func (p *roomConfigAccessPage) isInvalid() bool {
	return p.roomConfigPageBase.isInvalid()
}

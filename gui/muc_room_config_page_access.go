package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigAccessPage struct {
	*roomConfigPageBase

	roomPassword     gtki.Entry  `gtk-widget:"room-password"`
	roomMembersOnly  gtki.Switch `gtk-widget:"room-membersonly"`
	roomAllowInvites gtki.Switch `gtk-widget:"room-allowinvites"`
}

func (c *mucRoomConfigComponent) newRoomConfigAccessPage() mucRoomConfigPage {
	p := &roomConfigAccessPage{}
	p.roomConfigPageBase = c.newConfigPage("access", "MUCRoomConfigPageAccess", p, nil)

	p.initDefaultValues()

	return p
}

func (p *roomConfigAccessPage) initDefaultValues() {
	setEntryText(p.roomPassword, p.form.Password)
	setSwitchActive(p.roomMembersOnly, p.form.MembersOnly)
	setSwitchActive(p.roomAllowInvites, p.form.OccupantsCanInvite)
}

func (p *roomConfigAccessPage) collectData() {
	p.form.Password = getEntryText(p.roomPassword)
	p.form.PasswordProtected = p.form.Password != ""
	p.form.MembersOnly = getSwitchActive(p.roomMembersOnly)
	p.form.OccupantsCanInvite = getSwitchActive(p.roomAllowInvites)
}

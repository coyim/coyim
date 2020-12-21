package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigPage interface {
	getPageView() gtki.Box
	collectData()
	onRefresh(func())
	refresh()
	nofityError(string)
}

type roomConfigPageBase struct {
	u             *gtkUI
	box           gtki.Box
	notifications *notifications
	refreshList   []func()
	form          *muc.RoomConfigForm
	log           coylog.Logger
}

func (c *mucRoomConfigComponent) newConfigPage(b gtki.Box, nb gtki.Box) *roomConfigPageBase {
	cp := &roomConfigPageBase{
		u:    c.u,
		box:  b,
		form: c.form,
		log:  c.log,
	}

	cp.notifications = c.u.newNotifications(nb)

	return cp
}

func (p *roomConfigPageBase) getPageView() gtki.Box {
	return p.box
}

func (p *roomConfigPageBase) collectData() {
	panic("developer error: collectData()")
}

func (p *roomConfigPageBase) onRefresh(f func()) {
	p.refreshList = append(p.refreshList, f)
}

func (p *roomConfigPageBase) refresh() {
	p.box.ShowAll()

	for _, f := range p.refreshList {
		f()
	}
}

func (p *roomConfigPageBase) nofityError(m string) {
	p.notifications.notifyOnError(m)
}

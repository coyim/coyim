package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigPage interface {
	getContent() gtki.Box
	isValid() bool
	collectData()
	onRefresh(func())
	refresh()
	showLoadingOverlay()
	hideLoadingOverlay()
}

type roomConfigPageBase struct {
	u              *gtkUI
	content        gtki.Box
	loadingOverlay *loadingOverlayComponent
	notifications  *notifications
	refreshList    []func()
	form           *muc.RoomConfigForm
	log            coylog.Logger
}

func (c *mucRoomConfigComponent) newConfigPage(content gtki.Box, nb gtki.Box) *roomConfigPageBase {
	return &roomConfigPageBase{
		u:              c.u,
		content:        content,
		notifications:  c.u.newNotifications(nb),
		loadingOverlay: c.u.newLoadingOverlayComponent(),
		form:           c.form,
		log:            c.log,
	}
}

func (p *roomConfigPageBase) getContent() gtki.Box {
	return p.content
}

func (p *roomConfigPageBase) isValid() bool {
	return true
}

func (p *roomConfigPageBase) collectData() {
	// Nothing to do, just implement the interface
}

func (p *roomConfigPageBase) onRefresh(f func()) {
	p.refreshList = append(p.refreshList, f)
}

func (p *roomConfigPageBase) refresh() {
	p.content.ShowAll()

	for _, f := range p.refreshList {
		f()
	}
}

func (p *roomConfigPageBase) clearErrors() {
	p.notifications.clearErrors()
}

func (p *roomConfigPageBase) nofityError(m string) {
	p.notifications.notifyOnError(m)
}

func (p *roomConfigPageBase) showLoadingOverlay() {
	p.loadingOverlay.show()
}

func (p *roomConfigPageBase) hideLoadingOverlay() {
	p.loadingOverlay.hide()
}

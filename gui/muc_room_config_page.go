package gui

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type mucRoomConfigPage interface {
	pageView() gtki.Box
	isValid() bool
	collectData()
	refresh()
	showLoadingOverlay()
	hideLoadingOverlay()
}

type roomConfigPageBase struct {
	u    *gtkUI
	form *muc.RoomConfigForm

	page              gtki.Box `gtk-widget:"room-config-page"`
	content           gtki.Box `gtk-widget:"room-config-page-content"`
	notificationsArea gtki.Box `gtk-widget:"notifications-box"`

	notifications  *notifications
	loadingOverlay *loadingOverlayComponent
	onRefresh      *callbacksSet

	log coylog.Logger
}

func (c *mucRoomConfigComponent) newConfigPage(pageID, pageTemplate string, page interface{}, signals map[string]interface{}) *roomConfigPageBase {
	p := &roomConfigPageBase{
		u:              c.u,
		loadingOverlay: c.u.newLoadingOverlayComponent(),
		onRefresh:      newCallbacksSet(),
		form:           c.form,
		log: c.log.WithFields(log.Fields{
			"page":     pageID,
			"template": pageTemplate,
		}),
	}

	builder := newBuilder("MUCRoomConfigPage")
	panicOnDevError(builder.bindObjects(p))

	p.notifications = c.u.newNotificationsComponent()
	p.notificationsArea.Add(p.notifications.widget())

	builder = newBuilder(pageTemplate)
	panicOnDevError(builder.bindObjects(page))
	builder.ConnectSignals(signals)

	pc, err := builder.GetObject(fmt.Sprintf("room-config-%s-page", pageID))
	if err != nil {
		panic(fmt.Sprintf("developer error: the ID for \"%s\" page doesn't exists", pageID))
	}

	p.content.Add(pc.(gtki.Box))

	return p
}

// pageView implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) pageView() gtki.Box {
	return p.page
}

// isValid implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) isValid() bool {
	return true
}

// Nothing to do, just implement the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) collectData() {}

// refresh MUST be called from the UI thread
func (p *roomConfigPageBase) refresh() {
	p.page.ShowAll()
	p.hideLoadingOverlay()
	p.onRefresh.invokeAll()
}

// clearErrors MUST be called from the ui thread
func (p *roomConfigPageBase) clearErrors() {
	p.notifications.clearErrors()
}

// nofityError MUST be called from the ui thread
func (p *roomConfigPageBase) nofityError(m string) {
	p.notifications.notifyOnError(m)
}

// showLoadingOverlay MUST be called from the ui thread
func (p *roomConfigPageBase) showLoadingOverlay() {
	p.loadingOverlay.show()
}

// hideLoadingOverlay MUST be called from the ui thread
func (p *roomConfigPageBase) hideLoadingOverlay() {
	p.loadingOverlay.hide()
}

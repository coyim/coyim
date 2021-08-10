package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewWindow struct {
	window              gtki.Window   `gtk-widget:"room-window"`
	overlay             gtki.Overlay  `gtk-widget:"room-overlay"`
	privacityWarningBox gtki.Box      `gtk-widget:"room-privacity-warnings-box"`
	content             gtki.Box      `gtk-widget:"room-main-box"`
	notificationsArea   gtki.Revealer `gtk-widget:"room-notifications-revealer"`
}

func (v *roomView) newRoomViewWindow() *roomViewWindow {
	vw := &roomViewWindow{}

	builder := newBuilder("MUCRoomWindow")
	panicOnDevError(builder.bindObjects(vw))

	builder.ConnectSignals(map[string]interface{}{
		"on_destroy_window": v.onDestroyWindow,
	})

	vw.window.SetTitle(i18n.Localf("%[1]s [%[2]s]", v.roomID(), v.account.Account()))
	mucStyles.setRoomWindowStyle(vw.window)

	return vw
}

// onNewNotificationAdded MUST be called from the UI thread
func (vw *roomViewWindow) onNewNotificationAdded() {
	if !vw.notificationsArea.GetRevealChild() {
		vw.notificationsArea.SetRevealChild(true)
	}
}

// onNoNotifications MUST be called from the UI thread
func (vw *roomViewWindow) onNoNotifications() {
	vw.notificationsArea.SetRevealChild(false)
}

// addContent MUST be called from the UI thread
func (vw *roomViewWindow) addContentWidget(c gtki.Widget) {
	vw.content.Add(c)
}

// removeContentWidget MUST be called from the UI thread
func (vw *roomViewWindow) removeContentWidget(c gtki.Widget) {
	vw.content.Remove(c)
}

// present MUST be called from the UI thread
func (vw *roomViewWindow) present() {
	vw.window.Present()
}

// show MUST be called from the UI thread
func (vw *roomViewWindow) show() {
	vw.window.Show()
}

// destroy MUST be called from the UI thread
func (vw *roomViewWindow) destroy() {
	vw.window.Destroy()
}

func (vw *roomViewWindow) view() gtki.Window {
	return vw.window
}

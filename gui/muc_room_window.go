package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewWindow struct {
	roomView *roomView

	window            gtki.Window   `gtk-widget:"room-window"`
	overlay           gtki.Overlay  `gtk-widget:"room-overlay"`
	privacyWarningBox gtki.Box      `gtk-widget:"room-privacy-warnings-box"`
	content           gtki.Box      `gtk-widget:"room-main-box"`
	notificationsArea gtki.Revealer `gtk-widget:"room-notifications-revealer"`
}

func (v *roomView) newRoomViewWindow() *roomViewWindow {
	vw := &roomViewWindow{
		roomView: v,
	}

	vw.loadUIDefinition()
	vw.initDefaults(v.u)

	return vw
}

func (vw *roomViewWindow) loadUIDefinition() {
	buildUserInterface("MUCRoomWindow", vw, vw.connectUISignals)
}

func (vw *roomViewWindow) connectUISignals(b *builder) {
	b.ConnectSignals(map[string]interface{}{
		"on_destroy_window": vw.roomView.onDestroyWindow,
		"on_before_delete":  vw.onBeforeWindowClose,
	})
}

func (vw *roomViewWindow) initDefaults(u *gtkUI) {
	vw.window.SetTitle(i18n.Localf("%[1]s [%[2]s]", vw.roomView.roomID(), vw.roomView.account.Account()))
	mucStyles.setRoomWindowStyle(vw.window)

	u.connectShortcutsMucRoomWindow(vw.window, func(_ gtki.Window) {
		_ = vw.onBeforeWindowClose()
	})
}

const (
	roomWindowCloseStopEvent     = true // This will stop calling all the signals attached to `delete-event`
	roomWindowCloseContinueEvent = false
)

// onBeforeWindowClose MUST be called from the UI thread
func (vw *roomViewWindow) onBeforeWindowClose() bool {
	if vw.roomView.isSelfOccupantInTheRoom() {
		vw.roomView.confirmWindowClose()
		return roomWindowCloseStopEvent
	}

	return roomWindowCloseContinueEvent
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

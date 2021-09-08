package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

// initNotifications MUST be called from the UI thread
func (v *roomView) initNotifications() {
	v.notifications = v.newRoomNotifications()
}

// onNewNotificationAdded MUST be called from the UI thread
func (v *roomView) onNewNotificationAdded() {
	v.window.onNewNotificationAdded()
}

// onNoNotifications MUST be called from the UI thread
func (v *roomView) onNoNotifications() {
	v.window.onNoNotifications()
}

type roomNotifications struct {
	u             *gtkUI
	notifications *notificationsComponent
	roomView      *roomView
}

func (v *roomView) newRoomNotifications() *roomNotifications {
	notifications := v.u.newNotificationsComponent()
	notifications.setStacked(true)

	return &roomNotifications{
		u:             v.u,
		notifications: notifications,
		roomView:      v,
	}
}

type roomNotificationAction struct {
	label        string
	responseType gtki.ResponseType
	signals      map[string]interface{}
}

type roomNotificationActions []roomNotificationAction

type roomNotificationOptions struct {
	message     string
	showTime    bool
	showSpinner bool
	closeable   bool
	mustRemain  bool
	actions     roomNotificationActions
}

// other MUST be called from the UI thread
func (rn *roomNotifications) other(n roomNotificationOptions) {
	rn.newNotification(gtki.MESSAGE_OTHER, n)
}

// info MUST be called from the UI thread
func (rn *roomNotifications) info(n roomNotificationOptions) {
	rn.newNotification(gtki.MESSAGE_INFO, n)
}

// warning MUST be called from the UI thread
func (rn *roomNotifications) warning(n roomNotificationOptions) {
	rn.newNotification(gtki.MESSAGE_WARNING, n)
}

// error MUST be called from the UI thread
func (rn *roomNotifications) error(n roomNotificationOptions) {
	rn.newNotification(gtki.MESSAGE_ERROR, n)
}

// newNotification MUST be called from the UI thread
func (rn *roomNotifications) newNotification(mt gtki.MessageType, n roomNotificationOptions) *notificationBar {
	nb := rn.u.newNotificationBar(n.message, mt)

	if n.showTime {
		nb = rn.u.newNotificationBarWithTime(n.message, mt)
	}

	if n.showSpinner {
		nb.showSpinner()
	}

	if n.closeable {
		nb.whenRequestedToClose(func() {
			rn.remove(nb)
		})
	}

	nb.mustRemain = n.mustRemain

	for _, action := range n.actions {
		nb.addAction(action.label, action.responseType, action.signals)
	}

	rn.notifications.add(nb)
	rn.roomView.onNewNotificationAdded()

	return nb
}

// remove MUST be called from the UI thread
func (rn *roomNotifications) remove(nb *notificationBar) {
	rn.notifications.remove(nb)

	if rn.notifications.hasNoMessages() {
		rn.roomView.onNoNotifications()
	}
}

// notificationsBox MUST be called from the UI thread
func (rn *roomNotifications) notificationsBox() gtki.Widget {
	return rn.notifications.contentBox()
}

// clearAll MUST be called from the UI thread
func (rn *roomNotifications) clearAll() {
	rn.notifications.clearAll()
}

// clearErrors implements the "canNotifyErrors" interface
// clearErrors MUST be called from the UI thread
func (rn *roomNotifications) clearErrors() {
	rn.notifications.clearErrors()
}

// notifyOnError implements the "canNotifyErrors" interface
// notifyOnError MUST be called from the UI thread
func (rn *roomNotifications) notifyOnError(err string) {
	rn.error(roomNotificationOptions{
		message: err,
	})
}

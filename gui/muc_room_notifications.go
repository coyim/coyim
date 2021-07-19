package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

func (v *roomView) initNotifications() {
	v.notifications = v.newRoomNotifications()
	v.notificationsArea.Add(v.notifications.notificationsBox())
}

func (v *roomView) onNewNotificationAdded() {
	if !v.notificationsArea.GetRevealChild() {
		v.notificationsArea.SetRevealChild(true)
	}
}

func (v *roomView) onNoNotifications() {
	v.notificationsArea.SetRevealChild(false)
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
	messageType gtki.MessageType
	showTime    bool
	closeable   bool
	actions     roomNotificationActions
}

func (rn *roomNotifications) info(n roomNotificationOptions) {
	n.messageType = gtki.MESSAGE_INFO
	rn.newNotification(n)
}

func (rn *roomNotifications) warning(n roomNotificationOptions) {
	n.messageType = gtki.MESSAGE_WARNING
	rn.newNotification(n)
}

func (rn *roomNotifications) error(n roomNotificationOptions) {
	n.messageType = gtki.MESSAGE_ERROR
	rn.newNotification(n)
}

func (rn *roomNotifications) newNotification(n roomNotificationOptions) {
	nb := rn.u.newNotificationBar(n.message, n.messageType)

	if n.showTime {
		nb = rn.u.newNotificationBarWithTime(n.message, n.messageType)
	}

	if n.closeable {
		nb.whenRequestedToClose(func() {
			rn.remove(nb)
		})
	}

	for _, action := range n.actions {
		nb.addAction(action.label, action.responseType, action.signals)
	}

	rn.notifications.add(nb)
	rn.roomView.onNewNotificationAdded()
}

func (rn *roomNotifications) remove(nb *notificationBar) {
	rn.notifications.remove(nb)

	if rn.notifications.hasNoMessages() {
		rn.roomView.onNoNotifications()
	}
}

func (rn *roomNotifications) notificationsBox() gtki.Widget {
	return rn.notifications.contentBox()
}

func (rn *roomNotifications) clearErrors() {
	rn.notifications.clearErrors()
}

func (rn *roomNotifications) clearAll() {
	rn.notifications.clearAll()
}

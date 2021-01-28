package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

func (v *roomView) initNotifications() {
	v.notifications = v.newRoomNotifications()
	v.notificationsArea.Add(v.notifications.getNotificationsBox())
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

func (rn *roomNotifications) info(msg string) {
	rn.newNotification(msg, gtki.MESSAGE_INFO)
}

func (rn *roomNotifications) warning(msg string) {
	rn.newNotification(msg, gtki.MESSAGE_WARNING)
}

func (rn *roomNotifications) error(msg string) {
	rn.newNotification(msg, gtki.MESSAGE_ERROR)
}

func (rn *roomNotifications) newNotification(text string, messageType gtki.MessageType) {
	nb := rn.u.newNotificationBar(text, messageType)
	nb.setClosable(true)
	rn.add(nb)

	rn.roomView.onNewNotificationAdded()
}

func (rn *roomNotifications) add(nb *notificationBar) {
	if nb.isClosable() {
		nb.onClose(func() {
			rn.remove(nb)
		})
	}
	rn.notifications.add(nb)
}

func (rn *roomNotifications) remove(nb *notificationBar) {
	rn.notifications.remove(nb)

	if rn.notifications.hasNoMessages() {
		rn.roomView.onNoNotifications()
	}
}

func (rn *roomNotifications) getNotificationsBox() gtki.Widget {
	return rn.notifications.getBox()
}

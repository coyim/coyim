package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomNotifications struct {
	u             *gtkUI
	notifications *notifications
}

func (v *roomView) newRoomNotifications() *roomNotifications {
	notifications := v.u.newNotificationsComponent()
	notifications.setStacked(true)

	return &roomNotifications{
		u:             v.u,
		notifications: notifications,
	}
}

func (rn *roomNotifications) info(msg string) {
	nc := rn.u.newInfoBarComponent(msg, gtki.MESSAGE_INFO)
	nc.setClosable(true)
	rn.add(nc)
}

func (rn *roomNotifications) add(nc withNotification) {
	if nc.isClosable() {
		nc.onClose(func() {
			rn.remove(nc)
		})
	}
	rn.notifications.add(nc)
}

func (rn *roomNotifications) remove(nc withNotification) {
	rn.notifications.remove(nc)
}

func (rn *roomNotifications) widget() gtki.Widget {
	return rn.notifications.widget()
}

package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomNotifications struct {
	u             *gtkUI
	notifications *notifications
}

func (v *roomView) newRoomNotifications() *roomNotifications {
	rn := &roomNotifications{
		u:             v.u,
		notifications: v.u.newNotificationsComponent(),
	}

	return rn
}

func (rn *roomNotifications) info(msg string) {
	nc := rn.u.newInfoBarComponent(msg, gtki.MESSAGE_INFO)
	nc.setClosable(true)
	rn.notifications.add(nc)
}

func (rn *roomNotifications) add(nc withNotification) {
	rn.notifications.add(nc)
}

func (rn *roomNotifications) remove(nc withNotification) {
	rn.notifications.remove(nc)
}

func (rn *roomNotifications) widget() gtki.Widget {
	return rn.notifications.widget()
}

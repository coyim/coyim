package gui

import (
	"time"

	"github.com/coyim/gotk3adapter/gtki"
)

type notificationsComponent struct {
	u             *gtkUI
	box           gtki.Box
	notifications []*notificationBar
	stacked       bool
}

func (u *gtkUI) newNotificationsComponent() *notificationsComponent {
	b, _ := g.gtk.BoxNew(gtki.VerticalOrientation, 0)

	n := &notificationsComponent{
		u:   u,
		box: b,
	}

	return n
}

func (n *notificationsComponent) contentBox() gtki.Widget {
	return n.box
}

func (n *notificationsComponent) setStacked(v bool) {
	n.stacked = v
}

// add MUST be called from the ui thread
func (n *notificationsComponent) add(nb *notificationBar) {
	if !n.stacked {
		n.clearAll()
	}

	n.notifications = append(n.notifications, nb)

	n.box.PackStart(nb.view(), true, false, 0)
	n.box.ShowAll()
}

// remove MUST be called from the ui thread
func (n *notificationsComponent) remove(nb *notificationBar) {
	notifications := []*notificationBar{}
	for _, nbx := range n.notifications {
		if nb != nbx {
			notifications = append(notifications, nbx)
		}
	}

	n.notifications = notifications
	n.box.Remove(nb.view())
}

// clearAll MUST be called from the ui thread
func (n *notificationsComponent) clearAll() {
	notifications := n.notifications
	for _, nb := range notifications {
		n.remove(nb)
	}
}

// clearMessagesByType MUST be called from the ui thread
func (n *notificationsComponent) clearMessagesByType(mt gtki.MessageType) {
	notifications := n.notifications
	for _, nb := range notifications {
		if nb.messageType == mt {
			n.remove(nb)
		}
	}
}

// notify MUST be called from the UI thread
func (n *notificationsComponent) notify(text string, mt gtki.MessageType) {
	n.add(n.u.newNotificationBar(text, mt))
}

// warning MUST be called from the UI thread
func (n *notificationsComponent) warning(text string) {
	n.notify(text, gtki.MESSAGE_WARNING)
}

// error MUST be called from the UI thread
func (n *notificationsComponent) error(text string) {
	n.notify(text, gtki.MESSAGE_ERROR)
}

// info MUST be called from the ui thread
func (n *notificationsComponent) info(text string) {
	n.notify(text, gtki.MESSAGE_INFO)
}

// question MUST be called from the ui thread
func (n *notificationsComponent) question(text string) {
	n.notify(text, gtki.MESSAGE_QUESTION)
}

// message MUST be called from the ui thread
func (n *notificationsComponent) message(text string) {
	n.notify(text, gtki.MESSAGE_OTHER)
}

// notifyOnError is an alias for the "error" method and also
// implements the "canNotifyErrors" interface
//
// notifyOnError MUST be called from the ui thread
func (n *notificationsComponent) notifyOnError(err string) {
	n.error(err)
}

// clearErrors is an alias for the "clear" method and also
// implements the "canNotifyErrors" interface
//
// clearErrors MUST be called from the ui thread
func (n *notificationsComponent) clearErrors() {
	n.clearMessagesByType(gtki.MESSAGE_ERROR)
}

// hasNoMessages returns a boolean indicating if the notifications
// component has no messages
//
// hasNoMessages MUST be called from the ui thread
func (n *notificationsComponent) hasNoMessages() bool {
	return len(n.notifications) == 0
}

type notificationBar struct {
	*infoBarComponent
}

func (u *gtkUI) newNotificationBar(text string, messageType gtki.MessageType) *notificationBar {
	return &notificationBar{
		u.newInfoBarComponent(text, messageType),
	}
}

func (u *gtkUI) newNotificationBarWithTime(text string, messageType gtki.MessageType) *notificationBar {
	nb := u.newNotificationBar(text, messageType)
	nb.setTickerTime(time.Now())

	return nb
}

package gui

import "github.com/coyim/gotk3adapter/gtki"

type notification interface {
	widget
	message
}

type notifications struct {
	u        *gtkUI
	box      gtki.Box
	messages []notification
	stacked  bool
}

func (u *gtkUI) newNotificationsComponent() *notifications {
	b, _ := g.gtk.BoxNew(gtki.VerticalOrientation, 0)

	n := &notifications{
		u:   u,
		box: b,
	}

	return n
}

func (n *notifications) widget() gtki.Widget {
	return n.box
}

// add MUST be called from the ui thread
func (n *notifications) add(m notification) {
	if !n.stacked {
		n.clearAll()
	}

	n.messages = append(n.messages, m)

	n.box.Add(m.widget())
	n.box.ShowAll()
}

// remove MUST be called from the ui thread
func (n *notifications) remove(w gtki.Widget) {
	n.box.Remove(w)
}

// clearAll MUST be called from the ui thread
func (n *notifications) clearAll() {
	messages := n.messages
	for _, m := range messages {
		n.remove(m.widget())
	}
	n.messages = nil
}

// clearMessagesByType MUST be called from the ui thread
func (n *notifications) clearMessagesByType(mt gtki.MessageType) {
	messages := n.messages
	for i, m := range messages {
		if m.messageType() == mt {
			n.remove(m.widget())
			messages = append(messages[:i], messages[i+1:]...)
		}
	}
	n.messages = messages
}

// notify MUST be called from the UI thread
func (n *notifications) notify(text string, mt gtki.MessageType) {
	n.add(n.u.newNotificationBar(text, mt))
}

// warning MUST be called from the UI thread
func (n *notifications) warning(text string) {
	n.notify(text, gtki.MESSAGE_WARNING)
}

// error MUST be called from the UI thread
func (n *notifications) error(text string) {
	n.notify(text, gtki.MESSAGE_ERROR)
}

// info MUST be called from the ui thread
func (n *notifications) info(text string) {
	n.notify(text, gtki.MESSAGE_INFO)
}

// question MUST be called from the ui thread
func (n *notifications) question(text string) {
	n.notify(text, gtki.MESSAGE_QUESTION)
}

// message MUST be called from the ui thread
func (n *notifications) message(text string) {
	n.notify(text, gtki.MESSAGE_OTHER)
}

// notifyOnError is an alias for the "error" method and also
// implements the "canNotifyErrors" interface
func (n *notifications) notifyOnError(err string) {
	n.error(err)
}

// clearErrors is an alias for the "clear" method and also
// implements the "canNotifyErrors" interface
func (n *notifications) clearErrors() {
	n.clearMessagesByType(gtki.MESSAGE_ERROR)
}

type notificationBar struct {
	*infoBar
}

func (u *gtkUI) newNotificationBar(text string, mt gtki.MessageType) notification {
	return &notificationBar{
		u.newInfoBarComponent(text, mt),
	}
}

package gui

import (
	"sync"
	"time"

	"github.com/coyim/coyim/coylog"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

const mergeNotificationsThreshold = 7

func (u *gtkUI) lastActionTimeFor(f string) time.Time {
	return u.actionTimes[f]
}

func (u *gtkUI) registerLastActionTimeFor(f string, t time.Time) {
	u.actionTimes[f] = t
}

func (u *gtkUI) maybeNotify(timestamp time.Time, account *account, peer jid.WithoutResource, message string) {
	if u.deNotify == nil {
		return
	}

	dname := u.displayNameFor(account, peer)

	if timestamp.Before(u.lastActionTimeFor(peer.String()).Add(time.Duration(mergeNotificationsThreshold) * time.Second)) {
		u.log.Debug("Decided to not show notification, since the time is not ready")
		return
	}

	u.registerLastActionTimeFor(peer.String(), timestamp)

	err := u.deNotify.show(peer.String(), dname, message)
	if err != nil {
		u.log.WithError(err).Warn("Error when showing notification")
	}
}

func (u *gtkUI) showConnectAccountNotification(account *account) func() {
	var notification gtki.InfoBar

	doInUIThread(func() {
		notification = account.buildConnectionNotification()
		account.setCurrentNotification(notification, u.notificationArea)
	})

	return func() {
		doInUIThread(func() {
			account.removeCurrentNotificationIf(notification)
		})
	}
}

func (u *gtkUI) notifyTorIsNotRunning(account *account, moreInfo func()) {
	doInUIThread(func() {
		notification := account.buildTorNotRunningNotification(moreInfo)
		account.setCurrentNotification(notification, u.notificationArea)
	})
}

func (u *gtkUI) notifyConnectionFailure(account *account, moreInfo func()) {
	doInUIThread(func() {
		notification := account.buildConnectionFailureNotification(moreInfo)
		account.setCurrentNotification(notification, u.notificationArea)
	})
}

func (u *gtkUI) notify(title, message string) {
	builder := newBuilder("SimpleNotification")
	obj := builder.getObj("dialog")
	dlg := obj.(gtki.MessageDialog)

	_ = dlg.SetProperty("title", title)
	_ = dlg.SetProperty("text", message)
	dlg.SetTransientFor(u.window)

	doInUIThread(func() {
		dlg.Run()
		dlg.Destroy()
	})
}

type button struct {
	text         string
	responseType gtki.ResponseType
}

type widget interface {
	getWidget() gtki.Widget
}

type infoMessage interface {
	getMessageType() gtki.MessageType
}

type notificationWidget interface {
	widget
	infoMessage
}

type notifications struct {
	box      gtki.Box
	messages []notificationWidget
	options  map[string]interface{}

	lock sync.Mutex
	log  coylog.Logger

	showAllMessagesInStack bool
}

func (u *gtkUI) newNotifications(box gtki.Box) *notifications {
	n := &notifications{
		box: box,
		log: u.log.WithField("where", "notifications"),
	}

	return n
}

func (n *notifications) add(m notificationWidget) {
	if !n.showAllMessagesInStack {
		n.clearAll()
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	n.messages = append(n.messages, m)

	n.box.Add(m.getWidget())
	n.box.ShowAll()
}

func (n *notifications) remove(w gtki.Widget) {
	n.box.Remove(w)
}

func (n *notifications) clearAll() {
	n.lock.Lock()
	messages := n.messages
	n.messages = nil
	n.lock.Unlock()

	for _, m := range messages {
		n.remove(m.getWidget())
	}
}

func (n *notifications) clearMessageType(mt gtki.MessageType) {
	n.lock.Lock()
	messages := n.messages
	n.lock.Unlock()

	for i, m := range messages {
		if m.getMessageType() == mt {
			n.remove(m.getWidget())
			messages = append(messages[:i], messages[i+1:]...)
		}
	}

	n.lock.Lock()
	n.messages = messages
	n.lock.Unlock()
}

func (n *notifications) notify(text string, mt gtki.MessageType, b *button) {
	message := newInfoBar(text, mt)
	n.add(message)
}

func (n *notifications) warning(text string) {
	n.notify(text, gtki.MESSAGE_WARNING, nil)
}

func (n *notifications) error(text string) {
	n.notify(text, gtki.MESSAGE_ERROR, nil)
}

// notifyOnError is an alias for the "error" method and also
// implements the "canNotifyErrors" interface
func (n *notifications) notifyOnError(err string) {
	n.error(err)
}

// clearErrors is an alias for the "clear" method and also
// implements the "canNotifyErrors" interface
func (n *notifications) clearErrors() {
	n.clearMessageType(gtki.MESSAGE_ERROR)
}

func (n *notifications) info(text string) {
	n.notify(text, gtki.MESSAGE_INFO, nil)
}

func (n *notifications) question(text string) {
	n.notify(text, gtki.MESSAGE_QUESTION, nil)
}

func (n *notifications) message(text string) {
	n.notify(text, gtki.MESSAGE_OTHER, nil)
}

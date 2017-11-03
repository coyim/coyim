package gui

import (
	"fmt"
	"log"
	"time"

	"github.com/coyim/gotk3adapter/gtki"
)

const mergeNotificationsThreshold = 7

func (u *gtkUI) lastActionTimeFor(f string) time.Time {
	return u.actionTimes[f]
}

func (u *gtkUI) registerLastActionTimeFor(f string, t time.Time) {
	u.actionTimes[f] = t
}

func (u *gtkUI) maybeNotify(timestamp time.Time, account *account, from, message string) {
	if u.deNotify == nil {
		return
	}

	dname := u.displayNameFor(account, from)

	if timestamp.Before(u.lastActionTimeFor(from).Add(time.Duration(mergeNotificationsThreshold) * time.Second)) {
		fmt.Println("Decided to not show notification, since the time is not ready")
		return
	}

	u.registerLastActionTimeFor(from, timestamp)

	err := u.deNotify.show(from, dname, message)
	if err != nil {
		log.Println(err)
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

	dlg.SetProperty("title", title)
	dlg.SetProperty("text", message)
	dlg.SetTransientFor(u.window)

	doInUIThread(func() {
		dlg.Run()
		dlg.Destroy()
	})
}

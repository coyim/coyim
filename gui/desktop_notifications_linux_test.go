// +build linux

package gui

import . "gopkg.in/check.v1"

type desktopNotificationsSuite struct{}

var _ = Suite(&desktopNotificationsSuite{})

// func (s *desktopNotificationsSuite) Test_show_doesShowMessage(c *C) {
// 	dn := newDesktopNotifications()
// 	showMessage := true

// 	err := dn.show("someone@coy.im", "Some one", "Hi!", showMessage, true)
// 	notify.CloseNotification(dn.notifications["someone@coy.im"])

// 	c.Assert(err, Equals, nil)
// 	c.Assert(dn.notification.Summary, Equals, "From: Some one")
// 	c.Assert(dn.notification.Body, Equals, "Hi!")
// }

// func (s *desktopNotificationsSuite) Test_show_doesNotShowMessageWhenShowIsFalse(c *C) {
// 	dn := newDesktopNotifications()
// 	showMessage := false

// 	err := dn.show("someone@coy.im", "Some one", "Hi!", showMessage, true)
// 	notify.CloseNotification(dn.notifications["someone@coy.im"])

// 	c.Assert(err, Equals, nil)
// 	c.Assert(dn.notification.Summary, Equals, "New message!")
// 	c.Assert(dn.notification.Body, Equals, "From: <b>Some one</b>")
// }

// func (s *desktopNotificationsSuite) Test_show_doesNotShowMessageWhenMessageIsEmpty(c *C) {
// 	dn := newDesktopNotifications()
// 	showMessage := true

// 	err := dn.show("someone@coy.im", "Some one", "", showMessage, true)
// 	notify.CloseNotification(dn.notifications["someone@coy.im"])

// 	c.Assert(err, Equals, nil)
// 	c.Assert(dn.notification.Summary, Equals, "New message!")
// 	c.Assert(dn.notification.Body, Equals, "From: <b>Some one</b>")
// }

// func (s *desktopNotificationsSuite) Test_show_queuesSeveralNotifications(c *C) {
// 	dn := newDesktopNotifications()
// 	showMessage := true
// 	fullScreen := true

// 	err := dn.show("one@coy.im", "Some one", "One!", showMessage, fullScreen)
// 	c.Assert(err, Equals, nil)
// 	c.Assert(dn.notification.Summary, Equals, "From: Some one")
// 	c.Assert(dn.notification.Body, Equals, "One!")
// 	c.Assert(dn.notification.Hints[notify.HintUrgency], Equals, notify.UrgencyCritical)

// 	fullScreen = false
// 	err = dn.show("two@coy.im", "Some two", "Two!", showMessage, fullScreen)
// 	c.Assert(err, Equals, nil)
// 	c.Assert(dn.notification.Summary, Equals, "From: Some two")
// 	c.Assert(dn.notification.Body, Equals, "Two!")
// 	c.Assert(dn.notification.Hints[notify.HintUrgency], Equals, nil)

// 	notify.CloseNotification(dn.notifications["one@coy.im"])
// 	notify.CloseNotification(dn.notifications["two@coy.im"])
// }

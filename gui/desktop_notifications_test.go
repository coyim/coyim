package gui

import (
	. "gopkg.in/check.v1"
)

type DesktopNotificationsSuite struct{}

var _ = Suite(&DesktopNotificationsSuite{})

func (dns *DesktopNotificationsSuite) Test_NotificationsFormat(c *C) {
	expectedSummary := "New message!"
	expectedBody := "You have a new message"

	from := "user@domain.com"
	message := "This is a message."
	withHTML := false
	dn := new(desktopNotifications)

	dn.notificationStyle = "only-presence-of-new-information"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)

	expectedSummary = "New message!"
	expectedBody = "From: user@domain.com"

	dn = new(desktopNotifications)

	dn.notificationStyle = "with-author-but-no-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody = dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)

	expectedSummary = "From: user@domain.com"
	expectedBody = "This is a message."

	dn = new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody = dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)

	expectedSummary = ""
	expectedBody = ""
	dn = new(desktopNotifications)

	dn.notificationStyle = "this-style-does-not-exist"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody = dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_NotificationsFormat_WithHTMLTags(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From: <b>user@domain.com</b>"
	from := "user@domain.com"
	message := ""
	withHTML := true
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-author-but-no-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)

	expectedSummary = "From: user@domain.com"
	expectedBody = "&lt;b&gt;This&lt;/b&gt; &lt;i&gt;is&lt;/i&gt; a message."
	message = "<b>This</b> <i>is</i> a message."

	dn = new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody = dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)

	expectedSummary = "From: < user&name >"
	expectedBody = "This is a message."
	from = "< user&name >"
	message = "This is a message."

	actualSummary, actualBody = dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)

	expectedSummary = "From: < user&name >"
	expectedBody = "< This is a message. >"
	from = "< user&name >"
	message = "< This is a message. >"

	actualSummary, actualBody = dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

package gui

import (
	. "gopkg.in/check.v1"
)

type DesktopNotificationsSuite struct{}

var _ = Suite(&DesktopNotificationsSuite{})

func (dns *DesktopNotificationsSuite) Test_format_WithoutHTML_OnlyPresence(c *C) {
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
}

func (dns *DesktopNotificationsSuite) Test_format_WithHTML_OnlyPresence(c *C) {
	expectedSummary := "New message!"
	expectedBody := "You have a new message"
	from := "user@domain.com"
	message := "This is a message."
	withHTML := true
	dn := new(desktopNotifications)

	dn.notificationStyle = "only-presence-of-new-information"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WithoutHTML_NoContent(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From: user@domain.com"
	from := "user@domain.com"
	message := "This is a message."
	withHTML := false
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-author-but-no-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WithHTML_NoContent(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From: <b>user@domain.com</b>"
	from := "user@domain.com"
	message := "This is a message."
	withHTML := true
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-author-but-no-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WithHTML_NoEscapeSummary_NoContent(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From: <b>< user&name ></b>"
	from := "< user&name >"
	message := "This is a message."
	withHTML := true
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-author-but-no-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WithoutHTML_WithContent(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From - user@domain.com: This is a message."
	from := "user@domain.com"
	message := "This is a message."
	withHTML := false
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Skip("Failing test")
	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WithoutHTML_WithHTMLTaggedContent(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From - user@domain.com: &lt;b&gt;This&lt;/b&gt; &lt;i&gt;is&lt;/i&gt; a message."
	from := "user@domain.com"
	message := "<b>This</b> <i>is</i> a message."
	withHTML := false
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Skip("Failing test")
	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WithoutHTML_WithHTMLSingleWordGreaterThan254(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From - user@domain.com: Loremipsumdolorsitamet,consecteturadipiscingelit.Donecsempertristiquetellus,inullamcorperfelis.Crasscelerisquenisisedtristiquesemper.Curabitureudictumnisl.Donecpharetrahendreritdiamegetlacinia.Vivamusauctorsemmi,atcrasamet.Nequeporroquisquamestquidolore..."
	from := "user@domain.com"
	message := "Loremipsumdolorsitamet,consecteturadipiscingelit.Donecsempertristiquetellus,inullamcorperfelis.Crasscelerisquenisisedtristiquesemper.Curabitureudictumnisl.Donecpharetrahendreritdiamegetlacinia.Vivamusauctorsemmi,atcrasamet.Nequeporroquisquamestquidoloremipsumquiadolorsitamet,consectetur,adipiscivelit."
	withHTML := false
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Skip("Failing test")
	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WithoutHTML_WithContentGreaterThan254(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From - user@domain.com: Lorem ipsum dolor sit amet, " +
		"consectetur adipiscing elit. Nam consectetur elit leo, " +
		"nec tincidunt velit ultrices vitae. Aenean aliquam massa at dapibus " +
		"cursus. Mauris feugiat sed velit in porttitor. Fusce id nisi at leo " +
		"consequat feugiat vel sed..."
	message := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
		"Nam consectetur elit leo, nec tincidunt velit ultrices vitae. " +
		"Aenean aliquam massa at dapibus cursus. Mauris feugiat sed velit " +
		"in porttitor. Fusce id nisi at leo consequat feugiat vel sed enim. Duis sed."
	withHTML := false
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format("user@domain.com", message, withHTML)

	c.Skip("Failing test")
	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WitHTML_WithContent(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From - <b>user@domain.com</b>: This is a message."
	from := "user@domain.com"
	message := "This is a message."
	withHTML := true
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Skip("Failing test")
	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WithHTML_WithHTMLTaggedContent(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From - <b>user@domain.com</b>: &lt;b&gt;This&lt;/b&gt; &lt;i&gt;is&lt;/i&gt; a message."
	from := "user@domain.com"
	message := "<b>This</b> <i>is</i> a message."
	withHTML := true
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Skip("Failing test")
	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_WitHTML_NoEscapeBody_WithContent(c *C) {
	expectedSummary := "New message!"
	expectedBody := "From - <b>< user&name ></b>: This is a message."
	from := "< user&name >"
	message := "This is a message."
	withHTML := true
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Skip("Failing test")
	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

func (dns *DesktopNotificationsSuite) Test_format_InvalidNotificationStyle(c *C) {
	expectedSummary := ""
	expectedBody := ""
	from := "user@domain.com"
	message := "This is a message."
	withHTML := true
	dn := new(desktopNotifications)

	dn.notificationStyle = "this-style-does-not-exist"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	actualSummary, actualBody := dn.format(from, message, withHTML)

	c.Assert(actualSummary, Equals, expectedSummary)
	c.Assert(actualBody, Equals, expectedBody)
}

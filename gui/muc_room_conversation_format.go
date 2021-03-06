package gui

import (
	"time"

	"github.com/coyim/coyim/xmpp/jid"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

// displayCurrentTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayCurrentTimestamp() {
	c.displayTimestamp(time.Now())
}

// displayTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayTimestamp(timestamp time.Time) {
	c.addTextWithTag(i18n.Localf("[%s] ", formatTimestamp(timestamp)), conversationTagTimestamp)
}

// displayNotificationWhenOccupantJoinedRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	c.displayTextLineWithTimestamp(i18n.Localf("%s joined the room", nickname), conversationTagSomeoneJoinedRoom)
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	c.displayTextLineWithTimestamp(i18n.Localf("%s left the room", nickname), conversationTagSomeoneLeftRoom)
}

// displayNickname MUST be called from the UI thread
func (c *roomViewConversation) displayNickname(nickname string) {
	c.addTextWithTag(i18n.Localf("%s: ", nickname), conversationTagNickname)
}

// displayRoomSubject MUST be called from the UI thread
func (c *roomViewConversation) displayRoomSubject(subject string) {
	c.displayTextLineWithTimestamp(subject, conversationTagRoomSubject)
}

// displayMessage MUST be called from the UI thread
func (c *roomViewConversation) displayMessage(message string) {
	c.addTextWithTag(message, conversationTagMessage)
}

// displayInfoMessage MUST be called from the UI thread
func (c *roomViewConversation) displayInfoMessage(message string) {
	c.addTextWithTag(message, conversationTagInfo)
}

// displayConfigurationMessage MUST be called from the UI thread
func (c *roomViewConversation) displayConfigurationMessage(message string) {
	c.addTextWithTag(message, conversationTagRoomConfigChange)
}

// displayLiveMessage MUST be called from the UI thread
func (c *roomViewConversation) displayLiveMessage(nickname, message string, timestamp time.Time) {
	c.displayTimestamp(timestamp)

	c.displayNickname(nickname)
	c.displayMessage(message)

	c.addNewLine()
}

// displayDiscussionHistoryDate MUST be called from the UI thread
func (c *roomViewConversation) displayDiscussionHistoryDate(d time.Time) {
	j := c.chatTextView.GetJustification()
	defer c.chatTextView.SetJustification(j)

	c.chatTextView.SetJustification(gtki.JUSTIFY_CENTER)

	c.addTextWithTag(timeToFriendlyString(d), conversationTagDateGroup)
	c.addNewLine()
}

// displayDiscussionHistoryMessages MUST be called from the UI thread
func (c *roomViewConversation) displayDiscussionHistoryMessages(messages []*data.DelayedMessage) {
	for _, m := range messages {
		c.displayDelayedMessage(m.Nickname, m.Message, m.Timestamp)
	}
}

// displayDelayedMessage MUST be called from the UI thread
func (c *roomViewConversation) displayDelayedMessage(nickname, message string, timestamp time.Time) {
	c.displayTimestamp(timestamp)

	c.displayNickname(nickname)
	c.displayMessage(message)

	c.addNewLine()
}

// displayNewInfoMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewInfoMessage(message string) {
	c.displayCurrentTimestamp()
	c.displayInfoMessage(message)
	c.addNewLine()
}

// displayNewConfigurationMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewConfigurationMessage(message string) {
	c.displayCurrentTimestamp()
	c.displayConfigurationMessage(message)
	c.addNewLine()
}

// displayWarningMessage MUST be called from the UI thread
func (c *roomViewConversation) displayWarningMessage(message string) {
	c.displayTextLineWithTimestamp(message, conversationTagWarning)
}

// displayErrorMessage MUST be called from the UI thread
func (c *roomViewConversation) displayErrorMessage(message string) {
	c.displayTextLineWithTimestamp(message, conversationTagError)
}

const conversationDividerText = "――――――――――――――――――――"

// displayDivider MUST be called from the UI thread
func (c *roomViewConversation) displayDivider() {
	c.addTextWithTag(conversationDividerText, conversationTagDivider)
	c.addNewLine()
}

// displayTextLineWithTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayTextLineWithTimestamp(text string, tag conversationTag) {
	c.displayCurrentTimestamp()
	c.addTextWithTag(text, tag)
	c.addNewLine()
}

// displayNotificationWhenRoomDestroyed MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenRoomDestroyed(reason string, alternative jid.Bare, password string) {
	c.displayCurrentTimestamp()

	message := getDisplayForRoomDestroyed(&roomDestroyedData{reason, alternative, password})
	c.displayFormattedMessage(message, func(m string) {
		c.addTextWithTag(m, conversationTagWarning)
	})

	c.addNewLine()

	c.displayTextLineWithTimestamp(i18n.Local("You can no longer receive any messages in this room and the occupant list will not be updated anymore."), conversationTagWarning)
}

func formatTimestamp(t time.Time) string {
	return t.Format("15:04:05")
}

func formatTimeWithLayout(t time.Time, layout string) string {
	return t.Format(layout)
}

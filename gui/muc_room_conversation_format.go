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
	c.addTextWithTag(i18n.Localf("[%s] ", formatTimestamp(timestamp)), "timestamp")
}

// displayNotificationWhenOccupantJoinedRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	c.displayTextLineWithTimestamp(i18n.Localf("%s joined the room", nickname), "joinedRoom")
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	c.displayTextLineWithTimestamp(i18n.Localf("%s left the room", nickname), "leftRoom")
}

// displayNickname MUST be called from the UI thread
func (c *roomViewConversation) displayNickname(nickname string) {
	c.addTextWithTag(i18n.Localf("%s: ", nickname), "nickname")
}

// displayRoomSubject MUST be called from the UI thread
func (c *roomViewConversation) displayRoomSubject(subject string) {
	c.displayTextLineWithTimestamp(subject, "subject")
}

// displayMessage MUST be called from the UI thread
func (c *roomViewConversation) displayMessage(message string) {
	c.addTextWithTag(message, "message")
}

// displayInfoMessage MUST be called from the UI thread
func (c *roomViewConversation) displayInfoMessage(message string) {
	c.addTextWithTag(message, "infoMessage")
}

// displayConfigurationMessage MUST be called from the UI thread
func (c *roomViewConversation) displayConfigurationMessage(message string) {
	c.addTextWithTag(message, "configuration")
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

	c.addTextWithTag(timeToFriendlyString(d), "groupdate")
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
	c.displayTextLineWithTimestamp(message, "warning")
}

// displayErrorMessage MUST be called from the UI thread
func (c *roomViewConversation) displayErrorMessage(message string) {
	c.displayTextLineWithTimestamp(message, "error")
}

// displayDivider MUST be called from the UI thread
func (c *roomViewConversation) displayDivider() {
	c.addTextWithTag("――――――――――――――――――――", "divider")
	c.addNewLine()
}

// displayTextLineWithTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayTextLineWithTimestamp(text string, tag string) {
	c.displayCurrentTimestamp()
	c.addTextWithTag(text, tag)
	c.addNewLine()
}

// displayNotificationWhenRoomDestroyed MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenRoomDestroyed(reason string, alternative jid.Bare, password string) {
	c.displayAlternativeRoomInfo(reason, alternative, password)
	c.displayTextLineWithTimestamp(i18n.Local("You can no longer receive any messages in this room and the occupant list will not be updated anymore."), "warning")
}

func (c *roomViewConversation) displayAlternativeRoomInfoWithPassword(message, password string) {
	c.displayCurrentTimestamp()
	c.addTextWithTag(i18n.Localf("%s , with this password: \"", message), "warning")
	c.addTextWithTag(i18n.Localf("%s", password), "password")
	c.addTextWithTag(i18n.Local("\"."), "warning")
	c.addNewLine()
}

func (c *roomViewConversation) displayAlternativeRoomInfo(reason string, alternative jid.Bare, password string) {
	message := i18n.Local("The room was destroyed")

	if reason != "" {
		message = i18n.Localf("%s. The reason given was \"%s\"", message, reason)
	}

	if alternative != nil {
		message = i18n.Localf("%s. Discussions will continue in this room: \"%s\"", message, alternative)
		if password != "" {
			c.displayAlternativeRoomInfoWithPassword(message, password)
			return
		}
	}
		c.displayTextLineWithTimestamp(i18n.Localf("%s.", message), "warning")
	}
}

func formatTimestamp(t time.Time) string {
	return t.Format("15:04:05")
}

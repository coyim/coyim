package gui

import (
	"time"

	"github.com/coyim/coyim/i18n"
)

// displayTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayTimestamp() {
	c.addTextWithTag(i18n.Localf("[%s] ", getTimestamp()), "timestamp")
}

// displayNotificationWhenOccupantJoinedRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	c.addTextLineWithTimestamp(i18n.Localf("%s joined the room", nickname), "joinedRoom")
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	c.addTextLineWithTimestamp(i18n.Localf("%s left the room", nickname), "leftRoom")
}

// displayNickname MUST be called from the UI thread
func (c *roomViewConversation) displayNickname(nickname string) {
	c.addTextWithTag(i18n.Localf("%s: ", nickname), "nickname")
}

// displayRoomSubject MUST be called from the UI thread
func (c *roomViewConversation) displayRoomSubject(subject string) {
	c.displayTimestamp()
	c.addTextWithTag(subject, "subject")
	c.addNewLine()
}

// displayMessage MUST be called from the UI thread
func (c *roomViewConversation) displayMessage(message string) {
	c.addTextWithTag(message, "message")
}

// displayNewLiveMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewLiveMessage(nickname, subject, message string) {
	c.displayTimestamp()

	c.displayNickname(nickname)
	if subject != "" {
		c.displayRoomSubject(subject)
	}
	c.displayMessage(message)

	c.addNewLine()
}

// displayWarningMessage MUST be called from the UI thread
func (c *roomViewConversation) displayWarningMessage(message string) {
	c.addTextLineWithTimestamp(message, "warning")
}

// displayErrorMessage MUST be called from the UI thread
func (c *roomViewConversation) displayErrorMessage(message string) {
	c.addTextLineWithTimestamp(message, "error")
}

func getTimestamp() string {
	return time.Now().Format("15:04:05")
}

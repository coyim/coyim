package gui

import (
	"time"

	"github.com/coyim/coyim/text"
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
	c.addTextWithTag(messageForTimestamp(timestamp), conversationTagTimestamp)
}

// saveAndDisplayMessage MUST be called from the UI thread
func (c *roomViewConversation) saveAndDisplayMessage(nickname, message string, timestamp time.Time, messageType data.MessageType) {
	dm := data.NewDelayedMessage(nickname, message, timestamp, messageType)
	c.saveNotificationMessage(dm)
	c.displayMessageFromData(dm)
}

// handleOccupantJoinedRoom MUST be called from the UI thread
func (c *roomViewConversation) handleOccupantJoinedRoom(nickname string) {
	c.saveAndDisplayMessage(nickname, messageForSomeoneWhoJoinedTheRoom(nickname), time.Now(), data.Joined)
}

// displayMessageFromData MUST be called from the UI thread
func (c *roomViewConversation) displayMessageFromData(dm *data.DelayedMessage) {
	c.addNewLine()
	c.displayTimestamp(dm.Timestamp)
	c.displayFormattedMessage(dm.Message, messageTagBasedOnMessageType[dm.MessageType])
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	message := messageForSomeoneWhoLeftTheRoom(nickname)
	c.displayFormattedMessageWithTimestamp(message, conversationTagSomeoneLeftRoom)
}

// displayNickname MUST be called from the UI thread
func (c *roomViewConversation) displayNickname(nickname string) {
	c.addTextWithTag(messageForNickname(nickname), conversationTagNickname)
}

// displayRoomSubject MUST be called from the UI thread
func (c *roomViewConversation) displayRoomSubject(subject string) {
	c.displayTextLineWithTimestamp(subject, conversationTagRoomSubject)
}

// handleRoomSubject MUST be called from the UI thread
func (c *roomViewConversation) handleRoomSubject(subject string) {
	c.saveAndDisplayMessage("", subject, time.Now(), data.Subject)
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

// displayNewMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewMessage(nickname, message string, timestamp time.Time) {
	c.addNewLine()
	c.displayTimestamp(timestamp)
	c.displayNickname(nickname)
	c.displayMessage(message)
}

// displayDiscussionHistoryDate MUST be called from the UI thread
func (c *roomViewConversation) displayDiscussionHistoryDate(d time.Time) {
	c.addNewLine()

	j := c.chatTextView.GetJustification()
	defer c.chatTextView.SetJustification(j)

	c.chatTextView.SetJustification(gtki.JUSTIFY_CENTER)

	c.addTextWithTag(timeToFriendlyDate(d), conversationTagDateGroup)
}

// displayDiscussionHistoryMessages MUST be called from the UI thread
func (c *roomViewConversation) displayDiscussionHistoryMessages(messages []*data.DelayedMessage) {
	for _, m := range messages {
		if m.MessageType == data.Chat {
			c.displayNewMessage(m.Nickname, m.Message, m.Timestamp)
		} else {
			c.displayMessageFromData(m)
		}
	}
}

// displayNewInfoMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewInfoMessage(message string) {
	c.addNewLine()
	c.displayCurrentTimestamp()
	c.displayInfoMessage(message)
}

// displayNewConfigurationMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewConfigurationMessage(message string) {
	c.addNewLine()
	c.displayCurrentTimestamp()
	c.displayConfigurationMessage(message)
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
	c.addNewLine()
	c.addTextWithTag(conversationDividerText, conversationTagDivider)
}

// displayTextLineWithTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayTextLineWithTimestamp(text string, tag conversationTag) {
	c.addNewLine()
	c.displayCurrentTimestamp()
	c.addTextWithTag(text, tag)
}

// displayNotificationWhenRoomDestroyed MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenRoomDestroyed(reason string, alternative jid.Bare, password string) {
	c.saveAndDisplayMessage("", messageForRoomDestroyed(&roomDestroyedData{reason, alternative, password}), time.Now(), data.Warning)

	c.displayTextLineWithTimestamp(i18n.Local("You can no longer receive any messages in this room and the list of people in the room will not be updated anymore."), conversationTagWarning)
}

// displayOccupantUpdateMessageFor MUST be called from the UI thread
func (c *roomViewConversation) displayOccupantUpdateMessageFor(update interface{}) {
	c.saveAndDisplayMessage("", getMUCNotificationMessageFrom(update), time.Now(), data.OccupantInformationChanged)
}

// displayFormattedMessageWithTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayFormattedMessageWithTimestamp(message string, noFormattedTextTag conversationTag) {
	c.addNewLine()
	c.displayCurrentTimestamp()
	c.displayFormattedMessage(message, noFormattedTextTag)
}

// displayFormattedMessage MUST be called from the UI thread
func (c *roomViewConversation) displayFormattedMessage(message string, noFormattedTextTag conversationTag) {
	displayMessage := func(text string) {
		c.addTextWithTag(text, noFormattedTextTag)
	}

	if formatted, ok := text.ParseWithFormat(message); ok {
		text, formats := formatted.Join()

		lastDisplayedIndex := 0
		for _, format := range formats {
			previousTextBeforeFormat := text[lastDisplayedIndex:format.Start]
			displayMessage(previousTextBeforeFormat)

			textFormatSize := format.Start + format.Length
			textFormat := text[format.Start:textFormatSize]
			c.displayMessageFormatting(textFormat, format, noFormattedTextTag)

			lastDisplayedIndex = textFormatSize
		}

		restOfTheText := text[lastDisplayedIndex:]
		if restOfTheText != "" {
			displayMessage(restOfTheText)
		}
	} else {
		displayMessage(message)
	}
}

// displayMessageFormatting MUST be called from the UI thread
func (c *roomViewConversation) displayMessageFormatting(message string, format text.Formatting, tag conversationTag) {
	c.addTextWithTag(message, formattingTagName(format.Format, tag))
}

func formatTimestamp(t time.Time) string {
	return t.Format("15:04:05")
}

func formatTimeWithLayout(t time.Time, layout string) string {
	return t.Format(layout)
}

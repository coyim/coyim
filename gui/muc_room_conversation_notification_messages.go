package gui

import (
	"time"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
)

type roomConversationMessageType int

const (
	noMessagesBecauseRoomDestroyed roomConversationMessageType = iota
	noMessagesBecauseSelfOccupantBanned
	noMessagesBecauseSelfOccupantExpelled
	noMessagesBecauseSelfOccupantVisitor
	noMessagesBecauseForbidden
	noMessagesBecauseNotAcceptable
	noMessagesBecauseRoomConfiguration

	roomPubliclyLogged
	roomNotPubliclyLogged
	roomCanSeeRealJid
	roomCanNotSeeRealJid

	messageCouldNotBeSent
)

var roomConversationTexts map[roomConversationMessageType]string

func initMUCRoomConversationTexts() {
	roomConversationTexts = map[roomConversationMessageType]string{
		noMessagesBecauseRoomDestroyed:        i18n.Local("You can't send messages because this room has been destroyed."),
		noMessagesBecauseSelfOccupantBanned:   i18n.Local("You can't send messages because you have been banned."),
		noMessagesBecauseSelfOccupantExpelled: i18n.Local("You can't send messages because you were expelled from the room."),
		noMessagesBecauseSelfOccupantVisitor:  i18n.Local("As a visitor, you can't send messages in a moderated room."),
		noMessagesBecauseForbidden:            i18n.Local("You are forbidden to send messages to this room."),
		noMessagesBecauseNotAcceptable:        i18n.Local("Your messages to this room aren't accepted."),
		noMessagesBecauseRoomConfiguration:    i18n.Local("You can't send messages because the room configuration has been changed."),
		roomPubliclyLogged:                    i18n.Local("This room is now publicly logged."),
		roomNotPubliclyLogged:                 i18n.Local("This room is not publicly logged anymore."),
		roomCanSeeRealJid:                     i18n.Local("Your real JID can now be seen by anyone."),
		roomCanNotSeeRealJid:                  i18n.Local("Your real JID can now be seen only by moderators."),
		messageCouldNotBeSent:                 i18n.Local("The message couldn't be sent, please try again"),
	}
}

func messageForType(tp roomConversationMessageType) string {
	if message, ok := roomConversationTexts[tp]; ok {
		return message
	}

	panic("developer error: retrieving message for an unknown conversation event type")
}

func messageForRoomDestroyedEvent() string {
	return messageForType(noMessagesBecauseRoomDestroyed)
}

func messageForSelfOccupantBanned() string {
	return messageForType(noMessagesBecauseSelfOccupantBanned)
}

func messageForSelfOccupantExpelled() string {
	return messageForType(noMessagesBecauseSelfOccupantExpelled)
}

func messageForSelfOccupantVisitor() string {
	return messageForType(noMessagesBecauseSelfOccupantVisitor)
}

func messageForSelfOccupantForbidden() string {
	return messageForType(noMessagesBecauseForbidden)
}

func messageForSelfOccupantNotAcceptable() string {
	return messageForType(noMessagesBecauseNotAcceptable)
}

func messageForSelfOccupantRoomConfiguration() string {
	return messageForType(noMessagesBecauseRoomConfiguration)
}

func messageForRoomPubliclyLogged() string {
	return messageForType(roomPubliclyLogged)
}

func messageForRoomNotPubliclyLogged() string {
	return messageForType(roomNotPubliclyLogged)
}

func messageForRoomCanSeeRealJid() string {
	return messageForType(roomCanSeeRealJid)
}

func messageForRoomCanNotSeeRealJid() string {
	return messageForType(roomCanNotSeeRealJid)
}

func messageNotSent() string {
	return messageForType(messageCouldNotBeSent)
}

func messageForTimestamp(timestamp time.Time) string {
	return i18n.Localf("[%s] ", formatTimestamp(timestamp))
}

func messageForNickname(nickname string) string {
	return i18n.Localf("%s: ", nickname)
}

func messageForSomeoneWhoJoinedTheRoom(nickname string) string {
	return i18n.Localf("$nickname{%s} joined the room", nickname)
}

func messageForSomeoneWhoLeftTheRoom(nickname string) string {
	return i18n.Localf("$nickname{%s} left the room")
}

func messageForRoomSubjectUpdate(nickname, subject string) string {
	if nickname == "" {
		return i18n.Localf("Someone has updated the room subject to: \"%s\"", subject)
	}

	return i18n.Localf("$nickname{%[1]s} updated the room subject to: \"%[2]s\"", nickname, subject)
}

func messageForRoomSubject(subject string) string {
	if subject == "" {
		return i18n.Local("The room does not have a subject")
	}

	return i18n.Localf("The room subject is \"%s\"", subject)
}

type roomDestroyedData struct {
	reason      string
	alternative jid.Bare
	password    string
}

func (d *roomDestroyedData) hasReason() bool {
	return d.reason != ""
}

func (d *roomDestroyedData) hasAlternativeRoom() bool {
	return d.alternative != nil
}

func (d *roomDestroyedData) hasPassword() bool {
	return d.password != ""
}

func messageForRoomDestroyed(d *roomDestroyedData) string {
	switch {
	case d.hasReason() && d.hasAlternativeRoom() && d.hasPassword():
		return i18n.Localf("The room was destroyed. The reason given was \"%[1]s\". "+
			"Discussions will continue in this room: %[2]s, with this password: \"$password{%[3]s}\".",
			d.reason, d.alternative, d.password)

	case d.hasReason() && d.hasAlternativeRoom():
		return i18n.Localf("The room was destroyed. The reason given was \"%[1]s\". "+
			"Discussions will continue in this room: %[2]s.",
			d.reason, d.alternative)

	case d.hasReason():
		return i18n.Localf("The room was destroyed. The reason given was \"%s\".", d.reason)
	}

	return i18n.Local("The room was destroyed.")
}

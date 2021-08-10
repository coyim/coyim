package gui

import (
	"time"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
)

func messageForRoomDestroyedEvent() string {
	return i18n.Local("You can't send messages because this room has been destroyed.")
}

func messageForSelfOccupantBanned() string {
	return i18n.Local("You can't send messages because you have been banned.")
}

func messageForSelfOccupantExpelled() string {
	return i18n.Local("You can't send messages because you were expelled from the room.")
}

func messageForSelfOccupantVisitor() string {
	return i18n.Local("As a visitor, you can't send messages in a moderated room.")
}

func messageForSelfOccupantForbidden() string {
	return i18n.Local("You are forbidden to send messages to this room.")
}

func messageForSelfOccupantNotAcceptable() string {
	return i18n.Local("Your messages to this room aren't accepted.")
}

func messageForSelfOccupantRoomConfiguration() string {
	return i18n.Local("You can't send messages because the room configuration has been changed.")
}

func messageForRoomPubliclyLogged() string {
	return i18n.Local("This room is now publicly logged.")
}

func messageForRoomNotPubliclyLogged() string {
	return i18n.Local("This room is not publicly logged anymore.")
}

func messageForRoomCanSeeRealJid() string {
	return i18n.Local("Your real JID can now be seen by anyone.")
}

func messageForRoomCanNotSeeRealJid() string {
	return i18n.Local("Your real JID can now be seen only by moderators.")
}

func messageNotSent() string {
	return i18n.Local("The message couldn't be sent, please try again")
}

func messageForSelfOccupantDisconnected() string {
	return i18n.Local("You can't send messages because you lost connection.")
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
	return i18n.Localf("$nickname{%s} left the room", nickname)
}

func messageForRoomSubjectUpdate(nickname, subject string) string {
	if subject == "" {
		return messageForRoomSubjectRemoved(nickname)
	}

	if nickname == "" {
		return i18n.Localf("Someone has updated the room subject to: \"%s\"", subject)
	}

	return i18n.Localf("$nickname{%[1]s} updated the room subject to: \"%[2]s\"", nickname, subject)
}

func messageForRoomSubjectRemoved(nickname string) string {
	if nickname == "" {
		return i18n.Local("Someone removed the room subject")
	}

	return i18n.Localf("$nickname{%[1]s} removed the room subject", nickname)
}

func messageForRoomSubject(subject string) string {
	if subject == "" {
		return i18n.Local("The room does not have a subject")
	}

	return i18n.Localf("The room subject is \"%s\"", subject)
}

func messageForMembersOnlyRoom(nickname string) string {
	return i18n.Localf("$nickname{%s} was removed from this room because it's now a members only room.", nickname)
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

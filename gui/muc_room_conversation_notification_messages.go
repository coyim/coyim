package gui

import (
	"time"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
)

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

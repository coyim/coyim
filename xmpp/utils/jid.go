package utils

import "strings"

// RemoveResourceFromJid returns the user@domain portion of a JID.
func RemoveResourceFromJid(jid string) string {
	slash := strings.Index(jid, "/")
	if slash != -1 {
		return jid[:slash]
	}
	return jid
}

// DomainFromJid returns the domain of a full or bare JID.
func DomainFromJid(jid string) string {
	jid = RemoveResourceFromJid(jid)
	at := strings.Index(jid, "@")
	if at != -1 {
		return jid[at+1:]
	}
	return jid
}

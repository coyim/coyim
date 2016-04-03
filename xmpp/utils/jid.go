package utils

import (
	"fmt"
	"strings"
)

// SplitJid will return the clean JID and the resource as separate strings
func SplitJid(jid string) (string, string) {
	slash := strings.Index(jid, "/")
	if slash != -1 {
		return jid[:slash], jid[slash+1:]
	}
	return jid, ""
}

// RemoveResourceFromJid returns the user@domain portion of a JID.
func RemoveResourceFromJid(jid string) string {
	j, _ := SplitJid(jid)
	return j
}

// ResourceFromJid returns the resource portion of a JID.
func ResourceFromJid(jid string) string {
	_, r := SplitJid(jid)
	return r
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

// ComposeFullJid puts together a jid from the given bare jid and resource
func ComposeFullJid(jid, resource string) string {
	if resource == "" {
		return jid
	}
	return fmt.Sprintf("%s/%s", jid, resource)
}

package gui

import (
	"strings"

	"github.com/coyim/coyim/i18n"
)

// Method to validate a jabber id is correct according to the RFC-6122
// on Address Format
func verifyXMPPAddress(address string) (bool, string) {
	tldIndex := strings.LastIndex(address, ".")
	atIndex := strings.Index(address, "@")

	if !isValidAccount(address, atIndex) {
		err := i18n.Local("Validation failed: \nThe XMPP address is invalid. An XMPP address should look like this: local@domain.com.")
		return false, err
	}

	isValidDomain := verifyDomainPart(tldIndex, atIndex, address)
	isValidPart := verifyLocalPart(atIndex)

	switch {
	case !isValidDomain && !isValidPart:
		return false, i18n.Local("Validation failed:\nThe XMPP address has an invalid domain part, The XMMP address has an invalid local part. An XMPP address should look like this: local@domain.com.")
	case !isValidDomain:
		return false, i18n.Local("Validation failed:\nThe XMPP address has an invalid domain part. An XMPP address should look like this: local@domain.com.")
	case !isValidPart:
		return false, i18n.Local("Validation failed:\nThe XMMP address has an invalid local part. An XMPP address should look like this: local@domain.com.")
	}

	return true, ""
}

func isValidAccount(address string, atIndex int) bool {
	return len(address) != 0 && atIndex != -1
}

func verifyDomainPart(tldIndex, atIndex int, address string) bool {
	return !(tldIndex < 0 || tldIndex == len(address)-1 || tldIndex-atIndex < 2)
}

func verifyLocalPart(atIndex int) bool {
	return atIndex != 0
}

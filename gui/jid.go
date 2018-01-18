package gui

import (
	"strings"
)

// Method to validate a jabber id is correct according to the RFC-6122
// on Address Format
func verifyXMPPAddress(address string) (bool, string) {
	var err string
	var isValid bool

	tldIndex := strings.LastIndex(address, ".")
	atIndex := strings.Index(address, "@")

	if !isValidAccount(address, atIndex) {
		err = "Validation failed: \nThe XMPP address is invalid. An XMPP address should look like this: local@domain.com."
		return isValid, err
	}

	isValidDomain, errDomain := verifyDomainPart(tldIndex, atIndex, address)
	isValidPart, errLocal := verifyLocalPart(atIndex)
	isValid = isValidDomain && isValidPart

	if !isValid {
		result := "Validation failed:\n"
		sep := ""
		if errDomain != "" {
			result += errDomain
			sep = ", "
		}
		if errLocal != "" {
			result += sep + errLocal
		}
		err = result + ". An XMPP address should look like this: local@domain.com."
	}

	return isValid, err
}

func isValidAccount(address string, atIndex int) bool {
	if len(address) != 0 && atIndex != -1 {
		return true
	}

	return false
}

func verifyDomainPart(tldIndex, atIndex int, address string) (bool, string) {
	isValid := true
	var err string

	if tldIndex < 0 || tldIndex == len(address)-1 || tldIndex-atIndex < 2 {
		isValid = false
		err = "The XMPP address has an invalid domain part"
	}

	return isValid, err
}

func verifyLocalPart(atIndex int) (bool, string) {
	isValid := true
	var err string

	if atIndex == 0 {
		isValid = false
		err = "The XMMP address has an invalid local part"
	}

	return isValid, err
}

package gui

import (
	"strings"
)

// Method to validate a jabber id is correct according to the RFC-6122
// on Address Format
// TODO: verify the resource part
func verifyXMPPAddress(address string) (bool, string) {
	var err string

	tldIndex := strings.LastIndex(address, ".")
	atIndex := strings.Index(address, "@")
	isValidDomain, errDomain := verifyDomainPart(tldIndex, atIndex, address)
	isValidPart, errPart := verifyLocalPart(atIndex)
	isValid := isValidDomain && isValidPart

	if !isValid {
		result := "Validation failed:\n"
		sep := ""
		if errDomain != "" {
			result += errDomain
			sep = ", "
		}
		if errPart != "" {
			result += sep + errPart
		}

		err = result + ". An XMPP address should look like this: local@domain.com."
	}

	return isValid, err
}

func verifyDomainPart(tldIndex, atIndex int, address string) (bool, string) {
	isValid := true
	var err string
	if tldIndex < 0 || tldIndex == len(address)-1 || tldIndex-atIndex < 2 {
		isValid = false
		err = "XMPP address has an invalid domain part"
	}

	return isValid, err
}

func verifyLocalPart(atIndex int) (bool, string) {
	isValid := true
	var err string
	if atIndex == 0 {
		isValid = false
		err = "XMMP address has an invalid local part"
	}

	return isValid, err
}

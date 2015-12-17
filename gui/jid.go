package gui

import (
	"strings"
)

func verify(address string) (bool, string) {
	var err string

	tldIndex := strings.Index(address, ".")
	atIndex := strings.Index(address, "@")
	isValidDomain, errDomain := verifyDomainPart(tldIndex, atIndex, address)
	isValidPart, errPart := verifyLocalPart(atIndex)
	isValid := isValidDomain && isValidPart

	if !isValid {
		errs := []string{errDomain, errPart, "XMPP address should look like local@domain.com"}
		err = strings.Join(errs, ", ")
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
	if atIndex < 1 {
		isValid = false
		err = "XMMP address has an invalid local part"
	}

	return isValid, err
}

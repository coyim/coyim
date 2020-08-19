package jid

import (
	"net"
	"regexp"

	"github.com/xdg/stringprep"
)

func stringToStringprepSet(s string) stringprep.Set {
	var result stringprep.Set

	for _, r := range s {
		result = append(result, stringprep.RuneRange{r, r})
	}

	return result
}

var localPartCustomExcludedTable = stringToStringprepSet("\"&'/:<>@")

var localPartProfile = stringprep.Profile{
	Mappings: []stringprep.Mapping{
		stringprep.TableB1,
		stringprep.TableB2,
	},
	Normalize: true,
	Prohibits: []stringprep.Set{
		stringprep.TableC1_1,
		stringprep.TableC1_2,
		stringprep.TableC2_1,
		stringprep.TableC2_2,
		stringprep.TableC3,
		stringprep.TableC4,
		stringprep.TableC5,
		stringprep.TableC6,
		stringprep.TableC7,
		stringprep.TableC8,
		stringprep.TableC9,
		localPartCustomExcludedTable,
	},
	CheckBiDi: true,
}

var resourcePartProfile = stringprep.Profile{
	Mappings: []stringprep.Mapping{
		stringprep.TableB1,
	},
	Normalize: true,
	Prohibits: []stringprep.Set{
		stringprep.TableC1_2,
		stringprep.TableC2_1,
		stringprep.TableC2_2,
		stringprep.TableC3,
		stringprep.TableC4,
		stringprep.TableC5,
		stringprep.TableC6,
		stringprep.TableC7,
		stringprep.TableC8,
		stringprep.TableC9,
	},
	CheckBiDi: true,
}

// ValidLocal checks whether the given string is a valid local part of a JID A localpart is defined to be any string of
// length 1 to 1023, matching the UsernameCaseMapped profile from RFC7613, and excluding a few more characters
func ValidLocal(s string) bool {
	ps, err := localPartProfile.Prepare(s)
	if err != nil {
		return false
	}

	l := len(ps)
	if l == 0 || l > 1023 {
		return false
	}

	return true
}

// Patterns taken from the govalidator project
var dnsName = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
var dnsReg = regexp.MustCompile(dnsName)

// ValidDomain returns true if the given string is a valid domain part for a JID
func ValidDomain(s string) bool {
	l := len(s)
	if l == 0 || l > 1023 {
		return false
	}

	res := net.ParseIP(s)
	if res != nil {
		return true
	}

	return dnsReg.MatchString(s)
}

// ValidResource returns true if the given string is a valid resource part for a JID.
// Note that a resource part is allowed to contain / and @ characters
func ValidResource(s string) bool {
	ps, err := resourcePartProfile.Prepare(s)
	if err != nil {
		return false
	}

	l := len(ps)
	if l == 0 || l > 1023 {
		return false
	}

	return true
}

// ValidJID returns true if the given string is any of the possible JID types
func ValidJID(s string) bool {
	return ValidFullJID(s) || ValidBareJID(s) || ValidDomain(s) || ValidDomainWithResource(s)
}

// ValidBareJID returns true if the given string is a valid bare JID. This function will true for full JIDs as well as
// bare JIDs
func ValidBareJID(s string) bool {
	res, ok := TryParseBare(s)
	return ok && res.Valid()
}

// ValidFullJID returns true if the given string is a valid full JID
func ValidFullJID(s string) bool {
	res, ok := TryParseFull(s)
	return ok && res.Valid()
}

// ValidDomainWithResource returns true if the given string a valid domain with resource part. This wil return true for
// a full JID, as well as a domain with JID
func ValidDomainWithResource(s string) bool {
	res := Parse(s)

	switch res.(type) {
	case full:
		return res.Valid()
	case domainWithResource:
		return res.Valid()
	}

	return false
}

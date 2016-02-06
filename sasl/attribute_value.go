package sasl

import (
	"errors"
	"regexp"
	"strings"
)

var (
	keyValue       = regexp.MustCompile(`^([^=]+)=(.+)$`)
	quotedKeyValue = regexp.MustCompile(`^([^=]+)="(.+)"$`)

	//ErrMissingParameter indicates a missing parameter from the server challenge
	ErrMissingParameter = errors.New("missing parameter in server challenge")
)

// AttributeValuePairs represents atribute-value pairs
type AttributeValuePairs map[string]string

// ParseAttributeValuePairs parses a string of comma-separated attribute=value pairs
func ParseAttributeValuePairs(src []byte) AttributeValuePairs {
	ret := make(AttributeValuePairs)
	params := strings.Split(string(src), ",")

	for _, p := range params {
		m := quotedKeyValue.FindStringSubmatch(p)
		if len(m) != 3 {
			m = keyValue.FindStringSubmatch(p)
		}

		if len(m) != 3 {
			continue
		}

		key := m[1]
		value := m[2]
		ret[key] = value
	}

	return ret
}

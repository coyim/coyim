package data

import (
	"encoding/xml"

	. "gopkg.in/check.v1"
)

type SASLSuite struct{}

var _ = Suite(&SASLSuite{})

func (s *SASLSuite) Test_SaslFailure_Condition_returnsTheCondition(c *C) {
	sf := SaslFailure{
		DefinedCondition: Any{
			XMLName: xml.Name{
				Local: "malformed-request",
			},
		},
	}

	c.Assert(sf.Condition(), Equals, SASLMalformedRequest)
}

func (s *SASLSuite) Test_SaslFailure_String_usesTheTextIfAvailable(c *C) {
	sf := SaslFailure{
		Text: "blubba",
		DefinedCondition: Any{
			XMLName: xml.Name{
				Local: "temporary-auth-failure",
			},
		},
	}

	c.Assert(sf.String(), Equals, "temporary-auth-failure: \"blubba\"")
}

func (s *SASLSuite) Test_SaslFailure_String_returnsConditionIfNoTextAvailable(c *C) {
	sf := SaslFailure{
		DefinedCondition: Any{
			XMLName: xml.Name{
				Local: "account-disabled",
			},
		},
	}

	c.Assert(sf.String(), Equals, "account-disabled")
}

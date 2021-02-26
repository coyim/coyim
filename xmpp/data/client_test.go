package data

import (
	"encoding/xml"

	. "gopkg.in/check.v1"
)

type ClientSuite struct{}

var _ = Suite(&ClientSuite{})

func (s *ClientSuite) Test_Extension_GoString(c *C) {
	var te *Extension

	c.Assert(te.GoString(), Equals, "<nil>")
	te = &Extension{
		XMLName: xml.Name{
			Space: "foo.bar.com",
			Local: "hello",
		},
		Body: "blublu",
	}
	c.Assert(te.GoString(), Equals, "<Extension {foo.bar.com hello} body=\"blublu\">")
}

func (s *ClientSuite) Test_Extensions_GoString(c *C) {
	var te *Extensions

	c.Assert(te.GoString(), Equals, "<nil>")

	te = &Extensions{
		&Extension{
			XMLName: xml.Name{
				Space: "foo.bar.com",
				Local: "hello",
			},
			Body: "blublu",
		},
		&Extension{
			XMLName: xml.Name{
				Space: "stuff.org",
				Local: "something",
			},
		},
	}
	c.Assert(te.GoString(), Equals, "{<Extension {foo.bar.com hello} body=\"blublu\">, <Extension {stuff.org something} body=\"\">}")
}

func (s *ClientSuite) Test_StanzaError_AnyMUCError(c *C) {
	mec := &MUCNotAllowed{}
	se := &StanzaError{
		MUCNotAllowed: mec,
	}

	c.Assert(se.AnyMUCError(), Equals, mec)
	c.Assert((&StanzaError{}).AnyMUCError(), IsNil)
}

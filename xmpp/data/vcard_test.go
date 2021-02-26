package data

import (
	. "gopkg.in/check.v1"
)

type VCardSuite struct{}

var _ = Suite(&VCardSuite{})

func (s *VCardSuite) Test_ParseVCard_failsOnBadStanza(c *C) {
	st := Stanza{
		Value: "foo",
	}
	res, e := ParseVCard(st)
	c.Assert(e, ErrorMatches, "xmpp: vcard request resulted in tag of type.*")
	c.Assert(res, Equals, VCard{})
}

func (s *VCardSuite) Test_ParseVCard_failsOnBadXML(c *C) {
	st := Stanza{
		Value: &ClientIQ{
			Query: []byte("<badXML"),
		},
	}
	res, e := ParseVCard(st)
	c.Assert(e, ErrorMatches, "XML syntax error on.*")
	c.Assert(res, Equals, VCard{})
}

func (s *VCardSuite) Test_ParseVCard_succeeds(c *C) {
	x := `
<vCard xmlns='vcard-temp'>
  <FN>Hello</FN>
  <NICKNAME>Someone</NICKNAME>
</vCard>
`

	st := Stanza{
		Value: &ClientIQ{
			Query: []byte(x),
		},
	}
	res, e := ParseVCard(st)
	c.Assert(e, IsNil)
	c.Assert(res.FullName, Equals, "Hello")
	c.Assert(res.Nickname, Equals, "Someone")
}

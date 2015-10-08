package xmpp

import . "gopkg.in/check.v1"

type CapabilitiesXmppSuite struct{}

var _ = Suite(&CapabilitiesXmppSuite{})

func (s *CapabilitiesXmppSuite) Test_DiscoveryIdentity_xep0115Less_comparesCategory(c *C) {
	left := &DiscoveryIdentity{}
	right := &DiscoveryIdentity{}

	left.Category = "A"
	right.Category = "B"
	c.Assert(left.xep0115Less(right), Equals, true)
	c.Assert(right.xep0115Less(left), Equals, false)

	left.Category = "B"
	right.Category = "A"
	c.Assert(left.xep0115Less(right), Equals, false)
	c.Assert(right.xep0115Less(left), Equals, true)
}

func (s *CapabilitiesXmppSuite) Test_DiscoveryIdentity_xep0115Less_comparesType(c *C) {
	left := &DiscoveryIdentity{Category: "A"}
	right := &DiscoveryIdentity{Category: "A"}

	left.Type = "A"
	right.Type = "B"
	c.Assert(left.xep0115Less(right), Equals, true)
	c.Assert(right.xep0115Less(left), Equals, false)

	left.Type = "B"
	right.Type = "A"
	c.Assert(left.xep0115Less(right), Equals, false)
	c.Assert(right.xep0115Less(left), Equals, true)
}

func (s *CapabilitiesXmppSuite) Test_DiscoveryIdentity_xep0115Less_comparesLang(c *C) {
	left := &DiscoveryIdentity{Category: "A", Type: "B"}
	right := &DiscoveryIdentity{Category: "A", Type: "B"}

	left.Lang = "A"
	right.Lang = "B"
	c.Assert(left.xep0115Less(right), Equals, true)
	c.Assert(right.xep0115Less(left), Equals, false)

	left.Lang = "B"
	right.Lang = "A"
	c.Assert(left.xep0115Less(right), Equals, false)
	c.Assert(right.xep0115Less(left), Equals, true)
}

func (s *CapabilitiesXmppSuite) Test_formField_xep0115Less_comparesVar(c *C) {
	left := &formField{}
	right := &formField{}

	left.Var = "FORM_TYPE"
	right.Var = "FORM_TYPE2"
	c.Assert(left.xep0115Less(right), Equals, true)
	c.Assert(right.xep0115Less(left), Equals, false)

	left.Var = "FORM_TYPE2"
	right.Var = "FORM_TYPE"
	c.Assert(left.xep0115Less(right), Equals, false)
	c.Assert(right.xep0115Less(left), Equals, true)

	left.Var = "FORM_TYPE2"
	right.Var = "FORM_TYPE3"
	c.Assert(left.xep0115Less(right), Equals, true)
	c.Assert(right.xep0115Less(left), Equals, false)

	left.Var = "FORM_TYPE3"
	right.Var = "FORM_TYPE2"
	c.Assert(left.xep0115Less(right), Equals, false)
	c.Assert(right.xep0115Less(left), Equals, true)
}

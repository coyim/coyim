package xmpp

import (
	"github.com/twstrike/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type CapabilitiesXMPPSuite struct{}

var _ = Suite(&CapabilitiesXMPPSuite{})

func (s *CapabilitiesXMPPSuite) Test_DiscoveryIdentity_xep0115Less_comparesCategory(c *C) {
	left := &data.DiscoveryIdentity{}
	right := &data.DiscoveryIdentity{}

	left.Category = "A"
	right.Category = "B"
	c.Assert(xep0115Less(left, right), Equals, true)
	c.Assert(xep0115Less(right, left), Equals, false)

	left.Category = "B"
	right.Category = "A"
	c.Assert(xep0115Less(left, right), Equals, false)
	c.Assert(xep0115Less(right, left), Equals, true)
}

func (s *CapabilitiesXMPPSuite) Test_DiscoveryIdentity_xep0115Less_comparesType(c *C) {
	left := &data.DiscoveryIdentity{Category: "A"}
	right := &data.DiscoveryIdentity{Category: "A"}

	left.Type = "A"
	right.Type = "B"
	c.Assert(xep0115Less(left, right), Equals, true)
	c.Assert(xep0115Less(right, left), Equals, false)

	left.Type = "B"
	right.Type = "A"
	c.Assert(xep0115Less(left, right), Equals, false)
	c.Assert(xep0115Less(right, left), Equals, true)
}

func (s *CapabilitiesXMPPSuite) Test_DiscoveryIdentity_xep0115Less_comparesLang(c *C) {
	left := &data.DiscoveryIdentity{Category: "A", Type: "B"}
	right := &data.DiscoveryIdentity{Category: "A", Type: "B"}

	left.Lang = "A"
	right.Lang = "B"
	c.Assert(xep0115Less(left, right), Equals, true)
	c.Assert(xep0115Less(right, left), Equals, false)

	left.Lang = "B"
	right.Lang = "A"
	c.Assert(xep0115Less(left, right), Equals, false)
	c.Assert(xep0115Less(right, left), Equals, true)
}

func (s *CapabilitiesXMPPSuite) Test_formField_xep0115Less_comparesVar(c *C) {
	left := &data.FormFieldX{}
	right := &data.FormFieldX{}

	left.Var = "FORM_TYPE"
	right.Var = "FORM_TYPE2"
	c.Assert(xep0115Less(left, right), Equals, true)
	c.Assert(xep0115Less(right, left), Equals, false)

	left.Var = "FORM_TYPE2"
	right.Var = "FORM_TYPE"
	c.Assert(xep0115Less(left, right), Equals, false)
	c.Assert(xep0115Less(right, left), Equals, true)

	left.Var = "FORM_TYPE2"
	right.Var = "FORM_TYPE3"
	c.Assert(xep0115Less(left, right), Equals, true)
	c.Assert(xep0115Less(right, left), Equals, false)

	left.Var = "FORM_TYPE3"
	right.Var = "FORM_TYPE2"
	c.Assert(xep0115Less(left, right), Equals, false)
	c.Assert(xep0115Less(right, left), Equals, true)
}

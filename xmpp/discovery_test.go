package xmpp

import (
	"github.com/twstrike/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type DiscoveryXMPPSuite struct{}

var _ = Suite(&DiscoveryXMPPSuite{})

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreDuplicateIdentities(c *C) {
	reply := &data.DiscoveryReply{
		Identities: []data.DiscoveryIdentity{
			data.DiscoveryIdentity{
				Lang:     "en",
				Category: "stuff",
				Type:     "thing",
				Name:     "something",
			},
			data.DiscoveryIdentity{
				Lang:     "en",
				Category: "stuff",
				Type:     "thing",
				Name:     "something",
			},
		},
	}

	_, err := VerificationString(reply)
	c.Assert(err.Error(), Equals, "duplicate discovery identity")
}

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreDuplicateFeatures(c *C) {
	reply := &data.DiscoveryReply{
		Features: []data.DiscoveryFeature{
			data.DiscoveryFeature{
				Var: "foo",
			},
			data.DiscoveryFeature{
				Var: "foo",
			},
		},
	}

	_, err := VerificationString(reply)
	c.Assert(err.Error(), Equals, "duplicate discovery feature")
}

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreDuplicateFormTypes(c *C) {
	reply := &data.DiscoveryReply{
		Forms: []data.Form{
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:    "FORM_TYPE",
						Type:   "Foo",
						Values: []string{"foo"},
					},
				},
			},
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:    "FORM_TYPE",
						Type:   "Foo",
						Values: []string{"foo"},
					},
				},
			},
		},
	}

	_, err := VerificationString(reply)
	c.Assert(err.Error(), Equals, "multiple forms of the same type")
}

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreNoValues(c *C) {
	reply := &data.DiscoveryReply{
		Forms: []data.Form{
			data.Form{
				Type: "foo",
			},
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "FORM_TYPE2",
						Type: "Foo",
					},
				},
			},
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "FORM_TYPE",
						Type: "Foo",
					},
				},
			},
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "FORM_TYPE",
						Type: "Foo",
					},
				},
			},
		},
	}

	_, err := VerificationString(reply)
	c.Assert(err.Error(), Equals, "form does not have a single FORM_TYPE value")
}

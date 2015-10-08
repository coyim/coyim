package xmpp

import . "gopkg.in/check.v1"

type DiscoveryXmppSuite struct{}

var _ = Suite(&DiscoveryXmppSuite{})

func (s *DiscoveryXmppSuite) Test_VerificationString_failsIfThereAreDuplicateIdentities(c *C) {
	reply := DiscoveryReply{
		Identities: []DiscoveryIdentity{
			DiscoveryIdentity{
				Lang:     "en",
				Category: "stuff",
				Type:     "thing",
				Name:     "something",
			},
			DiscoveryIdentity{
				Lang:     "en",
				Category: "stuff",
				Type:     "thing",
				Name:     "something",
			},
		},
	}

	_, err := reply.VerificationString()
	c.Assert(err.Error(), Equals, "duplicate discovery identity")
}

func (s *DiscoveryXmppSuite) Test_VerificationString_failsIfThereAreDuplicateFeatures(c *C) {
	reply := DiscoveryReply{
		Features: []DiscoveryFeature{
			DiscoveryFeature{
				Var: "foo",
			},
			DiscoveryFeature{
				Var: "foo",
			},
		},
	}

	_, err := reply.VerificationString()
	c.Assert(err.Error(), Equals, "duplicate discovery feature")
}

func (s *DiscoveryXmppSuite) Test_VerificationString_failsIfThereAreDuplicateFormTypes(c *C) {
	reply := DiscoveryReply{
		Forms: []Form{
			Form{
				Type: "foo",
				Fields: []formField{
					formField{
						Var:    "FORM_TYPE",
						Type:   "Foo",
						Values: []string{"foo"},
					},
				},
			},
			Form{
				Type: "foo",
				Fields: []formField{
					formField{
						Var:    "FORM_TYPE",
						Type:   "Foo",
						Values: []string{"foo"},
					},
				},
			},
		},
	}

	_, err := reply.VerificationString()
	c.Assert(err.Error(), Equals, "multiple forms of the same type")
}

func (s *DiscoveryXmppSuite) Test_VerificationString_failsIfThereAreNoValues(c *C) {
	reply := DiscoveryReply{
		Forms: []Form{
			Form{
				Type: "foo",
			},
			Form{
				Type: "foo",
				Fields: []formField{
					formField{
						Var:  "FORM_TYPE2",
						Type: "Foo",
					},
				},
			},
			Form{
				Type: "foo",
				Fields: []formField{
					formField{
						Var:  "FORM_TYPE",
						Type: "Foo",
					},
				},
			},
			Form{
				Type: "foo",
				Fields: []formField{
					formField{
						Var:  "FORM_TYPE",
						Type: "Foo",
					},
				},
			},
		},
	}

	_, err := reply.VerificationString()
	c.Assert(err.Error(), Equals, "form does not have a single FORM_TYPE value")
}

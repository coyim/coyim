package config

import (
	"sort"

	. "gopkg.in/check.v1"
)

type AccountPeerSuite struct{}

var _ = Suite(&AccountPeerSuite{})

func (s *AccountPeerSuite) Test_Peer_ChangesNicknameAndGroups(c *C) {
	pid := "peer@coy.im"
	a := &Account{
		Account: "account@coy.im",
		Peers: []*Peer{
			&Peer{
				UserID: pid,
			},
		},
	}

	a.SavePeerDetails(pid, "bff", []string{"Foes"})

	p, _ := a.GetPeer(pid)
	c.Check(p.Nickname, Equals, "bff")
	c.Check(p.Groups, DeepEquals, []string{"Foes"})
}

func (s *AccountPeerSuite) Test_Fingerprint_MarshalJSON(c *C) {
	f := &Fingerprint{
		Fingerprint: []byte{0x01, 0x03, 0x02},
		Trusted:     true,
		Tag:         "something",
	}

	res, e := f.MarshalJSON()
	c.Assert(e, IsNil)
	c.Assert(string(res), Equals, "{\"FingerprintHex\":\"010302\",\"Trusted\":true,\"Tag\":\"something\"}")
}

func (s *AccountPeerSuite) Test_Fingerprint_UnmarshalJSON_succeeds(c *C) {
	f := &Fingerprint{}
	e := f.UnmarshalJSON([]byte("{\"FingerprintHex\":\"010302\",\"Trusted\":true,\"Tag\":\"something\"}"))
	c.Assert(e, IsNil)
	c.Assert(f.Trusted, Equals, true)
	c.Assert(f.Tag, Equals, "something")
	c.Assert(f.Fingerprint, DeepEquals, []byte{0x01, 0x03, 0x02})
}

func (s *AccountPeerSuite) Test_Fingerprint_UnmarshalJSON_failsOnBadJSON(c *C) {
	f := &Fingerprint{}
	e := f.UnmarshalJSON([]byte("{\"FingerprintHex\":\"010302\",\"Tr"))
	c.Assert(e, ErrorMatches, ".*unexpected end of JSON input.*")
}

func (s *AccountPeerSuite) Test_Fingerprint_UnmarshalJSON_failsOnBadHex(c *C) {
	f := &Fingerprint{}
	e := f.UnmarshalJSON([]byte("{\"FingerprintHex\":\"01q302\",\"Trusted\":true,\"Tag\":\"something\"}"))
	c.Assert(e, ErrorMatches, ".*encoding/hex.*")
}

func (s *AccountPeerSuite) Test_Fingerprint_ByNaturalOrder(c *C) {
	f1 := &Fingerprint{Fingerprint: []byte{0x01, 0x02}}
	f2 := &Fingerprint{Fingerprint: []byte{0x01, 0x01}}
	f3 := &Fingerprint{Fingerprint: []byte{0x02, 0x01}}
	one := []*Fingerprint{f1, f2, f3}
	sort.Sort(ByNaturalOrder(one))
	c.Assert(one[0], Equals, f2)
	c.Assert(one[1], Equals, f1)
	c.Assert(one[2], Equals, f3)
}

func (s *AccountPeerSuite) Test_Account_UpdateEncryptionRequired_Require(c *C) {
	a := &Account{}
	a.UpdateEncryptionRequired("some@one.com", true)

	c.Assert(a.DontEncryptWith, HasLen, 0)
	c.Assert(a.AlwaysEncryptWith, HasLen, 1)
	c.Assert(a.AlwaysEncryptWith[0], Equals, "some@one.com")
	c.Assert(a.Peers, HasLen, 1)
	c.Assert(a.Peers[0].EncryptionSettings, Equals, EncryptionSettings(AlwaysEncrypt))
}

func (s *AccountPeerSuite) Test_Account_UpdateEncryptionRequired_DontRequire(c *C) {
	a := &Account{}
	a.UpdateEncryptionRequired("some@one.com", false)

	c.Assert(a.AlwaysEncryptWith, HasLen, 0)
	c.Assert(a.DontEncryptWith, HasLen, 1)
	c.Assert(a.DontEncryptWith[0], Equals, "some@one.com")
	c.Assert(a.Peers, HasLen, 1)
	c.Assert(a.Peers[0].EncryptionSettings, Equals, EncryptionSettings(NeverEncrypt))
}

func (s *AccountPeerSuite) Test_Account_SavePeerDetails(c *C) {
	a := &Account{}
	a.SavePeerDetails("some@one.com", "hubba", []string{"left", "right"})

	c.Assert(a.Peers, HasLen, 1)
	c.Assert(a.Peers[0].UserID, Equals, "some@one.com")
	c.Assert(a.Peers[0].Nickname, Equals, "hubba")
	c.Assert(a.Peers[0].Groups, DeepEquals, []string{"left", "right"})
}

func (s *AccountPeerSuite) Test_Account_RemovePeer(c *C) {
	a := &Account{
		Peers: []*Peer{
			&Peer{UserID: "one"},
			&Peer{UserID: "two"},
			&Peer{UserID: "three"},
		},
	}

	a.RemovePeer("four")

	c.Assert(a.Peers, HasLen, 3)

	a.RemovePeer("one")

	c.Assert(a.Peers, HasLen, 2)
	c.Assert(a.Peers[0].UserID, Equals, "two")
	c.Assert(a.Peers[1].UserID, Equals, "three")
}

func (s *AccountPeerSuite) Test_Account_RemoveFingerprint_forUnknownPeer(c *C) {
	a := &Account{}
	res := a.RemoveFingerprint("four", []byte{0x01, 0x02})
	c.Assert(res, Equals, false)
}

func (s *AccountPeerSuite) Test_Account_RemoveFingerprint_thatPeerDoesntHave(c *C) {
	a := &Account{
		Peers: []*Peer{
			&Peer{UserID: "one"},
			&Peer{UserID: "four", Fingerprints: []*Fingerprint{}},
			&Peer{UserID: "three"},
		},
	}
	res := a.RemoveFingerprint("four", []byte{0x01, 0x02})
	c.Assert(res, Equals, false)
}

func (s *AccountPeerSuite) Test_Account_RemoveFingerprint_thatPeerDoesHave(c *C) {
	a := &Account{
		Peers: []*Peer{
			&Peer{UserID: "one"},
			&Peer{UserID: "four", Fingerprints: []*Fingerprint{
				&Fingerprint{Fingerprint: []byte{0xFF, 0xFE}},
				&Fingerprint{Fingerprint: []byte{0x01, 0x02}},
				&Fingerprint{Fingerprint: []byte{0x99, 0xAA}},
			}},
			&Peer{UserID: "three"},
		},
	}
	res := a.RemoveFingerprint("four", []byte{0x01, 0x02})
	c.Assert(res, Equals, true)
	c.Assert(a.Peers[1].Fingerprints, HasLen, 2)
}

func (s *AccountPeerSuite) Test_Account_AuthorizeFingerprint_returnsErrorOnAlreadyAuthorizedFingerprint(c *C) {
	a := &Account{
		Peers: []*Peer{
			&Peer{UserID: "one"},
			&Peer{UserID: "four", Fingerprints: []*Fingerprint{
				&Fingerprint{Fingerprint: []byte{0xFF, 0xFE}},
				&Fingerprint{Fingerprint: []byte{0x01, 0x02}, Trusted: true},
				&Fingerprint{Fingerprint: []byte{0x99, 0xAA}},
			}},
			&Peer{UserID: "three"},
		},
	}
	res := a.AuthorizeFingerprint("one", []byte{0x01, 0x02}, "somewhere")
	c.Assert(res, Equals, errFingerprintAlreadyAuthorized)
}

func (s *AccountPeerSuite) Test_Account_AuthorizeFingerprint_updatesTrustOnExistingFingerprint(c *C) {
	fpr := &Fingerprint{Fingerprint: []byte{0xFF, 0xFE}}
	a := &Account{
		Peers: []*Peer{
			&Peer{UserID: "one"},
			&Peer{UserID: "four", Fingerprints: []*Fingerprint{
				fpr,
				&Fingerprint{Fingerprint: []byte{0x01, 0x02}, Trusted: true},
				&Fingerprint{Fingerprint: []byte{0x99, 0xAA}},
			}},
			&Peer{UserID: "three"},
		},
	}
	res := a.AuthorizeFingerprint("four", []byte{0xFF, 0xFE}, "somewhere")
	c.Assert(res, IsNil)
	c.Assert(fpr.Trusted, Equals, true)
	c.Assert(fpr.Tag, Equals, "somewhere")
}

func (s *AccountPeerSuite) Test_Account_AuthorizeFingerprint_addsNewFingerprint(c *C) {
	a := &Account{
		Peers: []*Peer{
			&Peer{UserID: "one"},
			&Peer{UserID: "four", Fingerprints: []*Fingerprint{
				&Fingerprint{Fingerprint: []byte{0x01, 0x02}, Trusted: true},
				&Fingerprint{Fingerprint: []byte{0x99, 0xAA}},
			}},
			&Peer{UserID: "three"},
		},
	}
	res := a.AuthorizeFingerprint("four", []byte{0xFF, 0xFE}, "somewhere2")
	c.Assert(res, IsNil)
	c.Assert(a.Peers[1].Fingerprints, HasLen, 3)
	c.Assert(a.Peers[1].Fingerprints[2].Fingerprint, DeepEquals, []byte{0xFF, 0xFE})
	c.Assert(a.Peers[1].Fingerprints[2].Trusted, Equals, true)
	c.Assert(a.Peers[1].Fingerprints[2].Tag, Equals, "somewhere2")
}

func (s *AccountPeerSuite) Test_Account_HasFingerprint(c *C) {
	a := &Account{
		Peers: []*Peer{
			&Peer{UserID: "one"},
			&Peer{UserID: "four", Fingerprints: []*Fingerprint{
				&Fingerprint{Fingerprint: []byte{0x01, 0x02}, Trusted: true},
				&Fingerprint{Fingerprint: []byte{0x99, 0xAA}},
			}},
			&Peer{UserID: "three"},
		},
	}
	c.Assert(a.HasFingerprint("unknown"), Equals, false)
	c.Assert(a.HasFingerprint("one"), Equals, false)
	c.Assert(a.HasFingerprint("four"), Equals, true)
}

func (s *AccountPeerSuite) Test_Peer_GetFingerprint(c *C) {
	fpr := &Fingerprint{Fingerprint: []byte{0x01, 0x02}, Trusted: true}
	peer := &Peer{UserID: "four", Fingerprints: []*Fingerprint{
		&Fingerprint{Fingerprint: []byte{0x99, 0xAA}},
		fpr,
	}}

	res, ok := peer.GetFingerprint([]byte{0x01, 0x02})
	c.Assert(res, Equals, fpr)
	c.Assert(ok, Equals, true)
	res, ok = peer.GetFingerprint([]byte{0x01, 0x03})
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *AccountPeerSuite) Test_Account_removeEmptyFingerprints(c *C) {
	a := &Account{
		Peers: []*Peer{
			&Peer{UserID: "one", Fingerprints: []*Fingerprint{
				&Fingerprint{Fingerprint: []byte{0x99, 0xAA}},
				&Fingerprint{Fingerprint: []byte{}},
			}},
			&Peer{UserID: "four", Fingerprints: []*Fingerprint{
				&Fingerprint{Fingerprint: []byte{0x01, 0x02}, Trusted: true},
				&Fingerprint{Fingerprint: []byte{0x99, 0xAA}},
			}},
			&Peer{UserID: "three"},
		},
	}

	c.Assert(a.removeEmptyFingerprints(), Equals, true)
	c.Assert(a.Peers[0].Fingerprints, HasLen, 1)
	c.Assert(a.Peers[1].Fingerprints, HasLen, 2)
	c.Assert(a.Peers[2].Fingerprints, HasLen, 0)
}

func (s *AccountPeerSuite) Test_Account_updateFingerprintsToLatestVersion_returnsFalseIfNoLegacyFingerprints(c *C) {
	a := &Account{}

	c.Assert(a.updateFingerprintsToLatestVersion(), Equals, false)
}

func (s *AccountPeerSuite) Test_Account_updateFingerprintsToLatestVersion_addsLegacyFingerprints(c *C) {
	a := &Account{
		LegacyKnownFingerprints: []KnownFingerprint{
			KnownFingerprint{
				UserID:      "one@some.org",
				Fingerprint: []byte{0x01, 0x02, 0x03},
				Untrusted:   true,
			},
			KnownFingerprint{
				UserID:      "ignored@fingerprint.com",
				Fingerprint: []byte{},
				Untrusted:   true,
			},
			KnownFingerprint{
				UserID:      "one@some.org",
				Fingerprint: []byte{0x02, 0x02, 0x05},
				Untrusted:   false,
			},
		},
	}

	c.Assert(a.updateFingerprintsToLatestVersion(), Equals, true)
	c.Assert(a.LegacyKnownFingerprints, HasLen, 0)
	c.Assert(a.Peers, HasLen, 1)
	c.Assert(a.Peers[0].Fingerprints, HasLen, 2)
	c.Assert(a.Peers[0].Fingerprints[0].Fingerprint, DeepEquals, []byte{0x01, 0x02, 0x03})
	c.Assert(a.Peers[0].Fingerprints[0].Trusted, Equals, false)
	c.Assert(a.Peers[0].Fingerprints[0].Tag, Equals, "")
	c.Assert(a.Peers[0].Fingerprints[1].Fingerprint, DeepEquals, []byte{0x02, 0x02, 0x05})
	c.Assert(a.Peers[0].Fingerprints[1].Trusted, Equals, true)
	c.Assert(a.Peers[0].Fingerprints[1].Tag, Equals, "")
}

func (s *AccountPeerSuite) Test_Account_updateToLatestVersion_doesNothingForEmptyAccount(c *C) {
	a := &Account{}

	c.Assert(a.updateToLatestVersion(), Equals, false)
}

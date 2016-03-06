package config

import (
	"encoding/hex"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/otr3"

	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
)

type AccountXmppSuite struct{}

var _ = Suite(&AccountXmppSuite{})

func (s *AccountXmppSuite) Test_Account_Is_recognizesJids(c *C) {
	a := &Account{Account: "hello@bar.com"}
	c.Check(a.Is("foo"), Equals, false)
	c.Check(a.Is("hello@bar.com"), Equals, true)
	c.Check(a.Is("hello@bar.com/foo"), Equals, true)
}

func (s *AccountXmppSuite) Test_Account_ShouldEncryptTo(c *C) {
	a := &Account{Account: "hello@bar.com", AlwaysEncrypt: false, AlwaysEncryptWith: []string{"one@foo.com", "two@foo.com"}}
	a2 := &Account{Account: "hello@bar.com", AlwaysEncrypt: true, AlwaysEncryptWith: []string{"one@foo.com", "two@foo.com"}}
	c.Check(a.ShouldEncryptTo("foo"), Equals, false)
	c.Check(a.ShouldEncryptTo("hello@bar.com"), Equals, false)
	c.Check(a.ShouldEncryptTo("one@foo.com"), Equals, true)
	c.Check(a.ShouldEncryptTo("two@foo.com"), Equals, true)
	c.Check(a.ShouldEncryptTo("two@foo.com/blarg"), Equals, true)
	c.Check(a2.ShouldEncryptTo("foo"), Equals, true)
	c.Check(a2.ShouldEncryptTo("hello@bar.com"), Equals, true)
}

func (s *AccountXmppSuite) Test_NewAccount_ReturnsNewAccountWithSafeDefaults(c *C) {
	a, err := NewAccount()

	c.Check(err, IsNil)
	c.Check(len(a.PrivateKeys), Equals, 1)
	c.Check(a.AlwaysEncrypt, Equals, true)
	c.Check(a.OTRAutoStartSession, Equals, true)
	c.Check(a.OTRAutoTearDown, Equals, true)
	c.Check(a.Proxies, DeepEquals, []string{"tor-auto://"})
}

func (s *AccountXmppSuite) Test_SetOTRPoliciesFor_SetupOTRPolicies(c *C) {
	a, _ := NewAccount()
	conv := &otr3.Conversation{}

	expectedConv := &otr3.Conversation{}
	expectedPolicies := expectedConv.Policies
	expectedPolicies.AllowV2()
	expectedPolicies.AllowV3()
	expectedPolicies.SendWhitespaceTag()
	expectedPolicies.WhitespaceStartAKE()
	expectedPolicies.RequireEncryption()
	expectedPolicies.ErrorStartAKE()

	a.SetOTRPoliciesFor("someon@jabber.com", conv)
	c.Check(conv.Policies, Equals, expectedPolicies)
}

func (s *AccountXmppSuite) Test_SetOTRPoliciesFor_SetupOTRPoliciesWithOptionalEncription(c *C) {
	a, _ := NewAccount()
	a.AlwaysEncrypt = false
	conv := &otr3.Conversation{}

	expectedConv := &otr3.Conversation{}
	expectedPolicies := expectedConv.Policies
	expectedPolicies.AllowV2()
	expectedPolicies.AllowV3()
	expectedPolicies.SendWhitespaceTag()
	expectedPolicies.WhitespaceStartAKE()

	a.SetOTRPoliciesFor("someon@jabber.com", conv)
	c.Check(conv.Policies, Equals, expectedPolicies)
}

func (s *AccountXmppSuite) Test_EnsurePrivateKey_DoesNotUpdateIfKeyExists(c *C) {
	a, _ := NewAccount()
	changed, err := a.EnsurePrivateKey()

	c.Check(err, IsNil)
	c.Check(changed, Equals, false)
}

func (s *AccountXmppSuite) Test_EnsurePrivateKey_GeneratePrivateKeyIfMissing(c *C) {
	a := &Account{}
	changed, err := a.EnsurePrivateKey()

	c.Check(err, IsNil)
	c.Check(changed, Equals, true)
	c.Check(len(a.PrivateKeys), Equals, 1)
}

func (s *AccountXmppSuite) Test_ID_generatesID(c *C) {
	a := &Account{}
	c.Check(a.ID(), Not(HasLen), 0)
}

func (s *AccountXmppSuite) Test_ID_doesNotChangeID(c *C) {
	a := &Account{
		id: "existing",
	}
	c.Check(a.ID(), Equals, "existing")
}

func (s *AccountXmppSuite) Test_ServerCertificateHash_deserializeServerCertificateHash(c *C) {
	expectedCertificateHash := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
	}

	a := &Account{
		ServerCertificateSHA256: hex.EncodeToString(expectedCertificateHash),
	}

	serverHash, err := a.ServerCertificateHash()
	c.Check(err, IsNil)
	c.Check(serverHash, DeepEquals, expectedCertificateHash)
}

func (s *AccountXmppSuite) Test_ServerCertificateHash_ErrorWhenFailsToDeserializeHash(c *C) {
	a := &Account{
		ServerCertificateSHA256: "af3",
	}

	_, err := a.ServerCertificateHash()
	c.Check(err.Error(), Equals, "Failed to parse ServerCertificateSHA256 (should be hex string): encoding/hex: odd length hex string")
}

func (s *AccountXmppSuite) Test_ServerCertificateHash_ErrorWHenHashHasDifferentSize(c *C) {
	expectedCertificateHash := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0xAA,
	}

	a := &Account{
		ServerCertificateSHA256: hex.EncodeToString(expectedCertificateHash),
	}

	_, err := a.ServerCertificateHash()
	c.Check(err, Equals, errCertificateSizeMismatch)
}

package importer

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coyim/coyim/config"
	. "gopkg.in/check.v1"
)

type XMPPClientXMPPSuite struct{}

var _ = Suite(&XMPPClientXMPPSuite{})

func decode(in string) []byte {
	ret, _ := hex.DecodeString(in)
	return ret
}

func (s *XMPPClientXMPPSuite) Test_XmppClient_canImportXmppClientConfiguration(c *C) {
	importer := xmppClientImporter{}
	res, ok := importer.importFrom(testResourceFilename("xmpp_client_test_conf.json"))
	c.Assert(ok, Equals, true)

	c.Assert(len(res.Accounts), Equals, 1)

	c.Assert(res.Accounts[0].Account, Equals, "ox@coyim.net")
	c.Assert(res.Accounts[0].Server, Equals, "xmpp.coyim.net")
	c.Assert(res.Accounts[0].Proxies[0], Equals, "socks5://127.0.0.1:9051")
	c.Assert(res.Accounts[0].Password, Equals, "123547567846rghdfghdrghr6ythdt")
	c.Assert(res.Accounts[0].Port, Equals, 5223)
	c.Assert(res.Accounts[0].HideStatusUpdates, Equals, true)
	c.Assert(res.Accounts[0].OTRAutoTearDown, Equals, true)
	c.Assert(res.Accounts[0].OTRAutoAppendTag, Equals, true)
	c.Assert(res.Accounts[0].OTRAutoStartSession, Equals, true)
	c.Assert(res.Accounts[0].LegacyServerCertificateSHA256, Equals, "592f46183527ab40838882ab4cb4aef4e2cf916074ab01f9bc243931ca5c4ed1")
	c.Assert(res.Accounts[0].PrivateKeys[0], DeepEquals, []byte{0x00, 0x10, 0x80, 0x04, 0x20, 0x01})
	c.Assert(res.Accounts[0].AlwaysEncrypt, Equals, true)
	c.Assert(res.Accounts[0].AlwaysEncryptWith, DeepEquals, []string(nil))
	c.Assert(res.Accounts[0].InstanceTag, Equals, uint32(0))

	c.Assert(len(res.Accounts[0].LegacyKnownFingerprints), Equals, 0)
	// c.Assert(res.Accounts[0].KnownFingerprints[0].UserID, Equals, "arnold@jabber.ccc.de")
	// c.Assert(res.Accounts[0].KnownFingerprints[0].Fingerprint, DeepEquals, decode("c2a23b8e8852bff5335b39b674ceec13228be0af"))
	// c.Assert(res.Accounts[0].KnownFingerprints[0].Untrusted, Equals, false)
	// c.Assert(res.Accounts[0].KnownFingerprints[1].UserID, Equals, "some@one.com")
	// c.Assert(res.Accounts[0].KnownFingerprints[1].Fingerprint, DeepEquals, decode("410aad3ce865b83ed564b2e1ce52882b07b00976"))
	// c.Assert(res.Accounts[0].KnownFingerprints[1].Untrusted, Equals, false)
	// c.Assert(res.Accounts[0].KnownFingerprints[2].UserID, Equals, "hello@riseup.net")
	// c.Assert(res.Accounts[0].KnownFingerprints[2].Fingerprint, DeepEquals, decode("50ae9522641401e1a58de568fc4b265493d451b4"))
	// c.Assert(res.Accounts[0].KnownFingerprints[2].Untrusted, Equals, false)
	// c.Assert(res.Accounts[0].KnownFingerprints[3].UserID, Equals, "second.hello@riseup.net")
	// c.Assert(res.Accounts[0].KnownFingerprints[3].Fingerprint, DeepEquals, decode("4da542fdd60e077b38b05aa3485916d5d7c958aa"))
	// c.Assert(res.Accounts[0].KnownFingerprints[3].Untrusted, Equals, false)

	c.Assert(res.Bell, Equals, true)
	c.Assert(res.Display.MergeAccounts, Equals, false)
	c.Assert(res.Display.ShowOnlyOnline, Equals, false)
	c.Assert(res.RawLogFile, Equals, "bla")
	c.Assert(res.NotifyCommand, DeepEquals, []string{"hello"})
	c.Assert(res.IdleSecondsBeforeNotification, Equals, 42)
}

func (s *XMPPClientXMPPSuite) Test_xmppClientImporter_TryImport(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", origHome)
	}()
	os.Setenv("HOME", dir)

	input, _ := ioutil.ReadFile(testResourceFilename("xmpp_client_test_conf.json"))
	_ = ioutil.WriteFile(filepath.Join(dir, ".xmpp-client"), input, 0644)

	res := (&xmppClientImporter{}).TryImport()
	c.Assert(res, HasLen, 1)
}

func (s *XMPPClientXMPPSuite) Test_xmppClientImporter_TryImport_noFilesFound(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", origHome)
	}()
	os.Setenv("HOME", dir)

	res := (&xmppClientImporter{}).TryImport()
	c.Assert(res, HasLen, 0)
}

func (s *XMPPClientXMPPSuite) Test_xmppClientImporter_TryImport_failsOnBadFile(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", origHome)
	}()
	os.Setenv("HOME", dir)

	_ = ioutil.WriteFile(filepath.Join(dir, ".xmpp-client"), []byte("{"), 0644)

	res := (&xmppClientImporter{}).TryImport()
	c.Assert(res, HasLen, 0)
}

func (s *XMPPClientXMPPSuite) Test_xmppClientImporter_importFrom_failsOnReadingFile(c *C) {
	res, ok := (&xmppClientImporter{}).importFrom("bla-file-that-hopefully-doesnt-exists")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *XMPPClientXMPPSuite) Test_xmppClientImporter_importFrom_failsOnAddingAccount(c *C) {
	origAddNewAccount := addNewAccount
	defer func() {
		addNewAccount = origAddNewAccount
	}()

	addNewAccount = func(a *config.ApplicationConfig) (*config.Account, error) {
		return nil, errors.New("marker return")
	}

	res, ok := (&xmppClientImporter{}).importFrom(testResourceFilename("xmpp_client_test_conf.json"))
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

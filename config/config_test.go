package config

import (
	"encoding/json"
	"net"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ConfigXmppSuite struct{}

var _ = Suite(&ConfigXmppSuite{})

func (s *ConfigXmppSuite) TestDetectTor(c *C) {
	scannedForTor = false

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	c.Assert(err, IsNil)

	_, port, err := net.SplitHostPort(ln.Addr().String())
	c.Assert(err, IsNil)

	torPorts = []string{port}
	torAddress := detectTor()
	c.Assert(torAddress, Equals, ln.Addr().String())

	ln.Close()

	newAddr := detectTor()
	c.Assert(newAddr, Equals, torAddress)
}

func (s *ConfigXmppSuite) TestDetectTorConnectionRefused(c *C) {
	scannedForTor = false

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	c.Assert(err, IsNil)

	_, port, err := net.SplitHostPort(ln.Addr().String())
	c.Assert(err, IsNil)

	ln.Close()

	torPorts = []string{port}
	torAddress := detectTor()
	c.Assert(torAddress, Equals, "")
}

func (s *ConfigXmppSuite) TestParseYes(c *C) {
	c.Assert(ParseYes("Y"), Equals, true)
	c.Assert(ParseYes("y"), Equals, true)
	c.Assert(ParseYes("YES"), Equals, true)
	c.Assert(ParseYes("yes"), Equals, true)
	c.Assert(ParseYes("Yes"), Equals, true)
	c.Assert(ParseYes("anything"), Equals, false)
}

func (s *ConfigXmppSuite) TestSerializeMultiAccountConfig(c *C) {
	expected := `{
	"Accounts": [
		{
			"Account": "bob@riseup.net",
			"PrivateKey": null,
			"KnownFingerprints": null,
			"Bell": false,
			"HideStatusUpdates": false,
			"UseTor": true,
			"OTRAutoTearDown": false,
			"OTRAutoAppendTag": false,
			"OTRAutoStartSession": false,
			"AlwaysEncrypt": true
		},
		{
			"Account": "bob@riseup.net",
			"PrivateKey": null,
			"KnownFingerprints": null,
			"Bell": false,
			"HideStatusUpdates": false,
			"UseTor": false,
			"OTRAutoTearDown": false,
			"OTRAutoAppendTag": false,
			"OTRAutoStartSession": false
		}
	]
}`

	conf := MultiAccount{
		Accounts: []Config{
			Config{
				Account:       "bob@riseup.net",
				UseTor:        true,
				AlwaysEncrypt: true,
			},
			Config{
				Account: "bob@riseup.net",
			},
		},
	}

	contents, err := json.MarshalIndent(conf, "", "\t")
	c.Assert(err, IsNil)
	c.Assert(string(contents), Equals, expected)
}

func (s *ConfigXmppSuite) TestParseMultiAccount(c *C) {
	multiConf := &MultiAccount{
		Accounts: []Config{
			Config{
				Account:           "alice@riseup.net",
				HideStatusUpdates: true,
				OTRAutoTearDown:   true,
				OTRAutoAppendTag:  true,
			},
		},
	}

	singleConf := &Config{
		Account:       "bob@riseup.net",
		Bell:          true,
		UseTor:        true,
		AlwaysEncrypt: true,
	}

	multiConfFile, _ := json.Marshal(multiConf)
	singleConfFile, _ := json.Marshal(singleConf)

	conf, err := parseMultiAccount([]byte(singleConfFile))
	c.Assert(err, IsNil)
	c.Assert(conf.Accounts[0], DeepEquals, *singleConf)

	conf, err = parseMultiAccount([]byte(multiConfFile))
	c.Assert(err, IsNil)
	c.Assert(conf, DeepEquals, multiConf)
}

func (s *ConfigXmppSuite) TestFindConfigFile(c *C) {
	conf := FindConfigFile()
	c.Assert(conf, Equals, os.Getenv("HOME")+"/.config/coyim/accounts.json")
}

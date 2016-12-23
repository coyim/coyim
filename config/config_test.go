package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type ConfigXMPPSuite struct{}

var _ = Suite(&ConfigXMPPSuite{})

func (s *ConfigXMPPSuite) TestParseYes(c *C) {
	c.Assert(ParseYes("Y"), Equals, true)
	c.Assert(ParseYes("y"), Equals, true)
	c.Assert(ParseYes("YES"), Equals, true)
	c.Assert(ParseYes("yes"), Equals, true)
	c.Assert(ParseYes("Yes"), Equals, true)
	c.Assert(ParseYes("anything"), Equals, false)
}

func (s *ConfigXMPPSuite) TestSerializeAccountsConfig(c *C) {
	expected := `{
	"Accounts": [
		{
			"Account": "bob@riseup.net",
			"Peers": [
				{
					"UserID": "bob@riseup.net",
					"Nickname": "boby",
					"Fingerprints": null
				}
			],
			"HideStatusUpdates": false,
			"OTRAutoTearDown": false,
			"OTRAutoAppendTag": false,
			"OTRAutoStartSession": false,
			"AlwaysEncrypt": true,
			"ConnectAutomatically": false
		},
		{
			"Account": "bob@riseup.net",
			"Peers": null,
			"HideStatusUpdates": false,
			"OTRAutoTearDown": false,
			"OTRAutoAppendTag": false,
			"OTRAutoStartSession": false,
			"ConnectAutomatically": false
		}
	],
	"Bell": false,
	"ConnectAutomatically": false,
	"Display": {
		"MergeAccounts": false,
		"ShowOnlyOnline": false,
		"HideFeedbackBar": false,
		"ShowOnlyConfirmed": false,
		"SortByStatus": false
	},
	"AdvancedOptions": false,
	"UniqueConfigurationID": ""
}`

	conf := ApplicationConfig{
		Accounts: []*Account{
			&Account{
				Account: "bob@riseup.net",
				Peers: []*Peer{
					&Peer{
						UserID:   "bob@riseup.net",
						Nickname: "boby",
					},
				},
				AlwaysEncrypt: true,
			},
			&Account{
				Account: "bob@riseup.net",
			},
		},
	}

	contents, err := json.MarshalIndent(conf, "", "\t")
	c.Assert(err, IsNil)
	c.Assert(string(contents), Equals, expected)
}

func (s *ConfigXMPPSuite) TestFindConfigFile(c *C) {
	conf := findConfigFile("")
	if strings.HasSuffix(conf, ".enc") {
		c.Assert(conf, Equals, os.Getenv("HOME")+"/.config/coyim/accounts.json.enc")
	} else {
		c.Assert(conf, Equals, os.Getenv("HOME")+"/.config/coyim/accounts.json")
	}
}

package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseYes(t *testing.T) {
	if ok := parseYes("Y"); !ok {
		t.Errorf("parsed Y as %v", ok)
	}

	if ok := parseYes("y"); !ok {
		t.Errorf("parsed y as %v", ok)
	}

	if ok := parseYes("YES"); !ok {
		t.Errorf("parsed YES as %v", ok)
	}

	if ok := parseYes("yes"); !ok {
		t.Errorf("parsed yes as %v", ok)
	}

	if ok := parseYes("Yes"); !ok {
		t.Errorf("parsed yes as %v", ok)
	}

	if ok := parseYes("anything"); ok {
		t.Errorf("parsed something else as %v", ok)
	}
}

func TestSerializeMultiAccountConfig(t *testing.T) {
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

	conf := MultiAccountConfig{
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
	if err != nil {
		t.Errorf("failed to marshal config: %s", err)
	}

	if string(contents) != expected {
		t.Errorf("wrong serialized config: %s", string(contents))
	}
}

func TestParseMultiAccountConfig(t *testing.T) {
	multiConf := &MultiAccountConfig{
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

	c, err := parseMultiConfig([]byte(singleConfFile))
	if err != nil {
		t.Errorf("unexpected failure %s", err)
	}

	if !reflect.DeepEqual(c.Accounts[0], *singleConf) {
		t.Errorf("single account conf does not match %#v", c.Accounts[0])
	}

	c, err = parseMultiConfig([]byte(multiConfFile))
	if err != nil {
		t.Errorf("unexpected failure %s", err)
	}

	if !reflect.DeepEqual(c, multiConf) {
		t.Errorf("multi account conf does not match %#v", c)
	}
}

func TestFindConfigFile(t *testing.T) {
	if _, err := findConfigFile(""); err != errHomeDirNotSet {
		t.Errorf("unexpected error %s", err)
	}

	c, err := findConfigFile("/foo")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	if *c != "/foo/.xmpp-client" {
		t.Errorf("wrong config file path %s", *c)
	}
}

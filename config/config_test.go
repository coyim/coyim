package config

import (
	"encoding/json"
	"net"
	"reflect"
	"testing"
)

func TestDetectTor(t *testing.T) {
	scannedForTor = false

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	torPorts = []string{port}
	torAddress := detectTor()
	if torAddress != ln.Addr().String() {
		t.Fatalf("unexpected tor address %s", torAddress)
	}

	ln.Close()

	newAddr := detectTor()
	if newAddr != torAddress {
		t.Fatalf("unexpected tor address %s", torAddress)
	}
}

func TestDetectTorConnectionRefused(t *testing.T) {
	scannedForTor = false

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	ln.Close()

	torPorts = []string{port}
	torAddress := detectTor()
	if torAddress != "" {
		t.Fatalf("unexpected tor address %s", torAddress)
	}
}

func TestParseYes(t *testing.T) {
	if ok := ParseYes("Y"); !ok {
		t.Errorf("parsed Y as %v", ok)
	}

	if ok := ParseYes("y"); !ok {
		t.Errorf("parsed y as %v", ok)
	}

	if ok := ParseYes("YES"); !ok {
		t.Errorf("parsed YES as %v", ok)
	}

	if ok := ParseYes("yes"); !ok {
		t.Errorf("parsed yes as %v", ok)
	}

	if ok := ParseYes("Yes"); !ok {
		t.Errorf("parsed yes as %v", ok)
	}

	if ok := ParseYes("anything"); ok {
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
	if _, err := FindConfigFile(""); err != errHomeDirNotSet {
		t.Errorf("unexpected error %s", err)
	}

	c, err := FindConfigFile("/foo")
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	if *c != "/foo/.xmpp-client" {
		t.Errorf("wrong config file path %s", *c)
	}
}

package importer

import (
	"bufio"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/otr3"
)

// ImportKeysFromPidginStyle will try to read keys in Pidgin style from the given file
func ImportKeysFromPidginStyle(f string, protocolMatcher func(string) bool) (map[string][]byte, bool) {
	file, err := os.Open(f)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	acs, err := otr3.ImportKeys(file)
	if err != nil {
		return nil, false
	}

	res := make(map[string][]byte)
	for _, ac := range acs {
		if protocolMatcher(ac.Protocol) {
			res[strings.TrimSuffix(ac.Name, "/")] = ac.Key.Serialize()
		}
	}

	return res, true
}

// ImportFingerprintsFromPidginStyle will try to read fingerprints in Pidgin style from the given file
func ImportFingerprintsFromPidginStyle(f string, protocolMatcher func(string) bool) (map[string][]*config.KnownFingerprint, bool) {
	file, err := os.Open(f)
	if err != nil {
		return nil, false
	}

	defer file.Close()
	sc := bufio.NewScanner(file)
	result := make(map[string][]*config.KnownFingerprint)
	for sc.Scan() {
		ln := strings.Split(sc.Text(), "\t")
		if len(ln) < 4 {
			return nil, false
		}
		name := strings.TrimSuffix(ln[1], "/")
		if protocolMatcher(ln[2]) {
			vv, ok := result[name]
			if !ok {
				vv = make([]*config.KnownFingerprint, 0, 1)
			}

			fp, err := hex.DecodeString(ln[3])
			if err != nil {
				continue
			}

			result[name] = append(vv, &config.KnownFingerprint{
				UserID:      ln[0],
				Fingerprint: fp,
				Untrusted:   len(ln) < 5 || ln[4] != "verified",
			})
		}

	}

	return result, true
}

func importAccountsPidginStyle(f string) (map[string]*config.Account, bool) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, false
	}

	var a pidginAccountsXML
	err = xml.Unmarshal(content, &a)
	if err != nil {
		return nil, false
	}

	res := make(map[string]*config.Account)
	for _, ac := range a.Accounts {
		if ac.Protocol == "prpl-jabber" {
			nm := data.ParseJID(ac.Name).EnsureNoResource().Representation()
			a := &config.Account{}
			a.Account = nm
			a.Password = ac.Password
			settings := ac.settingsAsMap()
			a.Port = parseIntOr(settings["port"], 5222)

			a.Proxies = make([]string, 0)
			for _, px := range ac.Proxy {
				if px.Type == "tor" {
					a.Proxies = append(a.Proxies, "tor-auto://")
				}
				a.Proxies = append(a.Proxies,
					composeProxyString(px.Type, px.Username, px.Password, px.Host, strconv.Itoa(px.Port)),
				)
			}

			if settings["connect_server"] != "" {
				a.Server = settings["connect_server"]
			}

			res[nm] = a
		}
	}

	return res, true
}

func importPeerPrefsPidginStyle(f string) (map[string]map[string]*pidginOTRSettings, bool) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, false
	}

	var a pidginBlistXML
	err = xml.Unmarshal(content, &a)
	if err != nil {
		return nil, false
	}

	res := make(map[string]map[string]*pidginOTRSettings)

	for _, p := range a.Peers {
		if p.Protocol == "prpl-jabber" {
			haveOTR := false
			settings := &pidginOTRSettings{}
			for _, s := range p.Settings {
				switch s.Name {
				case "OTR/enabled":
					haveOTR = true
					settings.enabled = s.Value == "1"
				case "OTR/automatic":
					haveOTR = true
					settings.automatic = s.Value == "1"
				case "OTR/avoidloggingotr":
					haveOTR = true
					settings.avoidLoggingOTR = s.Value == "1"
				case "OTR/onlyprivate":
					haveOTR = true
					settings.onlyPrivate = s.Value == "1"
				}
			}
			if haveOTR {
				pp := data.ParseJID(p.Account).EnsureNoResource().Representation()
				pp2 := data.ParseJID(p.Name).EnsureNoResource().Representation()
				getOrMake(res, pp)[pp2] = settings

			}
		}
	}

	return res, true
}

func importGlobalPrefsPidginStyle(f string) (*pidginOTRSettings, bool) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, false
	}

	var a pidginPrefsXML
	err = xml.Unmarshal(content, &a)
	if err != nil {
		return nil, false
	}

	settings := pidginOTRSettings{}

	have := false
	if res, ok := a.lookup("OTR", "enabled"); ok {
		have = true
		settings.enabled = res.Value == "1"
	}
	if res, ok := a.lookup("OTR", "automatic"); ok {
		have = true
		settings.automatic = res.Value == "1"
	}
	if res, ok := a.lookup("OTR", "onlyprivate"); ok {
		have = true
		settings.onlyPrivate = res.Value == "1"
	}
	if res, ok := a.lookup("OTR", "avoidloggingotr"); ok {
		have = true
		settings.avoidLoggingOTR = res.Value == "1"
	}

	return &settings, have
}

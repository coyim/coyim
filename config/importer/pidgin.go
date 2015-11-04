package importer

import (
	"bufio"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/otr3"
)

// In $HOME or $APPDATA
const pidginConfigDir = ".purple"

// In Pidgin config directory
const pidginAccountsFile = "accounts.xml"
const pidginBuddyFile = "blist.xml"
const pidginPrefsFile = "prefs.xml"
const pidginOtrDataFingerprintsFile = "otr.fingerprints"
const pidginOtrDataKeyFile = "otr.private_key"

type pidginImporter struct{}

type pidginAccountsXML struct {
	Accounts []pidginAccountXML `xml:"account"`
}

type pidginAccountXML struct {
	Protocol string             `xml:"protocol"`
	Name     string             `xml:"name"`
	Password string             `xml:"password"`
	Proxy    []pidginProxyXML   `xml:"proxy"`
	Settings []pidginSettingXML `xml:"settings>setting"`
}

type pidginProxyXML struct {
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Port     int    `xml:"port"`
	Username string `xml:"username"`
	Password string `xml:"password"`
}

type pidginSettingXML struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

func (pax *pidginAccountXML) settingsAsMap() map[string]string {
	res := make(map[string]string)

	for _, s := range pax.Settings {
		res[s.Name] = s.Value
	}

	return res
}

func (p *pidginImporter) importKeysFrom(f string) (map[string][]byte, bool) {
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
		if ac.Protocol == "prpl-jabber" {
			res[strings.TrimSuffix(ac.Name, "/")] = ac.Key.Serialize()
		}
	}

	return res, true
}

func (p *pidginImporter) importFingerprintsFrom(f string) (map[string][]*config.KnownFingerprint, bool) {
	file, err := os.Open(f)
	if err != nil {
		return nil, false
	}

	defer file.Close()
	sc := bufio.NewScanner(file)
	result := make(map[string][]*config.KnownFingerprint)
	for sc.Scan() {
		ln := strings.Split(sc.Text(), "\t")
		name := strings.TrimSuffix(ln[1], "/")
		if ln[2] == "prpl-jabber" {
			vv, ok := result[name]
			if !ok {
				vv = make([]*config.KnownFingerprint, 0, 1)
			}
			result[name] = append(vv, &config.KnownFingerprint{
				UserID:         ln[0],
				FingerprintHex: ln[3],
				Untrusted:      len(ln) < 5 || ln[4] != "verified",
			})
		}

	}

	return result, true
}

func parseIntOr(s string, def int) int {
	if ret, e := strconv.Atoi(s); e == nil {
		return ret
	}
	return def
}

type pidginPrefsXML struct {
	Prefs []pidginPrefXML `xml:"pref"`
}

type pidginPrefXML struct {
	Name  string          `xml:"name,attr"`
	Type  string          `xml:"type,attr"`
	Value string          `xml:"value,attr"`
	Prefs []pidginPrefXML `xml:"pref"`
}

type pidginBlistXML struct {
	Peers []pidginPeerXML `xml:"blist>group>contact>buddy"`
}

type pidginPeerXML struct {
	Account  string             `xml:"account,attr"`
	Protocol string             `xml:"proto,attr"`
	Name     string             `xml:"name"`
	Settings []pidginSettingXML `xml:"setting"`
}

func (p *pidginPrefsXML) lookup(path ...string) (*pidginPrefXML, bool) {
	for _, pp := range p.Prefs {
		if pp.Name == path[0] {
			return pp.lookup(path[1:]...)
		}
	}
	return nil, false
}

func (p *pidginPrefXML) lookup(path ...string) (*pidginPrefXML, bool) {
	if p == nil {
		return nil, false
	}
	if len(path) == 0 {
		return p, true
	}
	for _, pp := range p.Prefs {
		if pp.Name == path[0] {
			return pp.lookup(path[1:]...)
		}
	}
	return nil, false
}

type pidginOTRSettings struct {
	enabled         bool
	automatic       bool
	onlyPrivate     bool
	avoidLoggingOTR bool
}

// only private means:
// prefsp->policy = OTRL_POLICY_ALWAYS;
// } else {
// prefsp->policy = OTRL_POLICY_OPPORTUNISTIC;
// automatic means: ALWAYS or OPP
//    otherwise it's MANUAL

// #define OTRL_POLICY_OPPORTUNISTIC \
// 	    ( OTRL_POLICY_ALLOW_V2 | \
// 	    OTRL_POLICY_ALLOW_V3 | \
// 	    OTRL_POLICY_SEND_WHITESPACE_TAG | \
// 	    OTRL_POLICY_WHITESPACE_START_AKE | \
// 	    OTRL_POLICY_ERROR_START_AKE )
// #define OTRL_POLICY_MANUAL \
// 	    ( OTRL_POLICY_ALLOW_V2 | \
// 	    OTRL_POLICY_ALLOW_V3)
// #define OTRL_POLICY_ALWAYS \
// 	    ( OTRL_POLICY_ALLOW_V2 | \
// 	    OTRL_POLICY_ALLOW_V3 | \
// 	    OTRL_POLICY_REQUIRE_ENCRYPTION | \
// 	    OTRL_POLICY_WHITESPACE_START_AKE | \
// 	    OTRL_POLICY_ERROR_START_AKE )

func getOrMake(m map[string]map[string]*pidginOTRSettings, nm string) map[string]*pidginOTRSettings {
	v, ok := m[nm]
	if !ok {
		v = make(map[string]*pidginOTRSettings)
		m[nm] = v
	}
	return v
}

func (p *pidginImporter) importPeerPrefs(f string) (map[string]map[string]*pidginOTRSettings, bool) {
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
				getOrMake(res, strings.TrimSuffix(p.Account, "/"))[strings.TrimSuffix(p.Name, "/")] = settings

			}
		}
	}

	return res, true
}

func (p *pidginImporter) importGlobalPrefs(f string) (*pidginOTRSettings, bool) {
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

	if res, ok := a.lookup("OTR", "enabled"); ok {
		settings.enabled = res.Value == "1"
	}
	if res, ok := a.lookup("OTR", "automatic"); ok {
		settings.automatic = res.Value == "1"
	}
	if res, ok := a.lookup("OTR", "onlyprivate"); ok {
		settings.onlyPrivate = res.Value == "1"
	}
	if res, ok := a.lookup("OTR", "avoidloggingotr"); ok {
		settings.avoidLoggingOTR = res.Value == "1"
	}

	return &settings, true
}

func (p *pidginImporter) importAccounts(f string) (map[string]*config.Account, bool) {
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
			nm := strings.TrimSuffix(ac.Name, "/")
			a := &config.Account{}
			a.Account = nm
			a.Password = ac.Password
			settings := ac.settingsAsMap()
			a.Port = parseIntOr(settings["port"], 5222)

			a.Proxies = make([]string, 0)
			for _, px := range ac.Proxy {
				if px.Type == "tor" {
					a.RequireTor = true
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

func (p *pidginImporter) importAllFrom(accountsFile, prefsFile, blistFile, keyFile, fprFile string) (*config.Accounts, bool) {
	accounts, ok1 := p.importAccounts(accountsFile)
	globalPrefs, ok2 := p.importGlobalPrefs(prefsFile)
	peerPrefs, ok3 := p.importPeerPrefs(blistFile)
	keys, ok4 := p.importKeysFrom(keyFile)
	fprs, ok5 := p.importFingerprintsFrom(fprFile)

	if !ok1 {
		return nil, false
	}

	res := &config.Accounts{}
	for name, ac := range accounts {
		res.Add(ac)
		if ok2 {
			if globalPrefs.enabled {
				if globalPrefs.onlyPrivate {
					ac.AlwaysEncrypt = true
					ac.OTRAutoStartSession = true
				} else if globalPrefs.automatic {
					ac.OTRAutoStartSession = true
					ac.OTRAutoAppendTag = true
				}
			} else {
				ac.AlwaysEncrypt = false
			}
		}
		if ok3 {
			if ss, ok := peerPrefs[name]; ok {
				for p, sp := range ss {
					if sp.enabled {
						if sp.onlyPrivate {
							ac.AlwaysEncryptWith = append(ac.AlwaysEncryptWith, p)
						}
					} else {
						ac.DontEncryptWith = append(ac.DontEncryptWith, p)
					}
				}
			}
		}
		if ok4 {
			if kk, ok := keys[name]; ok {
				ac.PrivateKey = kk
			}
		}
		if ok5 {
			if fprs, ok := fprs[name]; ok {
				ac.KnownFingerprints = make([]config.KnownFingerprint, len(fprs))
				sort.Sort(config.ByNaturalOrder(fprs))
				for ix, fpr := range fprs {
					ac.KnownFingerprints[ix] = *fpr
				}
			}
		}
	}

	sort.Sort(config.ByAccountNameAlphabetic(res.Accounts))

	return res, true
}

func (p *pidginImporter) findDir() (string, bool) {
	if fi, err := os.Stat(config.WithHome(filepath.Join(pidginConfigDir, pidginAccountsFile))); err == nil && !fi.IsDir() {
		return config.WithHome(pidginConfigDir), true
	}

	if config.IsWindows() {
		app := filepath.Join(config.SystemConfigDir(), pidginConfigDir)

		if fi, err := os.Stat(filepath.Join(app, pidginAccountsFile)); err == nil && !fi.IsDir() {
			return app, true
		}
	}

	return "", false
}

func (p *pidginImporter) composeFileNamesFrom(dir string) (accountsFile, prefsFile, blistFile, keyFile, fprFile string) {
	return filepath.Join(dir, pidginAccountsFile), filepath.Join(dir, pidginPrefsFile), filepath.Join(dir, pidginBuddyFile), filepath.Join(dir, pidginOtrDataKeyFile), filepath.Join(dir, pidginOtrDataFingerprintsFile)
}

func (p *pidginImporter) TryImport() []*config.Accounts {
	var res []*config.Accounts

	dd, ok := p.findDir()
	if ok {
		ac, ok := p.importAllFrom(p.composeFileNamesFrom(dd))
		if ok {
			res = append(res, ac)
		}
	}

	return res
}

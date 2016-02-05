package importer

import (
	"bufio"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hydrogen18/stalecucumber"
	"github.com/twstrike/coyim/config"
)

const gajimOtrDataKeyExtension = ".key3"
const gajimOtrDataFingerprintsExtension = ".fpr"

type gajimImporter struct{}

func getFilesMatching(dir, ext string) []string {
	var files []string
	fs, err := ioutil.ReadDir(dir)
	if err == nil {
		for _, fi := range fs {
			if !fi.IsDir() && path.Ext(fi.Name()) == ext {
				files = append(files, filepath.Join(dir, fi.Name()))
			}
		}
	}
	return files
}

func (g *gajimImporter) findFiles() (configFile string, pluginFile string, keyFiles []string, fingerprintFiles []string) {
	var configRoot, dataRoot string

	if config.IsWindows() {
		configRoot = filepath.Join(config.SystemConfigDir(), "Gajim")
		dataRoot = configRoot
	} else {
		configRoot = filepath.Join(config.XdgConfigHome(), "gajim")
		dataRoot = filepath.Join(config.XdgDataDir(), "gajim")
	}

	configFile = filepath.Join(configRoot, "config")
	pluginFile = filepath.Join(configRoot, "pluginsconfig/gotr")

	fingerprintFiles = getFilesMatching(dataRoot, gajimOtrDataFingerprintsExtension)
	keyFiles = getFilesMatching(dataRoot, gajimOtrDataKeyExtension)

	return
}

func (g *gajimImporter) importFingerprintsFrom(f string) (string, []*config.KnownFingerprint, bool) {
	file, err := os.Open(f)
	if err != nil {
		return "", nil, false
	}

	defer file.Close()
	sc := bufio.NewScanner(file)
	var result []*config.KnownFingerprint
	name := ""
	for sc.Scan() {
		ln := strings.Split(sc.Text(), "\t")
		name = ln[1]
		if ln[2] != "xmpp" {
			continue
		}

		fp, err := hex.DecodeString(ln[3])
		if err != nil {
			continue
		}

		result = append(result, &config.KnownFingerprint{
			UserID:      ln[0],
			Fingerprint: fp,
			Untrusted:   len(ln) < 5 || ln[4] != "verified",
		})

	}
	return name, result, true
}

func (g *gajimImporter) importKeyFrom(f string) (string, []byte, bool) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return "", nil, false
	}
	fn := strings.TrimSuffix(path.Base(f), path.Ext(f))
	return fn, content, true
}

type gajimOTRSettings struct {
	allowV1            bool
	allowV2            bool
	errorStartAke      bool
	requireEncryption  bool
	sendTag            bool
	whitespaceStartAke bool
}

type gajimAccountAndPeer struct {
	account string
	peer    string
}

func intoBool(val interface{}) bool {
	res, ok := val.(bool)
	return ok && res
}

func intoGajimOTRSettings(vv map[string]interface{}) gajimOTRSettings {
	return gajimOTRSettings{
		allowV1:            intoBool(vv["ALLOW_V1"]),
		allowV2:            intoBool(vv["ALLOW_V2"]),
		errorStartAke:      intoBool(vv["ERROR_START_AKE"]),
		requireEncryption:  intoBool(vv["REQUIRE_ENCRYPTION"]),
		sendTag:            intoBool(vv["SEND_TAG"]),
		whitespaceStartAke: intoBool(vv["WHITESPACE_START_AKE"]),
	}
}

func (g *gajimImporter) importOTRSettings(f string) (map[string]gajimOTRSettings, map[gajimAccountAndPeer]gajimOTRSettings, bool) {
	file, err := os.Open(f)
	if err != nil {
		return nil, nil, false
	}
	defer file.Close()
	res, err := stalecucumber.DictString(stalecucumber.Unpickle(file))
	if err != nil {
		return nil, nil, false
	}

	resAccount := make(map[string]gajimOTRSettings)
	resAccountToPeer := make(map[gajimAccountAndPeer]gajimOTRSettings)

	for k, v := range res {
		vv, err := stalecucumber.Dict(v, nil)
		if err == nil {
			for k2, v2 := range vv {
				vv2, _ := stalecucumber.DictString(v2, nil)
				settings := intoGajimOTRSettings(vv2)
				switch k2.(type) {
				case string:
					resAccountToPeer[gajimAccountAndPeer{k, k2.(string)}] = settings
				default:
					resAccount[k] = settings
				}
			}
		}
	}

	return resAccount, resAccountToPeer, true
}

func addNamedValue(m map[string]map[string]string, vals []string, prefix string) {
	newKey := strings.TrimPrefix(vals[0], prefix)
	sep := strings.LastIndex(newKey, ".")
	name := newKey[:sep]
	finalKeyName := newKey[sep+1:]
	hs, ok := m[name]
	if !ok {
		hs = make(map[string]string)
		m[name] = hs
	}
	hs[finalKeyName] = vals[1]
}

type gajimAccountInfo struct {
	accountNickName string
	password        string
	sslFingerprint  string
	hostname        string
	server          string
	name            string
	proxy           string
	port            string
}

func composeProxyStringFrom(proxy map[string]string) string {
	if proxy == nil {
		return ""
	}

	if proxy["useauth"] != "True" {
		proxy["user"] = ""
	}

	return composeProxyString(proxy["type"], proxy["user"], proxy["pass"], proxy["host"], proxy["port"])
}

func transformSettingsIntoAccount(user string, settings map[string]string, proxies map[string]map[string]string) gajimAccountInfo {
	res := gajimAccountInfo{
		accountNickName: user,
		password:        settings["password"],
		sslFingerprint:  strings.Replace(settings["ssl_fingerprint_sha1"], ":", "", -1),
		hostname:        settings["hostname"],
		server:          settings["custom_host"],
		name:            settings["name"],
		port:            settings["custom_port"],
	}

	if settings["use_custom_host"] != "True" {
		res.server = ""
	}

	p := settings["proxy"]

	if p != "" {
		res.proxy = composeProxyStringFrom(proxies[p])
	}

	return res
}

func (g *gajimImporter) importAccounts(f string) (map[string]gajimAccountInfo, bool) {
	file, err := os.Open(f)
	if err != nil {
		return nil, false
	}

	defer file.Close()
	sc := bufio.NewScanner(file)

	accountSettings := make(map[string]map[string]string)
	proxies := make(map[string]map[string]string)

	for sc.Scan() {
		val := sc.Text()
		ln := strings.SplitN(val, " = ", 2)
		if len(ln) == 2 {
			key := ln[0]
			switch {
			case strings.HasPrefix(key, "accounts."):
				addNamedValue(accountSettings, ln, "accounts.")
			case strings.HasPrefix(key, "proxies."):
				addNamedValue(proxies, ln, "proxies.")
			}
		}
	}

	accountInfo := make(map[string]gajimAccountInfo)

	for k, v := range accountSettings {
		accountInfo[k] = transformSettingsIntoAccount(k, v, proxies)
	}

	return accountInfo, true
}

func (g *gajimImporter) importAllFrom(configFile, gotrFile string, keyFiles []string, fprFiles []string) (*config.ApplicationConfig, bool) {
	accounts, ok1 := g.importAccounts(configFile)
	accountOTRSettings, accountAndPeerOTRSettings, _ := g.importOTRSettings(gotrFile)

	fprs := make(map[string][]*config.KnownFingerprint)
	for _, kk := range fprFiles {
		nm, res, ok := g.importFingerprintsFrom(kk)
		if ok {
			fprs[nm] = res
		}
	}
	keys := make(map[string][]byte)
	for _, kk := range keyFiles {
		nm, res, ok := g.importKeyFrom(kk)
		if ok {
			keys[nm] = res
		}
	}

	return mergeAll(accounts, accountOTRSettings, accountAndPeerOTRSettings, fprs, keys), ok1
}

func mergeAccountInformation(ac gajimAccountInfo, s gajimOTRSettings, s2 map[string]gajimOTRSettings, fprs []*config.KnownFingerprint, key []byte) *config.Account {
	res := &config.Account{}
	res.Password = ac.password

	// This is incorrect, since Gajim uses SHA-1 fingerprints, not SHA-256...
	//	res.ServerCertificateSHA256 = ac.sslFingerprint

	res.Account = ac.name + "@" + ac.hostname

	if ac.server != "" {
		res.Server = ac.server
	} else {
		res.Server = ac.hostname
	}

	if ac.port != "" {
		i, err := strconv.Atoi(ac.port)
		if err == nil {
			res.Port = i
		}
	}

	if ac.proxy != "" {
		res.Proxies = []string{ac.proxy}
		res.RequireTor = true
	}

	if strings.HasSuffix(res.Server, ".onion") {
		res.RequireTor = true
	}

	res.AlwaysEncrypt = s.requireEncryption
	res.OTRAutoAppendTag = s.sendTag
	res.OTRAutoStartSession = s.whitespaceStartAke
	res.OTRAutoTearDown = true

	if key != nil {
		res.PrivateKeys = [][]byte{key}
	}

	if fprs != nil {
		res.Peers = nil
		sort.Sort(config.LegacyByNaturalOrder(fprs))
		for _, kfpr := range fprs {
			fpr := res.EnsurePeer(kfpr.UserID).EnsureHasFingerprint(kfpr.Fingerprint)
			if !kfpr.Untrusted {
				fpr.Trusted = true
			}
		}
	}

	res.AlwaysEncryptWith = make([]string, 0)
	res.DontEncryptWith = make([]string, 0)
	for peer, settings := range s2 {
		if settings.requireEncryption {
			res.AlwaysEncryptWith = append(res.AlwaysEncryptWith, peer)
		} else {
			if !s.requireEncryption {
				res.DontEncryptWith = append(res.DontEncryptWith, peer)
			}
		}
	}
	sort.Sort(byAlpha(res.AlwaysEncryptWith))
	sort.Sort(byAlpha(res.DontEncryptWith))

	return res
}

type byAlpha []string

func (s byAlpha) Len() int { return len(s) }
func (s byAlpha) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s byAlpha) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func mergeAll(
	accounts map[string]gajimAccountInfo,
	accountOTRSettings map[string]gajimOTRSettings,
	accountAndPeerOTRSettings map[gajimAccountAndPeer]gajimOTRSettings,
	fprs map[string][]*config.KnownFingerprint,
	keys map[string][]byte,
) *config.ApplicationConfig {

	res := &config.ApplicationConfig{}

	for name, ac := range accounts {
		if name != "Local" {
			a := mergeAccountInformation(ac, accountOTRSettings[name], getPeerSettingsFor(name, accountAndPeerOTRSettings), fprs[name], keys[name])
			res.Add(a)
		}
	}

	sort.Sort(config.ByAccountNameAlphabetic(res.Accounts))

	return res
}

func getPeerSettingsFor(name string, s map[gajimAccountAndPeer]gajimOTRSettings) map[string]gajimOTRSettings {
	res := make(map[string]gajimOTRSettings)

	for k, v := range s {
		if k.account == name {
			res[k.peer] = v
		}
	}

	return res
}

func (g *gajimImporter) TryImport() []*config.ApplicationConfig {
	var res []*config.ApplicationConfig

	ac, ok := g.importAllFrom(g.findFiles())
	if ok {
		res = append(res, ac)
	}

	return res
}

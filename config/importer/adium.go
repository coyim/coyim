package importer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	plist "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/DHowett/go-plist"
	"github.com/twstrike/coyim/config"
)

// In $HOME
const adiumConfigDir = "Library/Application Support/Adium 2.0/Users/Default"

const adiumAccountMappingsFile = "Accounts.plist"

type adiumImporter struct{}

type adiumAccountMapping struct {
	objectID    string
	uid         string
	accountType string
}

func (p *adiumImporter) readAccountMappings(s string) (map[string]adiumAccountMapping, bool) {
	contents, _ := ioutil.ReadFile(s)
	var res map[string]interface{}
	plist.Unmarshal(contents, &res)

	result := make(map[string]adiumAccountMapping)

	acs := res["Accounts"].([]interface{})

	for _, v := range acs {
		vals := v.(map[string]interface{})

		var m adiumAccountMapping
		m.uid = vals["UID"].(string)
		m.objectID = vals["ObjectID"].(string)
		m.accountType = vals["Type"].(string)
		result[m.objectID] = m
	}

	return result, true
}

func (p *adiumImporter) protocolMatches(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), "libpurple-jabber")
}

func (p *adiumImporter) importKeysFrom(f string) (map[string][]byte, bool) {
	return ImportKeysFromPidginStyle(f, p.protocolMatches)
}

func (p *adiumImporter) importFingerprintsFrom(f string) (map[string][]*config.KnownFingerprint, bool) {
	return ImportFingerprintsFromPidginStyle(f, p.protocolMatches)
}

func (p *adiumImporter) importAccounts(f string) (map[string]*config.Account, bool) {
	return importAccountsPidginStyle(f)
}

func (p *adiumImporter) importPeerPrefs(f string) (map[string]map[string]*pidginOTRSettings, bool) {
	return importPeerPrefsPidginStyle(f)
}

func (p *adiumImporter) importGlobalPrefs(f string) (*pidginOTRSettings, bool) {
	return importGlobalPrefsPidginStyle(f)
}

func (p *adiumImporter) importAllFrom(accountMappingsFile, accountsFile, prefsFile, blistFile, keyFile, fprFile string) (*config.ApplicationConfig, bool) {
	accountMappings, ok0 := p.readAccountMappings(accountMappingsFile)
	accounts, ok1 := p.importAccounts(accountsFile)
	globalPrefs, ok2 := p.importGlobalPrefs(prefsFile)
	peerPrefs, ok3 := p.importPeerPrefs(blistFile)
	keysRaw, ok4 := p.importKeysFrom(keyFile)
	fprsRaw, ok5 := p.importFingerprintsFrom(fprFile)

	if !ok0 || !ok1 {
		return nil, false
	}

	keys := make(map[string][]byte)
	for k, v := range keysRaw {
		keys[accountMappings[k].uid] = v
	}

	fprs := make(map[string][]*config.KnownFingerprint)
	for k, v := range fprsRaw {
		fprs[accountMappings[k].uid] = v
	}

	res := &config.ApplicationConfig{}
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
				ac.PrivateKeys = [][]byte{kk}
			}
		}
		if ok5 {
			if fprs, ok := fprs[name]; ok {
				ac.Peers = nil
				sort.Sort(config.LegacyByNaturalOrder(fprs))
				for _, kfpr := range fprs {
					fpr, _ := ac.EnsurePeer(kfpr.UserID).EnsureHasFingerprint(kfpr.Fingerprint)
					if !kfpr.Untrusted {
						fpr.Trusted = true
					}
				}
			}
		}
	}

	sort.Sort(config.ByAccountNameAlphabetic(res.Accounts))

	return res, true
}

func (p *adiumImporter) findDir() (string, bool) {
	if fi, err := os.Stat(config.WithHome(filepath.Join(adiumConfigDir, adiumAccountMappingsFile))); err == nil && !fi.IsDir() {
		return config.WithHome(adiumConfigDir), true
	}

	return "", false
}

func (p *adiumImporter) composeFileNamesFrom(dir string) (accountMappingsFile, accountsFile, prefsFile, blistFile, keyFile, fprFile string) {
	return filepath.Join(dir, adiumAccountMappingsFile),
		filepath.Join(dir, "libpurple", pidginAccountsFile),
		filepath.Join(dir, "libpurple", pidginPrefsFile),
		filepath.Join(dir, "libpurple", pidginBuddyFile),
		filepath.Join(dir, pidginOtrDataKeyFile),
		filepath.Join(dir, pidginOtrDataFingerprintsFile)
}

func (p *adiumImporter) TryImport() []*config.ApplicationConfig {
	var res []*config.ApplicationConfig

	dd, ok := p.findDir()
	if ok {
		ac, ok := p.importAllFrom(p.composeFileNamesFrom(dd))
		if ok {
			res = append(res, ac)
		}
	}

	return res
}

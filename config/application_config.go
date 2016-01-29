package config

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
)

// ApplicationConfig contains the configuration for the application, including
// account information.
type ApplicationConfig struct {
	filename      string                `json:"-"`
	ShouldEncrypt bool                  `json:"-"`
	params        *EncryptionParameters `json:"-"`

	Accounts                      []*Account
	RawLogFile                    string   `json:",omitempty"`
	NotifyCommand                 []string `json:",omitempty"`
	IdleSecondsBeforeNotification int      `json:",omitempty"`
	Bell                          bool
	ConnectAutomatically          bool
	Display                       DisplayConfig `json:",omitempty"`
	AdvancedOptions               bool
}

var loadEntries []func(*ApplicationConfig)
var loadEntryLock = sync.Mutex{}

// WhenLoaded will ensure that the function f is not called until the configuration has been loaded
func (a *ApplicationConfig) WhenLoaded(f func(*ApplicationConfig)) {
	if a != nil {
		f(a)
		return
	}
	loadEntryLock.Lock()
	defer loadEntryLock.Unlock()

	loadEntries = append(loadEntries, f)
}

func (a *ApplicationConfig) accountLoaded() {
	loadEntryLock.Lock()
	defer loadEntryLock.Unlock()
	ourEntries := loadEntries
	loadEntries = []func(*ApplicationConfig){}
	for _, f := range ourEntries {
		go f(a)
	}
}

// LoadOrCreate will try to load the configuration from the given configuration file
// or from the standard configuration file. If no file exists or it is malformed,
// or it could not be decrypted, an error will be returned.
// However, the returned Accounts instance will always be usable
func LoadOrCreate(configFile string, ks KeySupplier) (a *ApplicationConfig, ok bool, e error) {
	a = new(ApplicationConfig)
	a.filename = findConfigFile(configFile)
	e = a.tryLoad(ks)
	ok = !(e == errNoPasswordSupplied || e == errDecryptionFailed)

	return
}

var (
	errInvalidConfigFile = errors.New("Failed to parse config file")
)

func (a *ApplicationConfig) tryLoad(ks KeySupplier) error {
	var contents []byte
	var err error

	contents, err = readFileOrTemporaryBackup(a.filename)
	if err != nil {
		return errInvalidConfigFile
	}
	_, err = parseEncryptedData(contents)
	switch err {
	case nil:
		a.ShouldEncrypt = true
		contents, a.params, err = decryptConfiguration(contents, ks)
		if err == errNoPasswordSupplied {
			return err
		} else if err == errDecryptionFailed {
			return err
		}
	case errDecryptionParamsEmpty:
		a.ShouldEncrypt = false
	default:
		return errInvalidConfigFile
	}

	if err = json.Unmarshal(contents, a); err != nil {
		return errInvalidConfigFile
	}

	if len(a.Accounts) == 0 {
		return errInvalidConfigFile
	}

	a.accountLoaded()

	return nil
}

// Add will add the account to the application configuration
func (a *ApplicationConfig) Add(ac *Account) {
	a.Accounts = append(a.Accounts, ac)
}

// Remove will update the accounts to exclude the account to remove, if
// it does exist
func (a *ApplicationConfig) Remove(toRemove *Account) {
	res := make([]*Account, len(a.Accounts)-1)
	found := false
	for i, ac := range a.Accounts {
		if ac.Is(toRemove.Account) {
			res = append(a.Accounts[:i], a.Accounts[i+1:]...)
			found = true
		}
	}
	if found {
		a.Accounts = res
	}
}

// AddNewAccount creates a new account and adds it to the list of accounts
func (a *ApplicationConfig) AddNewAccount() (ac *Account, err error) {
	ac, err = NewAccount()
	if err == nil {
		a.Add(ac)
	}
	return
}

// GetAccount will return the account with the given JID or not OK if it doesn't exist
func (a *ApplicationConfig) GetAccount(jid string) (*Account, bool) {
	for _, aa := range a.Accounts {
		if aa.Is(jid) {
			return aa, true
		}
	}
	return nil, false
}

// Save will save the application configuration
func (a *ApplicationConfig) Save(ks KeySupplier) error {
	contents, err := a.serialize()
	if err != nil {
		return err
	}

	if a.ShouldEncrypt && !strings.HasSuffix(a.filename, encryptedFileEnding) {
		a.filename = a.filename + encryptedFileEnding
	}

	if a.ShouldEncrypt {
		if a.params == nil {
			ps := newEncryptionParameters()
			a.params = &ps
		} else {
			a.params.regenerateNonce()
		}

		contents, err = encryptConfiguration(string(contents), a.params, ks)
		if err != nil {
			return err
		}
	}

	return safeWrite(a.filename, contents, 0600)
}

// UpdateToLatestVersion will run through all accounts and update their configuration to latest version
// for cases where we have changed the configuration format.
// It returns true if any changes were made
func (a *ApplicationConfig) UpdateToLatestVersion() bool {
	res := false

	for _, acc := range a.Accounts {
		res = acc.updateToLatestVersion() || res
	}

	return res
}

func (a *ApplicationConfig) serialize() ([]byte, error) {
	return json.MarshalIndent(a, "", "\t")
}

// ByAccountNameAlphabetic sorts the accounts based on their account names
type ByAccountNameAlphabetic []*Account

func (s ByAccountNameAlphabetic) Len() int { return len(s) }
func (s ByAccountNameAlphabetic) Less(i, j int) bool {
	return s[i].Account < s[j].Account
}
func (s ByAccountNameAlphabetic) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

package config

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
)

// Accounts contains the configuration for several accounts
type Accounts struct {
	filename      string                `json:"-"`
	ShouldEncrypt bool                  `json:"-"`
	params        *EncryptionParameters `json:"-"`

	Accounts                      []*Account
	RawLogFile                    string   `json:",omitempty"`
	NotifyCommand                 []string `json:",omitempty"`
	IdleSecondsBeforeNotification int      `json:",omitempty"`
	Bell                          bool
	MergeAccounts                 bool
	ShowOnlyOnline                bool
	ConnectAutomatically          bool
}

var loadEntries []func(*Accounts)
var loadEntryLock = sync.Mutex{}

// WhenLoaded will ensure that the function f is not called until the configuration has been loaded
func (a *Accounts) WhenLoaded(f func(*Accounts)) {
	if a != nil {
		f(a)
		return
	}
	loadEntryLock.Lock()
	defer loadEntryLock.Unlock()

	loadEntries = append(loadEntries, f)
}

func (a *Accounts) accountLoaded() {
	loadEntryLock.Lock()
	defer loadEntryLock.Unlock()
	ourEntries := loadEntries
	loadEntries = []func(*Accounts){}
	for _, f := range ourEntries {
		go f(a)
	}
}

// LoadOrCreate will try to load the configuration from the given configuration file
// or from the standard configuration file. If no file exists or it is malformed,
// or it could not be decrypted, an error will be returned.
// However, the returned Accounts instance will always be usable
func LoadOrCreate(configFile string, ks KeySupplier) (a *Accounts, ok bool, e error) {
	a = new(Accounts)
	a.filename = findConfigFile(configFile)
	e = a.tryLoad(ks)
	ok = !(e == errNoPasswordSupplied || e == errDecryptionFailed)

	return
}

var (
	errInvalidConfigFile = errors.New("Failed to parse config file")
)

func (a *Accounts) tryLoad(ks KeySupplier) error {
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

	for _, c := range a.Accounts {
		parseFingerprints(c)
	}

	a.accountLoaded()

	return nil
}

// Add will add the account to the configuration
func (a *Accounts) Add(ac *Account) {
	a.Accounts = append(a.Accounts, ac)
}

// AddNewAccount creates a new account and adds it to the list of accounts
func (a *Accounts) AddNewAccount() (ac *Account, err error) {
	ac, err = NewAccount()
	if err == nil {
		a.Add(ac)
	}
	return
}

// GetAccount will return the account with the given JID or not OK if it doesn't exist
func (a *Accounts) GetAccount(jid string) (*Account, bool) {
	for _, aa := range a.Accounts {
		if aa.Is(jid) {
			return aa, true
		}
	}
	return nil, false
}

// Save will save the account configuration
func (a *Accounts) Save(ks KeySupplier) error {
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

func (a *Accounts) serialize() ([]byte, error) {
	for _, account := range a.Accounts {
		account.serializeFingerprints()
	}

	return json.MarshalIndent(a, "", "\t")
}

// ByAccountNameAlphabetic sorts the accounts based on their account names
type ByAccountNameAlphabetic []*Account

func (s ByAccountNameAlphabetic) Len() int { return len(s) }
func (s ByAccountNameAlphabetic) Less(i, j int) bool {
	return s[i].Account < s[j].Account
}
func (s ByAccountNameAlphabetic) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

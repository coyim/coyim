package config

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"strings"

	"github.com/twstrike/otr3"
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
}

// LoadOrCreate will try to load the configuration from the given configuration file
// or from the standard configuration file. If no file exists or it is malformed,
// or it could not be decrypted, an error will be returned.
// However, the returned Accounts instance will always be usable
func LoadOrCreate(configFile string, ks KeySupplier) (a *Accounts, ok bool, e error) {
	shouldEncrypt := false
	if len(configFile) == 0 {
		configFile, shouldEncrypt = findConfigFile()
	}

	a = new(Accounts)
	a.filename = configFile
	a.ShouldEncrypt = shouldEncrypt
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

	if a.ShouldEncrypt {
		contents2, err2 := readFileOrTemporaryBackup(a.filename)
		if err2 != nil {
			err = err2
		} else {
			contents, a.params, err = decryptConfiguration(contents2, ks)

			if err == errNoPasswordSupplied {
				return err
			} else if err == errDecryptionFailed {
				return err
			}
		}
	} else {
		contents, err = readFileOrTemporaryBackup(a.filename)
	}

	if err != nil {
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

	return nil
}

// NewAccount creates a new account
func NewAccount() (*Account, error) {
	var torProxy []string
	torAddress := detectTor()

	if len(torAddress) != 0 {
		torProxy = []string{newTorProxy(torAddress)}
	}

	var priv otr3.PrivateKey

	err := priv.Generate(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Account{
		Proxies:    torProxy,
		RequireTor: torProxy != nil,

		PrivateKey:          priv.Serialize(),
		AlwaysEncrypt:       true,
		OTRAutoStartSession: true,
		OTRAutoTearDown:     true, //See #48
	}, nil
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

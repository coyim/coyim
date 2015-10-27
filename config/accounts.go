package config

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/twstrike/otr3"
)

// Accounts contains the configuration for several accounts
type Accounts struct {
	filename      string `json:"-"`
	ShouldEncrypt bool   `json:"-"`

	Accounts                      []*Account
	RawLogFile                    string   `json:",omitempty"`
	NotifyCommand                 []string `json:",omitempty"`
	IdleSecondsBeforeNotification int      `json:",omitempty"`
	Bell                          bool
	MergeAccounts                 bool
}

// LoadOrCreate will try to load the configuration from the given configuration file
// or from the standard configuration file. If no file exists or it is malformed, an error will
// be returned. However, the returned Accounts instance will always be usable
func LoadOrCreate(configFile string, ks KeySupplier) (a *Accounts, e error) {
	shouldEncrypt := false
	if len(configFile) == 0 {
		configFile, shouldEncrypt = findConfigFile()
	}

	a = new(Accounts)
	a.filename = configFile
	a.ShouldEncrypt = shouldEncrypt
	e = a.tryLoad(ks)

	return
}

var (
	errInvalidConfigFile = errors.New("Failed to parse config file")
)

func (a *Accounts) tryLoad(ks KeySupplier) error {
	var contents []byte
	var err error

	if a.ShouldEncrypt {
		contents2, err2 := ioutil.ReadFile(a.filename)
		if err2 != nil {
			err = err2
		} else {
			contents, err = decryptConfiguration(contents2, ks)
		}
	} else {
		contents, err = ioutil.ReadFile(a.filename)
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
func NewAccount() *Account {
	var torProxy []string
	torAddress := detectTor()

	if len(torAddress) != 0 {
		torProxy = []string{newTorProxy(torAddress)}
	}

	var priv otr3.PrivateKey

	//TODO: error
	priv.Generate(rand.Reader)

	return &Account{
		Proxies:    torProxy,
		RequireTor: torProxy != nil,

		PrivateKey:          priv.Serialize(),
		AlwaysEncrypt:       true,
		OTRAutoStartSession: true,
		OTRAutoTearDown:     true, //See #48
	}
}

// Add will add the account to the configuration
func (a *Accounts) Add(ac *Account) {
	a.Accounts = append(a.Accounts, ac)
}

// AddNewAccount creates a new account and adds it to the list of accounts
func (a *Accounts) AddNewAccount() *Account {
	ac := NewAccount()
	a.Add(ac)
	return ac
}

// Save will save the account configuration
func (a *Accounts) Save(ks KeySupplier) error {
	contents, err := a.serialize()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(a.filename, contents, 0600)
}

func (a *Accounts) serialize() ([]byte, error) {
	for _, account := range a.Accounts {
		account.serializeFingerprints()
	}

	return json.MarshalIndent(a, "", "\t")
}

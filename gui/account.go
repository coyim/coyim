package gui

import (
	"errors"

	"github.com/gotk3/gotk3/glib"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
)

// someone who knows how to persist account configuration
type configManager interface {
	Save() error
}

type Account struct {
	ConnectedSignal    *glib.Signal
	DisconnectedSignal *glib.Signal

	configManager
	*config.Config
	*session.Session
}

func (acc *Account) Connected() bool {
	return acc.ConnStatus == session.CONNECTED
}

var (
	errFingerprintAlreadyAuthorized = errors.New(i18n.Local("the fingerprint is already authorized"))
)

func (acc *Account) AuthorizeFingerprint(uid string, fingerprint []byte) error {
	existing := acc.UserIdForFingerprint(fingerprint)
	if len(existing) != 0 {
		return errFingerprintAlreadyAuthorized
	}

	acc.KnownFingerprints = append(acc.KnownFingerprints, config.KnownFingerprint{
		Fingerprint: fingerprint, UserId: uid,
	})

	return nil
}

func BuildAccountsFrom(multiAcc *config.MultiAccount, manager configManager, u *gtkUI) []*Account {
	accounts := make([]*Account, len(multiAcc.Accounts))

	for i := range multiAcc.Accounts {
		conf := &multiAcc.Accounts[i]
		u.ensureConfigHasKey(conf)

		account := newAccount(conf, manager)
		account.Session.SessionEventHandler = u
		accounts[i] = account
	}

	return accounts
}

func newAccount(conf *config.Config, m configManager) *Account {
	id := conf.ID()
	c, _ := glib.SignalNew(signalName(id, "connected"))
	d, _ := glib.SignalNew(signalName(id, "disconnected"))

	a := &Account{
		Config:        conf,
		Session:       session.NewSession(conf),
		configManager: m,

		ConnectedSignal:    c,
		DisconnectedSignal: d,
	}
	a.Session.Account = a

	return a
}

func signalName(id, signal string) string {
	return "coyim-account-" + signal + "-" + id
}

func (u *gtkUI) showAddAccountWindow() {
	account := newAccount(config.NewConfig(), u.configFileManager)
	accountDialog(account, func() error {
		err := u.configFileManager.Add(*account.Config)
		if err != nil {
			return err
		}

		return u.SaveConfig()
	})
}

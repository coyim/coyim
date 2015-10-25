package gui

import (
	"errors"

	"github.com/gotk3/gotk3/glib"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
)

type account struct {
	connectedSignal    *glib.Signal
	disconnectedSignal *glib.Signal

	session *session.Session
}

func (acc *account) connected() bool {
	return acc.session.ConnStatus == session.CONNECTED
}

var (
	errFingerprintAlreadyAuthorized = errors.New(i18n.Local("the fingerprint is already authorized"))
)

// TODO: this functionality is duplicated
func (acc *account) authorizeFingerprint(uid string, fingerprint []byte) error {
	a := acc.session.CurrentAccount
	existing := a.UserIDForFingerprint(fingerprint)
	if len(existing) != 0 {
		return errFingerprintAlreadyAuthorized
	}

	a.AddFingerprint(fingerprint, uid)

	return nil
}

func (u *gtkUI) buildAccounts() []*account {
	accounts := make([]*account, len(u.config.Accounts))

	for i, accountConf := range u.config.Accounts {
		u.ensureConfigHasKey(accountConf)

		account := newAccount(u.config, accountConf)
		account.session.SessionEventHandler = u
		accounts[i] = account
	}

	return accounts
}

func newAccount(conf *config.Accounts, currentConf *config.Account) *account {
	id := currentConf.ID()
	c, _ := glib.SignalNew(signalName(id, "connected"))
	d, _ := glib.SignalNew(signalName(id, "disconnected"))

	a := &account{
		session: session.NewSession(conf, currentConf),

		connectedSignal:    c,
		disconnectedSignal: d,
	}
	a.session.Account = a

	return a
}

func signalName(id, signal string) string {
	return "coyim-account-" + signal + "-" + id
}

func (u *gtkUI) showAddAccountWindow() {
	account := newAccount(u.config, config.NewAccount())
	accountDialog(account, func() error {
		u.config.Add(account.session.CurrentAccount)
		return u.SaveConfig()
	})
}

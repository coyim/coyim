package gui

import (
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/go-gtk/glib"
)

type Account struct {
	ID                 string
	ConnectedSignal    *glib.Signal
	DisconnectedSignal *glib.Signal

	*config.Config
	*session.Session
}

func (acc *Account) Connected() bool {
	return acc.ConnStatus == session.CONNECTED
}

var (
	errFingerprintAlreadyAuthorized = errors.New("the fingerprint is already authorized")
)

func (acc *Account) AuthorizeFingerprint(uid, fingerprint string) error {
	fpr, err := hex.DecodeString(fingerprint)

	if err != nil {
		return err
	}

	existing := acc.UserIdForFingerprint(fpr)
	if len(existing) != 0 {
		return errFingerprintAlreadyAuthorized
	}

	acc.KnownFingerprints = append(acc.KnownFingerprints, config.KnownFingerprint{
		Fingerprint: fpr, UserId: uid,
	})

	return nil
}

func BuildAccountsFrom(multiAccConfig *config.MultiAccountConfig) []Account {
	accounts := make([]Account, 0, len(multiAccConfig.Accounts))

	for i := range multiAccConfig.Accounts {
		conf := &multiAccConfig.Accounts[i]
		accounts = append(accounts, newAccount(conf))
	}

	return accounts
}

func newAccount(conf *config.Config) Account {
	//id := conf.Account + "-" + strconv.FormatUint(uint64(time.Now().UnixNano()), 10)
	id := strconv.FormatUint(uint64(time.Now().UnixNano()), 10)

	return Account{
		ID:      id,
		Config:  conf,
		Session: session.NewSession(conf),

		ConnectedSignal:    glib.NewSignal(signalName(id, "connected")),
		DisconnectedSignal: glib.NewSignal(signalName(id, "disconnected")),
	}
}

func signalName(id, signal string) string {
	return "coyim-account-" + signal + "-" + id
}

func (u *gtkUI) showAddAccountWindow() {
	conf := config.NewConfig()
	account := Account{
		Config: conf,
	}

	accountDialog(account, func() error {
		err := u.configFileManager.Add(*conf)
		if err != nil {
			return err
		}

		return u.SaveConfig()
	})
}

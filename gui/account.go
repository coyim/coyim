package gui

import (
	"strconv"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/go-gtk/glib"
)

type Account struct {
	ID           string
	Connected    *glib.Signal
	Disconnected *glib.Signal

	*config.Config
	*session.Session
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

		Connected:    glib.NewSignal(signalName(id, "connected")),
		Disconnected: glib.NewSignal(signalName(id, "disconnected")),
	}
}

func signalName(id, signal string) string {
	return "coyim-account-" + signal + "-" + id
}

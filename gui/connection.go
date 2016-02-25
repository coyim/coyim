package gui

import (
	"log"
	"math/rand"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/xmpp/errors"
)

func (u *gtkUI) connectAccount(account *account) {
	switch p := account.session.GetConfig().Password; p {
	case "":
		u.askForPasswordAndConnect(account)
	default:
		go u.connectWithPassword(account, p)
	}
}

func (u *gtkUI) connectWithPassword(account *account, password string) error {
	if !account.session.IsDisconnected() {
		return nil
	}

	removeNotification := u.showConnectAccountNotification(account)
	defer removeNotification()

	err := account.session.Connect(password)
	switch err {
	case config.ErrTorNotRunning:
		u.notifyTorIsNotRunning(account)
	case errors.ErrTCPBindingFailed:
		u.notifyConnectionFailure(account)
	case errors.ErrAuthenticationFailed:
		u.askForPasswordAndConnect(account)
	case errors.ErrConnectionFailed:
		u.notifyConnectionFailure(account)
	}

	return err
}

func (u *gtkUI) askForPasswordAndConnect(account *account) {
	if !account.IsAskingForPassword() {
		accountName := account.session.GetConfig().Account
		doInUIThread(func() {
			account.AskForPassword()
			u.askForPassword(accountName,
				func() {
					account.AskedForPassword()
				},
				func(password string) error {
					account.AskedForPassword()
					return u.connectWithPassword(account, password)
				})
		})
	}
}

func (u *gtkUI) connectWithRandomDelay(a *account) {
	sleepDelay := time.Duration(rand.Int31n(7643)) * time.Millisecond
	log.Printf("connectWithRandomDelay(%v, %v)\n", a.session.GetConfig().Account, sleepDelay)
	time.Sleep(sleepDelay)
	a.session.SetWantToBeOnline(true)
	a.Connect()
}

func (u *gtkUI) connectAllAutomatics(all bool) {
	log.Printf("connectAllAutomatics(%v)\n", all)
	var acc []*account
	for _, a := range u.accounts {
		if (all || a.session.GetConfig().ConnectAutomatically) && a.session.IsDisconnected() {
			acc = append(acc, a)
		}
	}

	for _, a := range acc {
		go u.connectWithRandomDelay(a)
	}
}

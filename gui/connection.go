package gui

import (
	"log"
	"math/rand"
	"time"

	"github.com/gotk3/gotk3/glib"
	"../config"
	"../xmpp"
)

func (u *gtkUI) connectAccount(account *account) {
	switch p := account.session.CurrentAccount.Password; p {
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

	u.showConnectAccountNotification(account)
	defer u.removeConnectAccountNotification(account)

	err := account.session.Connect(password)
	switch err {
	case config.ErrTorNotRunning:
		glib.IdleAdd(u.alertTorIsNotRunning)
	case xmpp.ErrTCPBindingFailed:
		u.askForServerDetailsAndConnect(account, password)
	case xmpp.ErrAuthenticationFailed:
		//TODO: notify authentication failure?
		u.askForPasswordAndConnect(account)
	case xmpp.ErrConnectionFailed:
		u.notifyConnectionFailure(account)
	}

	return err
}

func (u *gtkUI) askForPasswordAndConnect(account *account) {
	accountName := account.session.CurrentAccount.Account
	glib.IdleAdd(func() {
		u.askForPassword(accountName, func(password string) error {
			return u.connectWithPassword(account, password)
		})
	})
}

func (u *gtkUI) askForServerDetailsAndConnect(account *account, password string) {
	conf := account.session.CurrentAccount
	glib.IdleAdd(func() {
		u.askForServerDetails(conf, func() error {
			return u.connectWithPassword(account, password)
		})
	})
}

func (u *gtkUI) connectWithRandomDelay(a *account) {
	sleepDelay := time.Duration(rand.Int31n(7643)) * time.Millisecond
	log.Printf("connectWithRandomDelay(%v, %vms)\n", a.session.CurrentAccount.Account, sleepDelay)
	time.Sleep(sleepDelay)
	a.connect()
}

func (u *gtkUI) connectAllAutomatics(all bool) {
	log.Printf("connectAllAutomatics(%v)\n", all)
	var acc []*account
	for _, a := range u.accounts {
		if (all || a.session.CurrentAccount.ConnectAutomatically) && a.session.IsDisconnected() {
			acc = append(acc, a)
		}
	}

	for _, a := range acc {
		go u.connectWithRandomDelay(a)
	}
}

package gui

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/xmpp/errors"
)

func (u *gtkUI) connectAccount(account *account) {
	p := account.session.GetConfig().Password
	if p == "" {
		u.askForPasswordAndConnect(account, false)
	} else {
		go u.connectWithPassword(account, p)
	}
}

func (u *gtkUI) connectionFailureMoreInfoConnectionLost() {
	u.notify(i18n.Local("Connection lost"), i18n.Local("We lost connection to the server for unknown reasons.\n\nWe will try to reconnect."))
}

func (u *gtkUI) connectionFailureMoreInfoTCPBindingFailed() {
	u.notify(i18n.Local("Connection failure"), i18n.Local("We couldn't connect to the server because we couldn't determine a server address for the given domain.\n\nWe will try to reconnect."))
}

func (u *gtkUI) connectionFailureMoreInfoConnectionFailedGeneric() {
	u.notify(i18n.Local("Connection failure"), i18n.Local("We couldn't connect to the server for some reason - verify that the server address is correct and that you are actually connected to the internet.\n\nWe will try to reconnect."))
}

func (u *gtkUI) connectionFailureMoreInfoConnectionFailed(ee error) func() {
	return func() {
		u.notify(i18n.Local("Connection failure"),
			fmt.Sprintf(i18n.Local("We couldn't connect to the server - verify that the server address is correct and that you are actually connected to the internet.\n\nThis is the error we got: %s\n\nWe will try to reconnect."), ee.Error()))
	}
}

func (u *gtkUI) connectWithPassword(account *account, password string) error {
	if !account.session.IsDisconnected() {
		return nil
	}

	removeNotification := u.showConnectAccountNotification(account)
	defer removeNotification()

	err := account.session.Connect(password, u.verifierFor(account))
	switch err {
	case config.ErrTorNotRunning:
		u.notifyTorIsNotRunning(account)
	case errors.ErrTCPBindingFailed:
		u.notifyConnectionFailure(account, u.connectionFailureMoreInfoTCPBindingFailed)
	case errors.ErrAuthenticationFailed:
		account.cachedPassword = ""
		u.askForPasswordAndConnect(account, false)
	case errors.ErrGoogleAuthenticationFailed:
		account.cachedPassword = ""
		u.askForPasswordAndConnect(account, true)
	case errors.ErrConnectionFailed:
		u.notifyConnectionFailure(account, u.connectionFailureMoreInfoConnectionFailedGeneric)
	default:
		ff, ok := err.(*errors.ErrFailedToConnect)
		if ok {
			u.notifyConnectionFailure(account, u.connectionFailureMoreInfoConnectionFailed(ff))
		}
	}

	return err
}

func (u *gtkUI) askForPasswordAndConnect(account *account, addGoogleWarning bool) {
	if !account.IsAskingForPassword() {
		accountName := account.session.GetConfig().Account
		if account.cachedPassword != "" {
			u.connectWithPassword(account, account.cachedPassword)
			return
		}

		doInUIThread(func() {
			account.AskForPassword()
			u.askForPassword(accountName, addGoogleWarning,
				func() {
					account.session.SetWantToBeOnline(false)
					account.AskedForPassword()
				},
				func(password string) error {
					account.cachedPassword = password
					account.AskedForPassword()
					return u.connectWithPassword(account, password)
				},
				func(password string) {
					account.session.GetConfig().Password = password
					u.SaveConfig()
				},
			)
		})
	}
}

func (u *gtkUI) connectWithRandomDelay(a *account) {
	sleepDelay := time.Duration(rand.Int31n(7643)) * time.Millisecond
	time.Sleep(sleepDelay)
	a.session.SetWantToBeOnline(true)
	a.Connect()
}

func (u *gtkUI) connectAllAutomatics(all bool) {
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

func (u *gtkUI) disconnectAll() {
	for _, a := range u.accounts {
		ca := a
		go func() {
			ca.disconnect()
		}()
	}
}

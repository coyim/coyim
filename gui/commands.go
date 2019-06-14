package gui

import (
	"github.com/coyim/coyim/otr_client"
)

type executable interface {
	execute(u *gtkUI)
}

type connectAccountCmd struct{ a *account }
type disconnectAccountCmd struct{ a *account }
type connectionInfoCmd struct{ a *account }
type editAccountCmd struct{ a *account }
type changePasswordAccountCmd struct{ a *account }
type removeAccountCmd struct{ a *account }
type toggleAutoConnectCmd struct{ a *account }
type toggleAlwaysEncryptCmd struct{ a *account }

func (u *gtkUI) ExecuteCmd(c interface{}) {
	u.commands <- c
}

func (c connectAccountCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.connectAccount(c.a)
	})
}

func (c disconnectAccountCmd) execute(u *gtkUI) {
	go c.a.session.Close()
}

func (c connectionInfoCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.connectionInfoDialog(c.a)
	})
}

func (c editAccountCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.editAccount(c.a)
	})
}

func (c changePasswordAccountCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.changePasswordDialog(c.a)
	})
}

func (c removeAccountCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.removeAccount(c.a)
	})
}

func (c toggleAutoConnectCmd) execute(u *gtkUI) {
	go u.toggleAutoConnectAccount(c.a)
}

func (c toggleAlwaysEncryptCmd) execute(u *gtkUI) {
	go u.toggleAlwaysEncryptAccount(c.a)
}

func (u *gtkUI) watchCommands() {
	for command := range u.commands {
		switch c := command.(type) {
		case executable:
			c.execute(u)
		case otr_client.AuthorizeFingerprintCmd:
			c.Account.AuthorizeFingerprint(c.Peer.String(), c.Fingerprint)
			u.ExecuteCmd(otr_client.SaveApplicationConfigCmd{})
		case otr_client.SaveInstanceTagCmd:
			account := c.Account
			account.InstanceTag = c.InstanceTag
			u.ExecuteCmd(otr_client.SaveApplicationConfigCmd{})
		case otr_client.SaveApplicationConfigCmd:
			u.SaveConfig()
		}
	}
}

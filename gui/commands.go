package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/twstrike/coyim/client"
)

type connectAccountCmd *account
type disconnectAccountCmd *account
type editAccountCmd *account
type removeAccountCmd *account
type toggleAutoConnectCmd *account
type toggleAlwaysEncryptCmd *account

func (u *gtkUI) ExecuteCmd(c interface{}) {
	u.commands <- c
}

func (u *gtkUI) watchCommands() {
	for command := range u.commands {
		switch c := command.(type) {
		case connectAccountCmd:
			glib.IdleAdd(func() {
				u.connectAccount(c)
			})
		case disconnectAccountCmd:
			glib.IdleAdd(func() {
				u.disconnectAccount(c)
			})
		case editAccountCmd:
			glib.IdleAdd(func() {
				u.editAccount(c)
			})
		case removeAccountCmd:
			glib.IdleAdd(func() {
				u.removeAccount(c)
			})
		case toggleAutoConnectCmd:
			go u.toggleAutoConnectAccount(c)
		case toggleAlwaysEncryptCmd:
			go u.toggleAlwaysEncryptAccount(c)
		case client.AuthorizeFingerprintCmd:
			account := c.Account
			uid := c.Peer
			fpr := c.Fingerprint

			//TODO: it could be a different pointer,
			//find the account by ID()
			account.AuthorizeFingerprint(uid, fpr)
			u.ExecuteCmd(client.SaveApplicationConfigCmd{})
		case client.SaveInstanceTagCmd:
			account := c.Account
			account.InstanceTag = c.InstanceTag
			u.ExecuteCmd(client.SaveApplicationConfigCmd{})
		case client.SaveApplicationConfigCmd:
			u.SaveConfig()
		}
	}
}

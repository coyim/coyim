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

func (u *gtkUI) showAddAccountWindow() {
	c := config.NewAccount()
	accountDialog(c, func() {
		u.config.Add(c)
		u.SaveConfig()
	})
}

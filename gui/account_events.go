package gui

import "github.com/twstrike/coyim/session"

func (u *gtkUI) observeAccountEvents() {
	//TODO: check for channel close
	for ev := range u.events {
		switch ev.EventType {
		case session.Connected:
			account := u.findAccountForSession(ev.Session)
			if account == nil {
				return
			}

			u.window.Emit(account.connectedSignal.String())
		case session.Disconnected:
			account := u.findAccountForSession(ev.Session)
			if account != nil {
				u.window.Emit(account.disconnectedSignal.String())
			}

			for _, acc := range u.accounts {
				if acc.session.ConnStatus == session.CONNECTED {
					return
				}
			}

			u.roster.disconnected()
		}
	}
}

package muc

type contact struct {
	*rosterItem
}

type account struct {
	*rosterItem
	contacts []*contact
	rooms    []string
}

type accountManager struct {
	accounts []*account
}

func (m *accountManager) addAccount(a *account) {
	m.accounts = append(m.accounts, a)
}

func (u *gtkUI) initDemoAccounts() {
	u.accountManager = &accountManager{}

	accounts := fakeAccounts()
	for _, a := range accounts {
		u.accountManager.addAccount(a)
	}
}

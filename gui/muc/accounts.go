package muc

type mucAccountContact struct {
	*mucRosterItem
}

type mucAccount struct {
	*mucRosterItem
	contacts []*mucAccountContact
	rooms    []string
}

func (a *mucAccount) isOffline() bool {
	return false
}

func (a *mucAccount) displayName() string {
	if a.name != "" {
		return a.name
	}

	return a.id
}

type mucAccountManager struct {
	accounts []*mucAccount
}

func (m *mucAccountManager) addAccount(a *mucAccount) {
	m.accounts = append(m.accounts, a)
}

func (m *mucUI) initDemoAccounts() {
	m.accountManager = &mucAccountManager{}

	accounts := fakeAccounts()
	for _, a := range accounts {
		m.accountManager.addAccount(a)
	}
}

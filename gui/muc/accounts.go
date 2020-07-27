package muc

type contact struct {
	*rosterItem
}

type accountRoom struct {
	id     string
	status peerStatus
}

type group struct {
	*rosterItem
	contacts []*contact
	rooms    []*accountRoom
}

type account struct {
	*rosterItem
	groups   []*group
	contacts []*contact
	rooms    []*accountRoom
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

func (g *group) hasStatus() bool {
	return false
}

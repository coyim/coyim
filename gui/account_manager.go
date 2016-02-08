package gui

import (
	"sync"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/xmpp/interfaces"
)

type accountManager struct {
	accounts []*account
	events   chan interface{}
	contacts map[*account]*rosters.List

	client.CommandManager

	sync.RWMutex
}

func newAccountManager(c client.CommandManager) *accountManager {
	return &accountManager{
		events:         make(chan interface{}, 10),
		accounts:       make([]*account, 0, 5),
		contacts:       make(map[*account]*rosters.List),
		CommandManager: c,
	}
}

func (m *accountManager) disconnectAll() {
	for _, acc := range m.accounts {
		acc.disconnect()
	}
}

func (m *accountManager) getAccountByID(ID string) (*account, bool) {
	m.RLock()
	defer m.RUnlock()

	for _, acc := range m.accounts {
		if acc.ID() == ID {
			return acc, true
		}
	}

	return nil, false
}

func (m *accountManager) findAccountForSession(s access.Session) *account {
	acc, _ := m.getAccountByID(s.GetConfig().ID())
	return acc
}

func (m *accountManager) getAllContacts() map[*account]*rosters.List {
	return m.contacts
}

func (m *accountManager) getContacts(acc *account) *rosters.List {
	return m.contacts[acc]
}

func (m *accountManager) setContacts(account *account, contacts *rosters.List) {
	if account == nil {
		panic("Developer error: account should never be nil")
	}
	m.contacts[account] = contacts
}

func (m *accountManager) addAccount(appConfig *config.ApplicationConfig, account *config.Account, sf access.Factory, df func() interfaces.Dialer) {
	m.Lock()
	defer m.Unlock()

	acc := newAccount(appConfig, account, sf, df)
	acc.session.Subscribe(m.events)
	acc.session.SetCommandManager(m)
	acc.session.SetConnector(acc)

	m.accounts = append(m.accounts, acc)
	m.setContacts(acc, rosters.New())
}

func (m *accountManager) removeAccount(conf *config.Account, k func()) {
	toRemove, exists := m.getAccountByID(conf.ID())
	if !exists {
		return
	}

	m.Lock()
	defer m.Unlock()

	accs := make([]*account, 0, len(m.accounts)-1)
	for _, acc := range m.accounts {
		if acc == toRemove {
			delete(m.contacts, acc)
			continue
		}

		accs = append(accs, acc)
	}

	m.accounts = accs

	k()
}

func (m *accountManager) buildAccounts(appConfig *config.ApplicationConfig, sf access.Factory, df func() interfaces.Dialer) {
	hasConfUpdates := false
	for _, accountConf := range appConfig.Accounts {
		if _, ok := m.getAccountByID(accountConf.ID()); ok {
			continue
		}

		hasUpdate, err := accountConf.EnsurePrivateKey()
		if err != nil {
			continue
		}

		hasConfUpdates = hasConfUpdates || hasUpdate
		m.addAccount(appConfig, accountConf, sf, df)
	}

	if hasConfUpdates {
		m.ExecuteCmd(client.SaveApplicationConfigCmd{})
	}
}

func (m *accountManager) addNewAccountsFromConfig(appConfig *config.ApplicationConfig, sf access.Factory, df func() interfaces.Dialer) {
	for _, configAccount := range appConfig.Accounts {
		_, found := m.getAccountByID(configAccount.ID())
		if found {
			continue
		}

		m.addAccount(appConfig, configAccount, sf, df)
	}
}

func (m *accountManager) removePeer(account *account, peer string) {
	m.Lock()
	defer m.Unlock()

	l, ok := m.contacts[account]
	if !ok {
		return
	}

	l.Remove(peer)
}

func (m *accountManager) getPeer(account *account, peer string) (*rosters.Peer, bool) {
	m.RLock()
	defer m.RUnlock()

	l, ok := m.contacts[account]
	if !ok {
		return nil, false
	}

	return l.Get(peer)
}

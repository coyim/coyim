package gui

import (
	"sync"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/session"
)

type accountManager struct {
	accounts []*account
	events   chan interface{}

	client.CommandManager

	sync.RWMutex
}

func newAccountManager(c client.CommandManager) *accountManager {
	return &accountManager{
		events:         make(chan interface{}, 10),
		accounts:       make([]*account, 0, 5),
		CommandManager: c,
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

func (m *accountManager) findAccountForSession(s *session.Session) *account {
	acc, _ := m.getAccountByID(s.GetConfig().ID())
	return acc
}

func (m *accountManager) addAccount(appConfig *config.ApplicationConfig, account *config.Account) {
	m.Lock()
	defer m.Unlock()

	acc, err := newAccount(appConfig, account)
	if err != nil {
		//TODO error
		return
	}

	acc.session.Subscribe(m.events)
	acc.session.CommandManager = m

	m.accounts = append(m.accounts, acc)
}

func (m *accountManager) removeAccount(conf *config.Account) {
	toRemove, exists := m.getAccountByID(conf.ID())
	if !exists {
		return
	}

	m.Lock()
	defer m.Unlock()

	accs := make([]*account, 0, len(m.accounts)-1)
	for _, acc := range m.accounts {
		if acc == toRemove {
			continue
		}

		accs = append(accs, acc)
	}

	m.accounts = accs
}

func (m *accountManager) buildAccounts(appConfig *config.ApplicationConfig) {
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
		m.addAccount(appConfig, accountConf)
	}

	if hasConfUpdates {
		m.ExecuteCmd(client.SaveApplicationConfigCmd{})
	}
}

func (m *accountManager) addNewAccountsFromConfig(appConfig *config.ApplicationConfig) {
	for _, configAccount := range appConfig.Accounts {
		_, found := m.getAccountByID(configAccount.ID())
		if found {
			continue
		}

		m.addAccount(appConfig, configAccount)
	}
}

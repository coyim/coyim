package gui

import (
	"sync"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/session"
)

var accountsLock sync.Mutex

//TODO: once glib signals are removed from account, it could be a sessionManager
// and this could be in a client package (ui agnostic)
type accountManager struct {
	accounts          []*account
	events            chan interface{}
	saveConfiguration func()

	client.CommandManager
}

func newAccountManager(c client.CommandManager) *accountManager {
	return &accountManager{
		events:         make(chan interface{}, 10),
		CommandManager: c,
	}
}

func (m *accountManager) addAccount(appConfig *config.ApplicationConfig, account *config.Account) {
	acc, err := newAccount(appConfig, account)
	if err != nil {
		//TODO error
		return
	}

	//We dont need this anymore, only CLI
	//acc.session.SaveConfiguration = m.saveConfiguration
	acc.session.Subscribe(m.events)

	//TODO: remove this duplication
	acc.CommandManager = m
	acc.session.CommandManager = m

	m.accounts = append(m.accounts, acc)
}

func (m *accountManager) buildAccounts(appConfig *config.ApplicationConfig) {
	m.accounts = make([]*account, 0, len(appConfig.Accounts))
	hasConfUpdates := false
	for _, accountConf := range appConfig.Accounts {
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

func (m *accountManager) findAccountForSession(s *session.Session) *account {
	for _, a := range m.accounts {
		if a.session == s {
			return a
		}
	}

	return nil
}

func (m *accountManager) findAccountForUsername(s string) *account {
	for _, a := range m.accounts {
		if a.session.CurrentAccount.Is(s) {
			return a
		}
	}

	return nil
}

func (m *accountManager) addNewAccountsFromConfig(appConfig *config.ApplicationConfig) {
	for _, configAccount := range appConfig.Accounts {
		var found bool
		for _, acc := range m.accounts {
			if acc.session.CurrentAccount.ID() == configAccount.ID() {
				found = true
				break
			}
		}

		if found {
			continue
		}

		m.addAccount(appConfig, configAccount)
	}
}

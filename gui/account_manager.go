package gui

import (
	"fmt"
	"sync"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/otrclient"
	rosters "github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

type accountManager struct {
	accounts []*account
	events   chan interface{}
	contacts map[*account]*rosters.List

	otrclient.CommandManager

	log coylog.Logger

	sync.RWMutex
}

func newAccountManager(c otrclient.CommandManager, log coylog.Logger) *accountManager {
	return &accountManager{
		events:         make(chan interface{}, 10),
		accounts:       make([]*account, 0, 5),
		contacts:       make(map[*account]*rosters.List),
		CommandManager: c,
		log:            log,
	}
}

func (m *accountManager) getAllAccounts() []*account {
	m.RLock()
	defer m.RUnlock()

	return append([]*account(nil), m.accounts...)
}

func (m *accountManager) getAllConnectedAccounts() []*account {
	m.RLock()
	defer m.RUnlock()

	accounts := make([]*account, 0, len(m.accounts))
	for _, acc := range m.accounts {
		if !acc.connected() {
			continue
		}

		accounts = append(accounts, acc)
	}

	return accounts
}

func (m *accountManager) disconnectAll() {
	m.RLock()
	defer m.RUnlock()

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

func (u *gtkUI) addAccount(appConfig *config.ApplicationConfig, account *config.Account, sf access.Factory, df interfaces.DialerFactory) {
	u.Lock()
	defer u.Unlock()

	acc := newAccount(appConfig, account, sf, df)
	go u.observeAccountEvents(acc)
	acc.log = u.log.WithField("account", account.Account)
	acc.session.Subscribe(acc.events)
	acc.session.SetCommandManager(u)
	acc.session.SetConnector(acc)

	u.accounts = append(u.accounts, acc)
	u.setContacts(acc, rosters.New())
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

func (u *gtkUI) buildAccounts(appConfig *config.ApplicationConfig, sf access.Factory, df interfaces.DialerFactory) {
	hasConfUpdates := false
	for _, accountConf := range appConfig.Accounts {
		if _, ok := u.getAccountByID(accountConf.ID()); ok {
			continue
		}

		hasUpdate, err := accountConf.EnsurePrivateKey()
		if err != nil {
			continue
		}

		hasConfUpdates = hasConfUpdates || hasUpdate
		u.addAccount(appConfig, accountConf, sf, df)
	}

	if hasConfUpdates {
		u.ExecuteCmd(otrclient.SaveApplicationConfigCmd{})
	}
}

func (u *gtkUI) addNewAccountsFromConfig(appConfig *config.ApplicationConfig, sf access.Factory, df interfaces.DialerFactory) {
	for _, configAccount := range appConfig.Accounts {
		_, found := u.getAccountByID(configAccount.ID())
		if found {
			continue
		}

		u.addAccount(appConfig, configAccount, sf, df)
	}
}

func (m *accountManager) removePeer(account *account, peer jid.WithoutResource) {
	m.Lock()
	defer m.Unlock()

	l, ok := m.contacts[account]
	if !ok {
		return
	}

	l.Remove(peer)
}

func (m *accountManager) getPeer(account *account, peer jid.WithoutResource) (*rosters.Peer, bool) {
	m.RLock()
	defer m.RUnlock()

	l, ok := m.contacts[account]
	if !ok {
		account.log.Warn("Failure to look up account from account manager")
		return nil, false
	}

	return l.Get(peer)
}

func (m *accountManager) displayNameFor(account *account, peer jid.WithoutResource) string {
	p, ok := m.getPeer(account, peer)
	if !ok {
		return peer.String()
	}

	return p.NameForPresentation()
}

func (m *accountManager) debugPeersFor(account *account) {
	m.RLock()
	defer m.RUnlock()

	rs, ok := m.getAllContacts()[account]
	if ok {
		rs.Iter(func(_ int, item *rosters.Peer) {
			fmt.Printf("->   %s\n", item.Dump())
		})
	}

	fmt.Printf(" ************************************** \n")
	fmt.Println()
}

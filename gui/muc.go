package gui

import "github.com/coyim/gotk3adapter/gtki"

type mucUI struct {
	u              *gtkUI
	accountManager *mucAccountManager
	roster         *mucRoster
	builder        *builder
}

type mucAccountStatus string

var (
	mucStatusConnecting mucAccountStatus = "connecting"
	mucStatusOnline     mucAccountStatus = "online"
	mucStatusOffline    mucAccountStatus = "offline"
)

type notificationStatus string

type mucAccount struct {
	id       string
	name     string
	status   mucAccountStatus
	accounts []*mucAccountStatus
}

type mucAccountManager struct {
	accounts []*mucAccount
}

type mucRoster struct {
	widget gtki.ScrolledWindow `gtk-widget:"roster"`
	model  gtki.ListStore      `gtk-widget:"roster-model"`
	view   gtki.TreeView       `gtk-widget:"roster-tree"`

	u *gtkUI
}

func (u *gtkUI) initMUCMockups() {
	builder := newBuilder("MUC")

	m := &mucUI{
		u:       u,
		builder: builder,
	}

	m.init()

	m.addAccountsToRoster()

	m.showWindow()
}

func (m *mucUI) init() {
	m.initDemoAccounts()
	m.initRoster()
}

func (m *mucUI) initDemoAccounts() {
	m.accountManager = &mucAccountManager{}

	accounts := []*mucAccount{
		&mucAccount{
			id:     "sandy@autonomia.digital",
			status: mucStatusOnline,
		},
		&mucAccount{
			id:     "pedro@autonomia.digital",
			status: mucStatusOnline,
		},
		&mucAccount{
			id:     "pedro@coy.im",
			name:   "Pedro CoyIM",
			status: mucStatusOffline,
		},
	}

	for _, a := range accounts {
		m.accountManager.addAccount(a)
	}
}

func (m *mucUI) initRoster() {
	r := &mucRoster{
		u: m.u,
	}

	panicOnDevError(m.builder.bindObjects(r))

	m.roster = r
}

func (m *mucUI) addAccountsToRoster() {
	for _, a := range m.accountManager.accounts {
		m.roster.add(a)
	}
}

func (m *mucUI) showWindow() {
	win := m.builder.get("mucWindow").(gtki.Window)
	win.Show()
}

func (m *mucAccountManager) addAccount(a *mucAccount) {
	m.accounts = append(m.accounts, a)
}

func (r *mucRoster) add(account *mucAccount) {
	r.addAccount(account)
	r.view.ExpandAll()
}

func (r *mucRoster) addAccount(account *mucAccount) {
	cs := r.u.currentColorSet()
	parentIter := r.model.Append()

	accountCounter := &counter{}

	displayName := account.displayName()

	_ = r.model.SetValue(parentIter, indexParentJid, displayName)
	_ = r.model.SetValue(parentIter, indexAccountID, account.id)
	_ = r.model.SetValue(parentIter, indexRowType, "account")
	_ = r.model.SetValue(parentIter, indexWeight, 700)

	bgcolor := cs.rosterAccountOnlineBackground
	if account.isOffline() {
		bgcolor = cs.rosterAccountOfflineBackground
	}
	_ = r.model.SetValue(parentIter, indexBackgroundColor, bgcolor)

	_ = r.model.SetValue(parentIter, indexStatusIcon, statusIcons[account.getStatus()].GetPixbuf())
	_ = r.model.SetValue(parentIter, indexParentDisplayName, createGroupDisplayName(displayName, accountCounter, true))
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

func (a *mucAccount) getStatus() string {
	if a.status == mucStatusConnecting {
		return "connecting"
	}

	if a.status == mucStatusOffline {
		return "offline"
	}

	return "available"
}

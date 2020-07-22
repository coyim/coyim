package gui

import (
	"fmt"
	"html"

	"github.com/coyim/gotk3adapter/gtki"
)

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

type mucAccountContact struct {
	id     string
	name   string
	status mucAccountStatus
}

type mucAccount struct {
	id       string
	name     string
	status   mucAccountStatus
	contacts []*mucAccountContact
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
			contacts: []*mucAccountContact{
				&mucAccountContact{
					id:     "pedro@autonomia.digital",
					name:   "Pedro Enrique",
					status: mucStatusOnline,
				},
				&mucAccountContact{
					id:     "rafael@autonomia.digital",
					status: mucStatusOnline,
				},
				&mucAccountContact{
					id:     "cristina@autonomia.digital",
					name:   "Cristina Salcedo",
					status: mucStatusOffline,
				},
			},
		},
		&mucAccount{
			id:     "pedro@autonomia.digital",
			status: mucStatusOnline,
			contacts: []*mucAccountContact{
				&mucAccountContact{
					id:     "sandy@autonomia.digital",
					name:   "Sandy Acurio",
					status: mucStatusOnline,
				},
				&mucAccountContact{
					id:     "rafael@autonomia.digital",
					status: mucStatusOnline,
				},
				&mucAccountContact{
					id:     "cristina@autonomia.digital",
					name:   "Cristina Salcedo",
					status: mucStatusOffline,
				},
			},
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

	// Contacts for this account
	r.addAccountContacts(account.contacts, parentIter, accountCounter)

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

func (r *mucRoster) addAccountContacts(contacts []*mucAccountContact, parentIter gtki.TreeIter, accountCounter *counter) {
	groupCounter := &counter{}

	for _, item := range contacts {
		o := item.isOnline()
		accountCounter.inc(true, o)
		groupCounter.inc(true, o)
		r.addItem(item)
	}
}

func (r *mucRoster) addItem(item *mucAccountContact) {
	cs := r.u.currentColorSet()
	iter := r.model.Append()

	setValues(
		r.model,
		iter,
		item.id,
		item.displayName(),
		"BelongsTo",
		decideColorForContact(cs, item),
		cs.rosterPeerBackground,
		nil,
		createTooltipForContact(item),
	)
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

func (c *mucAccountContact) displayName() string {
	if c.name != "" {
		return c.name
	}

	return c.id
}

func (c *mucAccountContact) isOnline() bool {
	return c.status == mucStatusOnline
}

func setValues(v gtki.ListStore, iter gtki.TreeIter, values ...interface{}) {
	for i, val := range values {
		if val != nil {
			_ = v.SetValue(iter, i, val)
		}
	}
}

func decideColorForContact(cs colorSet, c *mucAccountContact) string {
	if !c.isOnline() {
		return cs.rosterPeerOfflineForeground
	}
	return cs.rosterPeerOnlineForeground
}

func createTooltipForContact(c *mucAccountContact) string {
	pname := html.EscapeString(c.displayName())
	jid := html.EscapeString(c.id)
	if pname != jid {
		return fmt.Sprintf("%s (%s)", pname, jid)
	}
	return jid
}

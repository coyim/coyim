package gui

import (
	"sort"
	"time"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtk_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session/events"
	smock "github.com/twstrike/coyim/session/mock"
)

type AccountSuite struct{}

var _ = Suite(&AccountSuite{})

type namedSessionMock struct {
	smock.SessionMock
	name string
}

func (v *namedSessionMock) GetConfig() *config.Account {
	return &config.Account{Account: v.name}
}

func (*AccountSuite) Test_account_sorting(c *C) {
	one := &account{session: &namedSessionMock{name: "bca"}}
	two := &account{session: &namedSessionMock{name: "abc"}}
	three := &account{session: &namedSessionMock{name: "cba"}}

	accounts := []*account{one, two, three}

	sort.Sort(byAccountNameAlphabetic(accounts))

	c.Assert(accounts, DeepEquals, []*account{two, one, three})
}

type accountInfoBarMock struct {
	gtk_mock.MockInfoBar

	hideCalled, destroyCalled int
}

func (v *accountInfoBarMock) Hide() {
	v.hideCalled++
}

func (v *accountInfoBarMock) Destroy() {
	v.destroyCalled++
}

func (*AccountSuite) Test_account_removeCurrentNotification_doesNothingIfItIsNil(c *C) {
	ac := &account{currentNotification: nil}
	ac.removeCurrentNotification()

	c.Assert(ac.currentNotification, IsNil)
}

func (*AccountSuite) Test_account_removeCurrentNotification_removesNotificationIfExists(c *C) {
	one := &accountInfoBarMock{}
	ac := &account{currentNotification: one}
	ac.removeCurrentNotification()

	c.Assert(ac.currentNotification, IsNil)
	c.Assert(one.hideCalled, Equals, 1)
	c.Assert(one.destroyCalled, Equals, 1)
}

func (*AccountSuite) Test_account_removeCurrentNotificationIf_doesNothingIfItIsntTheSameNotification(c *C) {
	one := &accountInfoBarMock{}
	two := &accountInfoBarMock{}
	ac := &account{currentNotification: one}
	ac.removeCurrentNotificationIf(two)

	c.Assert(ac.currentNotification, Equals, one)
	c.Assert(one.hideCalled, Equals, 0)
	c.Assert(one.destroyCalled, Equals, 0)
}

func (*AccountSuite) Test_account_removeCurrentNotificationIf_removesTheNotificationIfItMatches(c *C) {
	one := &accountInfoBarMock{}
	ac := &account{currentNotification: one}
	ac.removeCurrentNotificationIf(one)

	c.Assert(ac.currentNotification, IsNil)
	c.Assert(one.hideCalled, Equals, 1)
	c.Assert(one.destroyCalled, Equals, 1)
}

func (*AccountSuite) Test_account_IsAskingForPassword(c *C) {
	c.Assert((&account{askingForPassword: true}).IsAskingForPassword(), Equals, true)
	c.Assert((&account{askingForPassword: false}).IsAskingForPassword(), Equals, false)
}

func (*AccountSuite) Test_account_AskForPassword(c *C) {
	a := &account{}
	a.AskForPassword()
	c.Assert(a.askingForPassword, Equals, true)
}

func (*AccountSuite) Test_account_AskedForPassword(c *C) {
	a := &account{askingForPassword: true}
	a.AskedForPassword()
	c.Assert(a.askingForPassword, Equals, false)
}

type accountMockGtk struct {
	gtk_mock.Mock
}

func (*accountMockGtk) MenuNew() (gtki.Menu, error) {
	return &accountMockMenu{}, nil
}

func (*accountMockGtk) MenuItemNewWithMnemonic(mnem string) (gtki.MenuItem, error) {
	return &accountMockMenuItem{mnemonic: mnem, sensitive: true}, nil
}

func (*accountMockGtk) CheckMenuItemNewWithMnemonic(mnem string) (gtki.CheckMenuItem, error) {
	return &accountMockCheckMenuItem{mnemonic: mnem}, nil
}

func (*accountMockGtk) SeparatorMenuItemNew() (gtki.SeparatorMenuItem, error) {
	return &accountMockSeparatorMenuItem{}, nil
}

type accountMockMenu struct {
	gtk_mock.MockMenu

	menuItems []gtki.MenuItem
}

func (v *accountMockMenu) Append(v1 gtki.MenuItem) {
	v.menuItems = append(v.menuItems, v1)
}

type accountMockMenuItem struct {
	gtk_mock.MockMenuItem

	mnemonic  string
	sensitive bool

	onActivate interface{}
}

type accountMockCheckMenuItem struct {
	gtk_mock.MockCheckMenuItem

	mnemonic string
	active   bool

	onActivate interface{}
}

func (v *accountMockCheckMenuItem) SetActive(v1 bool) {
	v.active = v1
}

func (v *accountMockMenuItem) Connect(p string, v1 interface{}, v2 ...interface{}) (glibi.SignalHandle, error) {
	if p == "activate" {
		v.onActivate = v1
	}

	return glibi.SignalHandle(0), nil
}

func (v *accountMockMenuItem) SetSensitive(v1 bool) {
	v.sensitive = v1
}

func (v *accountMockCheckMenuItem) Connect(p string, v1 interface{}, v2 ...interface{}) (glibi.SignalHandle, error) {
	if p == "activate" {
		v.onActivate = v1
	}

	return glibi.SignalHandle(0), nil
}

type accountMockSeparatorMenuItem struct {
	gtk_mock.MockSeparatorMenuItem
}

type accountMockGlib struct {
	glib_mock.Mock
}

func (*accountMockGlib) Local(vx string) string {
	return "[localized] " + vx
}

func (*AccountSuite) Test_account_createSubmenu_createsTheGeneralStructure(c *C) {
	i18n.InitLocalization(&accountMockGlib{})
	g = Graphics{gtk: &accountMockGtk{}}

	sess := &accountMockSession{config: &config.Account{}}
	a := &account{session: sess}
	menu := a.createSubmenu()
	c.Assert(menu, Not(IsNil))

	createdMenu := menu.(*accountMockMenu)

	c.Assert(createdMenu.menuItems, Not(IsNil))
	c.Assert(createdMenu.menuItems[0].(*accountMockMenuItem).mnemonic, Equals, "[localized] _Connect")
	c.Assert(createdMenu.menuItems[1].(*accountMockMenuItem).mnemonic, Equals, "[localized] _Disconnect")

	_, ok := createdMenu.menuItems[2].(*accountMockSeparatorMenuItem)
	c.Assert(ok, Equals, true)

	c.Assert(createdMenu.menuItems[4].(*accountMockMenuItem).mnemonic, Equals, "[localized] _Edit...")
	c.Assert(createdMenu.menuItems[5].(*accountMockMenuItem).mnemonic, Equals, "[localized] _Remove")

	_, ok = createdMenu.menuItems[6].(*accountMockSeparatorMenuItem)
	c.Assert(ok, Equals, true)

	c.Assert(createdMenu.menuItems[7].(*accountMockCheckMenuItem).mnemonic, Equals, "[localized] Connect _Automatically")
	c.Assert(createdMenu.menuItems[8].(*accountMockCheckMenuItem).mnemonic, Equals, "[localized] Always Encrypt Conversation")
}

func (*AccountSuite) Test_account_createSubmenu_setsTheCheckboxesCorrectly(c *C) {
	i18n.InitLocalization(&accountMockGlib{})
	g = Graphics{gtk: &accountMockGtk{}}

	conf := &config.Account{ConnectAutomatically: true, AlwaysEncrypt: true}
	sess := &accountMockSession{config: conf}
	a := &account{session: sess}

	menu := a.createSubmenu()
	createdMenu := menu.(*accountMockMenu)
	c.Assert(createdMenu.menuItems[7].(*accountMockCheckMenuItem).active, Equals, true)
	c.Assert(createdMenu.menuItems[8].(*accountMockCheckMenuItem).active, Equals, true)

	conf.AlwaysEncrypt = false
	menu = a.createSubmenu()
	createdMenu = menu.(*accountMockMenu)
	c.Assert(createdMenu.menuItems[7].(*accountMockCheckMenuItem).active, Equals, true)
	c.Assert(createdMenu.menuItems[8].(*accountMockCheckMenuItem).active, Equals, false)

	conf.ConnectAutomatically = false
	menu = a.createSubmenu()
	createdMenu = menu.(*accountMockMenu)
	c.Assert(createdMenu.menuItems[7].(*accountMockCheckMenuItem).active, Equals, false)
	c.Assert(createdMenu.menuItems[8].(*accountMockCheckMenuItem).active, Equals, false)

	conf.AlwaysEncrypt = true
	menu = a.createSubmenu()
	createdMenu = menu.(*accountMockMenu)
	c.Assert(createdMenu.menuItems[7].(*accountMockCheckMenuItem).active, Equals, false)
	c.Assert(createdMenu.menuItems[8].(*accountMockCheckMenuItem).active, Equals, true)
}

func (*AccountSuite) Test_account_createSubmenu_setsActivationCorrectly(c *C) {
	i18n.InitLocalization(&accountMockGlib{})
	g = Graphics{gtk: &accountMockGtk{}}

	sess := &accountMockSession{config: &config.Account{}}
	a := &account{session: sess}

	menu := a.createSubmenu()
	createdMenu := menu.(*accountMockMenu)

	// We can't really check that these things are set to the correct functions, just that they are set
	// It might be possible to try invoking them and see that they do the right things, at some point
	// For now, too much bother.

	c.Assert(createdMenu.menuItems[0].(*accountMockMenuItem).onActivate, Not(IsNil))
	c.Assert(createdMenu.menuItems[1].(*accountMockMenuItem).onActivate, Not(IsNil))

	c.Assert(createdMenu.menuItems[4].(*accountMockMenuItem).onActivate, Not(IsNil))
	c.Assert(createdMenu.menuItems[5].(*accountMockMenuItem).onActivate, Not(IsNil))

	c.Assert(createdMenu.menuItems[7].(*accountMockCheckMenuItem).onActivate, Not(IsNil))
	c.Assert(createdMenu.menuItems[8].(*accountMockCheckMenuItem).onActivate, Not(IsNil))
}

type accountMockSession struct {
	smock.SessionMock

	isDisconnected bool
	config         *config.Account
	events         []chan<- interface{}
}

func (v *accountMockSession) IsDisconnected() bool {
	return v.isDisconnected
}

func (v *accountMockSession) GetConfig() *config.Account {
	return v.config
}

func (v *accountMockSession) Subscribe(v1 chan<- interface{}) {
	v.events = append(v.events, v1)
}

func (*AccountSuite) Test_account_createSubmenu_setsConnectAndDisconnectSensitivity(c *C) {
	i18n.InitLocalization(&accountMockGlib{})
	g = Graphics{gtk: &accountMockGtk{}}

	sess := &accountMockSession{isDisconnected: true, config: &config.Account{}}
	a := &account{session: sess}

	menu := a.createSubmenu()
	createdMenu := menu.(*accountMockMenu)
	c.Assert(createdMenu.menuItems[0].(*accountMockMenuItem).sensitive, Equals, true)
	c.Assert(createdMenu.menuItems[1].(*accountMockMenuItem).sensitive, Equals, false)

	sess.isDisconnected = false
	menu = a.createSubmenu()
	createdMenu = menu.(*accountMockMenu)
	c.Assert(createdMenu.menuItems[0].(*accountMockMenuItem).sensitive, Equals, false)
	c.Assert(createdMenu.menuItems[1].(*accountMockMenuItem).sensitive, Equals, true)
}

func (*AccountSuite) Test_account_createSubmenu_willWatchForThingsToChangeTheConnectSensitivity(c *C) {
	i18n.InitLocalization(&accountMockGlib{})
	g = Graphics{gtk: &accountMockGtk{}}

	sess := &accountMockSession{isDisconnected: true, config: &config.Account{}}
	a := &account{session: sess}

	menu := a.createSubmenu()
	connectItem := menu.(*accountMockMenu).menuItems[0].(*accountMockMenuItem)

	c.Assert(connectItem.sensitive, Equals, true)

	sess.isDisconnected = false
	for _, cc := range sess.events {
		cc <- events.Event{
			Type: events.Connecting,
		}
	}

	waitFor(c, func() bool { return !connectItem.sensitive })

	sess.isDisconnected = false
	for _, cc := range sess.events {
		cc <- events.Event{
			Type: events.Connected,
		}
	}

	waitFor(c, func() bool { return !connectItem.sensitive })

	sess.isDisconnected = true
	for _, cc := range sess.events {
		cc <- events.Event{
			Type: events.Disconnected,
		}
	}

	waitFor(c, func() bool { return connectItem.sensitive })
}

func waitFor(c *C, f func() bool) {
	cx := make(chan bool)

	go func() {
		for !f() {
			time.Sleep(time.Duration(20) * time.Millisecond)
		}
		cx <- true
	}()

	select {
	case <-time.After(5 * time.Second):
		c.Assert(f(), Equals, true)
	case <-cx:
		c.Assert(f(), Equals, true)
	}
}

func (*AccountSuite) Test_account_createSubmenu_willWatchForThingsToChangeTheDisconnectSensitivity(c *C) {
	i18n.InitLocalization(&accountMockGlib{})
	g = Graphics{gtk: &accountMockGtk{}}

	sess := &accountMockSession{isDisconnected: true, config: &config.Account{}}
	a := &account{session: sess}

	menu := a.createSubmenu()
	disconnectItem := menu.(*accountMockMenu).menuItems[1].(*accountMockMenuItem)

	c.Assert(disconnectItem.sensitive, Equals, false)

	sess.isDisconnected = false
	for _, cc := range sess.events {
		cc <- events.Event{
			Type: events.Connecting,
		}
	}

	waitFor(c, func() bool { return disconnectItem.sensitive })

	sess.isDisconnected = false
	for _, cc := range sess.events {
		cc <- events.Event{
			Type: events.Connected,
		}
	}

	waitFor(c, func() bool { return disconnectItem.sensitive })

	sess.isDisconnected = true
	for _, cc := range sess.events {
		cc <- events.Event{
			Type: events.Disconnected,
		}
	}

	waitFor(c, func() bool { return !disconnectItem.sensitive })
}

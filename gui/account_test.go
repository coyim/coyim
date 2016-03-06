package gui

import (
	"sort"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtk_mock"
	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/session/mock"
)

type AccountSuite struct{}

var _ = Suite(&AccountSuite{})

type namedSessionMock struct {
	mocks.SessionMock
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

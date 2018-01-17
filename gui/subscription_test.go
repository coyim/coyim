package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtk_mock"
	"github.com/coyim/gotk3adapter/gtki"
	. "gopkg.in/check.v1"
)

type SubscriptionSuite struct{}

var _ = Suite(&SubscriptionSuite{})

type subscriptionGtkMock struct {
	gtk_mock.Mock
	builder *subscriptionBuilderMock
}

func (v *subscriptionGtkMock) BuilderNew() (gtki.Builder, error) {
	v.builder = &subscriptionBuilderMock{}
	return v.builder, nil
}

type subscriptionBuilderMock struct {
	gtk_mock.MockBuilder
	dialog *subscriptionMessageDialogMock
}

func (v *subscriptionBuilderMock) GetObject(v1 string) (glibi.Object, error) {
	if v1 == "dialog" {
		v.dialog = &subscriptionMessageDialogMock{}
		return v.dialog, nil
	}
	return nil, nil
}

type subscriptionMessageDialogMock struct {
	gtk_mock.MockMessageDialog

	propertyType, propertyValue string
	transientFor                gtki.Window
}

func (v *subscriptionMessageDialogMock) SetProperty(v1 string, v2 interface{}) error {
	v.propertyType = v1
	v.propertyValue = v2.(string)
	return nil
}

func (v *subscriptionMessageDialogMock) SetTransientFor(v2 gtki.Window) {
	v.transientFor = v2
}

func (*SubscriptionSuite) Test_authorizePresenceSubscriptionDialog_setsTextPropertyCorrectly(c *C) {
	sm := &subscriptionGtkMock{}
	g = Graphics{gtk: sm}
	authorizePresenceSubscriptionDialog(nil, jid.NR("hello@world.org"))
	c.Assert(sm.builder.dialog.propertyType, Equals, "text")
	c.Assert(sm.builder.dialog.propertyValue, Equals, "hello@world.org wants to talk to you. Is that ok?")
}

func (*SubscriptionSuite) Test_authorizePresenceSubscriptionDialog_setsTransientForCorrectly(c *C) {
	sm := &subscriptionGtkMock{}
	g = Graphics{gtk: sm}
	w := &gtk_mock.MockWindow{}
	authorizePresenceSubscriptionDialog(w, jid.NR("hello@world.org"))
	c.Assert(sm.builder.dialog.transientFor, Equals, w)
}

func (*SubscriptionSuite) Test_authorizePresenceSubscriptionDialog_returnsTheDialog(c *C) {
	sm := &subscriptionGtkMock{}
	g = Graphics{gtk: sm}
	w := &gtk_mock.MockWindow{}
	ret := authorizePresenceSubscriptionDialog(w, jid.NR("hello@world.org"))
	c.Assert(sm.builder.dialog, Equals, ret)
}

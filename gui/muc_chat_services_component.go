package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type chatServicesComponent struct {
	u *gtkUI

	chatServicesBox  gtki.Box          `gtk-widget:"chat-services-content"`
	chatServices     gtki.ComboBoxText `gtk-widget:"chat-services-combobox-text"`
	chatServiceEntry gtki.Entry        `gtk-widget:"chat-service-entry"`

	previousUpdateChannel chan bool

	onServiceChanged func()
}

func (u *gtkUI) createChatServicesComponent(onServiceChanged func()) *chatServicesComponent {
	c := &chatServicesComponent{
		u:                u,
		onServiceChanged: onServiceChanged,
	}

	c.initBuilder()

	return c
}

func (c *chatServicesComponent) initBuilder() {
	b := newBuilder("MUCChatServices")
	panicOnDevError(b.bindObjects(c))

	b.ConnectSignals(map[string]interface{}{
		"on_service_changed": c.onServiceChanged,
	})
}

func (c *chatServicesComponent) updateServicesBasedOnAccount(ca *account) {
	if c.previousUpdateChannel != nil {
		c.previousUpdateChannel <- true
	}

	c.previousUpdateChannel = make(chan bool)

	csc, ec, endEarly := ca.session.GetChatServices(jid.ParseDomain(ca.Account()))

	go c.updateChatServices(ca, csc, ec, endEarly)
}

func (c *chatServicesComponent) updateChatServices(ca *account, csc <-chan jid.Domain, ec <-chan error, endEarly func()) {
	hadAny := false
	ts := make(chan jid.Domain)

	doInUIThread(func() {
		t := c.currentService()
		ts <- t
		c.removeAll()
	})

	typedService := <-ts

	defer func() {
		c.onUpdateChatServicesFinished(hadAny, typedService)
	}()

	for {
		select {
		case <-c.previousUpdateChannel:
			c.removeAll()
			endEarly()
			return
		case err, _ := <-ec:
			if err != nil {
				ca.log.WithError(err).Error("Something went wrong trying to get chat services")
			}
			return
		case cs, ok := <-csc:
			if !ok {
				return
			}

			hadAny = true
			doInUIThread(func() {
				c.chatServices.AppendText(cs.String())
			})
		}
	}
}

func (c *chatServicesComponent) onUpdateChatServicesFinished(hadAny bool, typedService jid.Domain) {
	if hadAny && len(typedService.String()) == 0 {
		c.setActive(0)
	}

	c.previousUpdateChannel = nil
}

func (c *chatServicesComponent) getView() gtki.Box {
	return c.chatServicesBox
}

func (c *chatServicesComponent) currentService() jid.Domain {
	cs, _ := c.chatServiceEntry.GetText()
	return jid.ParseDomain(cs)
}

func (c *chatServicesComponent) setActive(index int) {
	doInUIThread(func() {
		c.chatServices.SetActive(index)
	})
}

func (c *chatServicesComponent) removeAll() {
	doInUIThread(c.chatServices.RemoveAll)
}

// enableServiceInput MUST be called from the UI thread
func (c *chatServicesComponent) enableServiceInput() {
	c.chatServices.SetSensitive(true)
}

// disableServiceInput MUST be called from the UI thread
func (c *chatServicesComponent) disableServiceInput() {
	c.chatServices.SetSensitive(false)
}

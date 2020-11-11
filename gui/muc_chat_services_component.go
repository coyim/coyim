package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type chatServicesComponent struct {
	chatServicesBox       gtki.Box
	chatServicesList      gtki.ComboBoxText
	chatServiceEntry      gtki.Entry
	previousUpdateChannel chan bool
}

func (u *gtkUI) createChatServicesComponent(list gtki.ComboBoxText, entry gtki.Entry, onServiceChanged func()) *chatServicesComponent {
	c := &chatServicesComponent{}

	c.chatServiceEntry = entry

	c.chatServicesList = list
	c.chatServicesList.Connect("changed", func() {
		if onServiceChanged != nil {
			onServiceChanged()
		}
	})

	return c
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
			doInUIThread(c.removeAll)
			endEarly()
			return
		case err := <-ec:
			if err != nil {
				ca.log.WithError(err).Error("Something went wrong trying to get the available chat services")
			}
			return
		case cs, ok := <-csc:
			if !ok {
				return
			}

			hadAny = true
			doInUIThread(func() {
				c.addService(cs)
			})
		}
	}
}

func (c *chatServicesComponent) onUpdateChatServicesFinished(hadAny bool, typedService jid.Domain) {
	if hadAny && len(typedService.String()) == 0 {
		doInUIThread(func() {
			c.setActive(0)
		})
	}

	c.previousUpdateChannel = nil
}

// currentService MUST be called from the UI thread
func (c *chatServicesComponent) currentService() jid.Domain {
	cs, _ := c.chatServiceEntry.GetText()
	return jid.ParseDomain(cs)
}

// setActive MUST be called from the UI thread
func (c *chatServicesComponent) setActive(index int) {
	c.chatServicesList.SetActive(index)
}

// addService MUST be called from the UI thread
func (c *chatServicesComponent) addService(s jid.Domain) {
	c.chatServicesList.AppendText(s.String())
}

// removeAll MUST be called from the UI thread
func (c *chatServicesComponent) removeAll() {
	c.chatServicesList.RemoveAll()
}

// enableServiceInput MUST be called from the UI thread
func (c *chatServicesComponent) enableServiceInput() {
	c.chatServicesList.SetSensitive(true)
}

// disableServiceInput MUST be called from the UI thread
func (c *chatServicesComponent) disableServiceInput() {
	c.chatServicesList.SetSensitive(false)
}

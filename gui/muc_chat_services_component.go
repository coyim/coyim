package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	chatServicesModelIDColumn int = iota
	chatServicesModelTextColumn
)

type chatServicesComponent struct {
	currentAccount        *account
	currentValue          string
	services              map[int]string
	servicesList          gtki.ComboBoxText
	servicesListModel     gtki.ListStore
	serviceEntry          gtki.Entry
	previousUpdateChannel chan bool
}

func (u *gtkUI) createChatServicesComponent(list gtki.ComboBoxText, entry gtki.Entry, onServiceChanged func()) *chatServicesComponent {
	c := &chatServicesComponent{
		serviceEntry: entry,
		services:     make(map[int]string),
	}

	var err error
	// The following creates a list store model with two columns
	// one for the "ID" and the another for the "text"
	c.servicesListModel, err = g.gtk.ListStoreNew(glibi.TYPE_STRING, glibi.TYPE_STRING)
	if err != nil {
		panic(err)
	}

	onServiceChangedFinal := onServiceChanged
	onServiceChanged = func() {
		if currentValue, ok := c.services[c.servicesList.GetActive()]; ok {
			c.currentValue = currentValue
		} else {
			c.currentValue = c.servicesList.GetActiveText()
		}

		if onServiceChangedFinal != nil {
			onServiceChangedFinal()
		}
	}

	c.servicesList = list
	c.servicesList.Connect("changed", onServiceChanged)

	c.servicesList.SetModel(c.servicesListModel)
	c.servicesList.SetIDColumn(chatServicesModelIDColumn)
	c.servicesList.SetEntryTextColumn(chatServicesModelTextColumn)

	return c
}

func (c *chatServicesComponent) updateServicesBasedOnAccount(ca *account) {
	if c.currentAccount == nil || c.currentAccount.Account() != ca.Account() {
		c.currentAccount = ca

		if c.previousUpdateChannel != nil {
			c.previousUpdateChannel <- true
		}

		c.previousUpdateChannel = make(chan bool)

		csc, ec, endEarly := ca.session.GetChatServices(jid.ParseDomain(ca.Account()))

		go c.updateChatServices(ca, csc, ec, endEarly)
	}
}

func (c *chatServicesComponent) updateChatServices(ca *account, csc <-chan jid.Domain, ec <-chan error, endEarly func()) {
	doInUIThread(c.removeAll)

	defer func() {
		c.previousUpdateChannel = nil
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

			doInUIThread(func() {
				c.addService(cs)
				if c.currentValue == "" {
					c.setActive(0)
				}
			})
		}
	}
}

// currentServiceValue MUST be called from the UI thread
func (c *chatServicesComponent) currentServiceValue() string {
	return c.currentValue
}

// currentService MUST be called from the UI thread
func (c *chatServicesComponent) currentService() jid.Domain {
	return jid.ParseDomain(c.currentServiceValue())
}

// setCurrentService MUST be called from the UI thread
func (c *chatServicesComponent) setCurrentService(s jid.Domain) {
	for i, ss := range c.services {
		if ss == s.String() {
			c.setActive(i)
			return
		}
	}
	c.serviceEntry.SetText(s.String())
}

// setActive MUST be called from the UI thread
func (c *chatServicesComponent) setActive(index int) {
	if len(c.services) > 0 && len(c.services) < index {
		c.servicesList.SetActive(index)
	}
}

// addService MUST be called from the UI thread
func (c *chatServicesComponent) addService(s jid.Domain) {
	iter := c.servicesListModel.Append()

	_ = c.servicesListModel.SetValue(iter, chatServicesModelIDColumn, s.String())
	_ = c.servicesListModel.SetValue(iter, chatServicesModelTextColumn, s.String())

	c.services[len(c.services)] = s.String()
}

// removeAll MUST be called from the UI thread
func (c *chatServicesComponent) removeAll() {
	c.currentValue = ""
	c.services = make(map[int]string)
	c.servicesList.RemoveAll()
}

// enableServiceInput MUST be called from the UI thread
func (c *chatServicesComponent) enableServiceInput() {
	c.servicesList.SetSensitive(true)
}

// disableServiceInput MUST be called from the UI thread
func (c *chatServicesComponent) disableServiceInput() {
	c.servicesList.SetSensitive(false)
}

// resetToDefault MUST be called from the UI thread
func (c *chatServicesComponent) resetToDefault() {
	c.serviceEntry.SetText("")
	if len(c.services) > 0 {
		c.setActive(0)
	}
}

// hasServiceValue MUST be called from the UI thread
func (c *chatServicesComponent) hasServiceValue() bool {
	return c.currentServiceValue() != ""
}

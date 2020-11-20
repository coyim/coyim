package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

// This is a component that can be used from other dialogs in order to make it easy to have a
// list of all connected accounts, and having it updated automatically and so on

type connectedAccountsComponent struct {
	accountsModel      gtki.ListStore
	accountsInput      gtki.ComboBox
	accountsList       []*account
	accounts           map[string]*account
	currentlyActive    int
	onDestroyFunc      func()
	errorNotifications canNotifyErrors
	onNoAccounts       func()
	onAccountsUpdated  func(*account)
}

// disableAccountInput should ONLY be called from the UI thread
func (cac *connectedAccountsComponent) disableAccountInput() {
	cac.accountsInput.SetSensitive(false)
}

// enableAccountInput should ONLY be called from the UI thread
func (cac *connectedAccountsComponent) enableAccountInput() {
	cac.accountsInput.SetSensitive(true)
}

// onDestroy should ONLY be called from the UI thread
func (cac *connectedAccountsComponent) onDestroy() {
	if cac.onDestroyFunc != nil {
		cac.onDestroyFunc()
	}
}

// currentAccount is safe to call from the UI thread, or from outside the UI thread
func (cac *connectedAccountsComponent) currentAccount() *account {
	if cac.currentlyActive != -1 && len(cac.accountsList) > 0 && cac.currentlyActive < len(cac.accountsList) {
		return cac.accountsList[cac.currentlyActive]
	}
	return nil
}

// hasAccounts is safe to call from the UI thread, or from outside the UI thread
func (cac *connectedAccountsComponent) hasAccounts() bool {
	return len(cac.accounts) > 0
}

// initOrReplaceAccounts should ONLY be called from the UI thread
// We don't need to do any locking here, since it's already in the UI thread
// so no other things can happen to the UI in the meantime
func (cac *connectedAccountsComponent) initOrReplaceAccounts(accounts []*account) {
	if len(accounts) == 0 {
		cac.errorNotifications.notifyOnError(i18n.Local("There are no connected accounts"))
	}

	currentlyActive := 0
	oldActive := cac.accountsInput.GetActive()
	if cac.accounts != nil && oldActive >= 0 {
		currentlyActiveAccount := cac.accountsList[oldActive]
		for ix, a := range accounts {
			if currentlyActiveAccount == a {
				currentlyActive = ix
				cac.currentlyActive = currentlyActive
			}
		}
		cac.accountsModel.Clear()
	}

	cac.accountsList = accounts
	cac.accounts = make(map[string]*account)
	for _, acc := range accounts {
		iter := cac.accountsModel.Append()
		_ = cac.accountsModel.SetValue(iter, 0, acc.Account())
		_ = cac.accountsModel.SetValue(iter, 1, acc.ID())
		cac.accounts[acc.ID()] = acc
	}

	if len(accounts) > 0 {
		cac.accountsInput.SetActive(currentlyActive)
		cac.errorNotifications.clearErrors()
		cac.onAccountsUpdated(cac.currentAccount())
	} else {
		go cac.onNoAccounts()
	}
}

// createConnectedAccountsComponent will create a component for managing a combo box of connected accounts.
// It should ONLY be called from the UI thread
// it needs the combo box to populate. It will create its own model to populate with accounts.
// The errorNot argument is anything that can notify on errors, and clear those notifications. The onAccountsUpdated callback
// will be called when the the user moves to a different account in the list. The onNoAccounts callback will be called if
// there are no connected accounts available at all.
// Remember to call onDestroy() on the component when the window is destroyed, to not create to many update notifications that lead
// nowhere.
// The onAccountsUpdated function will NOT be called from the UI thread.
// The onNoAccounts function will NOT be called from the UI thread
func (u *gtkUI) createConnectedAccountsComponent(input gtki.ComboBox, errorNot canNotifyErrors, onAccountsUpdated func(*account), onNoAccounts func()) *connectedAccountsComponent {
	result := &connectedAccountsComponent{}
	result.errorNotifications = errorNot
	result.onNoAccounts = onNoAccounts

	onAccountsUpdatedFinal := onAccountsUpdated
	if onAccountsUpdatedFinal == nil {
		onAccountsUpdatedFinal = func(*account) {}
	}

	result.onAccountsUpdated = onAccountsUpdatedFinal

	var e error
	// These two arguments are: account name and account id
	result.accountsModel, e = g.gtk.ListStoreNew(glibi.TYPE_STRING, glibi.TYPE_STRING)
	if e != nil {
		panic(e)
	}

	result.accountsInput = input
	result.accountsInput.SetModel(result.accountsModel)

	result.initOrReplaceAccounts(u.getAllConnectedAccounts())

	accountsObserverToken := u.onChangeOfConnectedAccounts(func() {
		doInUIThread(func() {
			result.initOrReplaceAccounts(u.getAllConnectedAccounts())
		})
	})

	_, _ = result.accountsInput.Connect("changed", func() {
		act := result.accountsInput.GetActive()
		if act >= 0 && act < len(result.accountsList) && act != result.currentlyActive {
			result.currentlyActive = act
			go result.onAccountsUpdated(result.currentAccount())
		}
	})

	result.currentlyActive = -1
	if len(result.accountsList) > 0 {
		result.currentlyActive = 0
		result.onAccountsUpdated(result.currentAccount())
	}

	result.onDestroyFunc = func() {
		u.removeConnectedAccountsObserver(accountsObserverToken)
	}

	return result
}

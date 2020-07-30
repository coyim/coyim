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
}

func (cac *connectedAccountsComponent) onDestroy() {
	if cac.onDestroyFunc != nil {
		cac.onDestroyFunc()
	}
}

func (cac *connectedAccountsComponent) currentAccount() *account {
	return cac.accountsList[cac.currentlyActive]
}

// initOrReplaceAccounts should be called from the UI thread
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
	} else {
		cac.onNoAccounts()
	}
}

// createConnectedAccountsComponent will create a component for managing a combo box of connected accounts.
// it needs the combo box to populate. It will create its own model to populate with accounts.
// The errorNot argument is anything that can notify on errors, and clear those notifications. The onAccountsUpdated callback
// will be called when the the user moves to a different account in the list. The onNoAccounts callback will be called if
// there are no connected accounts available at all.
// Remember to call onDestroy() on the component when the window is destroyed, to not create to many update notifications that lead
// nowhere.
func (u *gtkUI) createConnectedAccountsComponent(input gtki.ComboBox, errorNot canNotifyErrors, onAccountsUpdated func(*account), onNoAccounts func()) *connectedAccountsComponent {
	result := &connectedAccountsComponent{}
	result.errorNotifications = errorNot
	result.onNoAccounts = onNoAccounts

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
			onAccountsUpdated(result.currentAccount())
		}
	})

	result.currentlyActive = -1
	if len(result.accountsList) > 0 {
		result.currentlyActive = 0
		onAccountsUpdated(result.currentAccount())
	}

	result.onDestroyFunc = func() {
		u.removeConnectedAccountsObserver(accountsObserverToken)
	}

	return result
}

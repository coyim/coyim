package gui

import (
	"log"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type editContactDialog struct {
	builder           *builder
	dialog            gtki.Dialog
	accountName       gtki.Label
	contactJID        gtki.Label
	nickname          gtki.Entry
	requireEncryption gtki.CheckButton
	currentGroups     gtki.ListStore
	existingGroups    gtki.Menu
	addGroup          gtki.MenuItem
	removeGroup       gtki.Button
	currentGroupsView gtki.TreeView
	save              gtki.Button
}

func (ecd *editContactDialog) init() {
	ecd.builder = newBuilder("PeerDetails")
	ecd.builder.getItems(
		"dialog", &ecd.dialog,
		"account-name", &ecd.accountName,
		"jid", &ecd.contactJID,
		"nickname", &ecd.nickname,
		"require-encryption", &ecd.requireEncryption,
		"current-groups", &ecd.currentGroups,
		"groups-menu", &ecd.existingGroups,
		"addGroup", &ecd.addGroup,
		"remove-btn", &ecd.removeGroup,
		"groups-view", &ecd.currentGroupsView,
		"btn-save", &ecd.save,
	)
}

func (ecd *editContactDialog) addCurrentGroup(g string) {
	ecd.currentGroups.SetValue(ecd.currentGroups.Append(), 0, g)
}

func (ecd *editContactDialog) initAllGroups(allGroups []string) {
	for _, gr := range allGroups {
		menu, _ := g.gtk.MenuItemNewWithLabel(gr)
		menu.SetVisible(true)
		menu.Connect("activate", func(m gtki.MenuItem) {
			ecd.addCurrentGroup(m.GetLabel())
		})

		ecd.existingGroups.Add(menu)
	}

	if len(allGroups) > 0 {
		sep, _ := g.gtk.SeparatorMenuItemNew()
		sep.SetVisible(true)
		ecd.existingGroups.Add(sep)
	}
}

func (ecd *editContactDialog) initCurrentGroups(groups []string) {
	ecd.currentGroups.Clear()

	for _, g := range groups {
		ecd.addCurrentGroup(g)
	}
}

func (r *roster) openEditContactDialog(jid string, acc *account) {
	assertInUIThread()
	peer, ok := r.ui.accountManager.getPeer(acc, jid)
	if !ok {
		log.Printf("Couldn't find peer %s in account %v", jid, acc)
		return
	}

	ecd := &editContactDialog{}
	ecd.init()
	conf := acc.session.GetConfig()

	ecd.accountName.SetText(conf.Account)
	ecd.contactJID.SetText(jid)

	//nickNameEntry.SetText(peer.Name)
	if peer, ok := r.ui.getPeer(acc, jid); ok {
		ecd.nickname.SetText(peer.Nickname)
	}

	shouldEncryptTo := conf.ShouldEncryptTo(jid)
	ecd.requireEncryption.SetActive(shouldEncryptTo)

	ecd.initCurrentGroups(sortedGroupNames(peer.Groups))
	ecd.initAllGroups(r.getGroupNamesFor(acc))

	ecd.existingGroups.Add(ecd.addGroup)
	ecd.removeGroup.SetSensitive(false)

	ecd.builder.ConnectSignals(map[string]interface{}{
		"on-add-new-group": func() {
			r.addGroupDialog(ecd.currentGroups)
		},
		"on-group-selection-changed": func() {
			ts, _ := ecd.currentGroupsView.GetSelection()
			_, _, ok := ts.GetSelected()
			ecd.removeGroup.SetSensitive(ok)
		},
		"on-remove-group": func() {
			ts, _ := ecd.currentGroupsView.GetSelection()
			if _, iter, ok := ts.GetSelected(); ok {
				ecd.currentGroups.Remove(iter)
			}
		},
		"on-cancel": ecd.dialog.Destroy,
		"on-save": func() {
			assertInUIThread()
			defer ecd.dialog.Destroy()

			groups := toArray(ecd.currentGroups)
			nickname, _ := ecd.nickname.GetText()

			err := r.updatePeer(acc, jid, nickname, groups, ecd.requireEncryption.GetActive() != shouldEncryptTo, ecd.requireEncryption.GetActive())
			if err != nil {
				log.Println(err)
			}
		},
	})

	ecd.save.GrabDefault()
	ecd.dialog.SetTransientFor(r.ui.window)
	ecd.dialog.ShowAll()
}

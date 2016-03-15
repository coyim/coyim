package gui

import (
	"log"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

func (r *roster) openEditContactDialog(jid string, acc *account) {
	peer, ok := r.ui.accountManager.getPeer(acc, jid)
	if !ok {
		log.Printf("Couldn't find peer %s in account %v", jid, acc)
		return
	}

	builder := newBuilder("PeerDetails")
	dialog := builder.getObj("dialog").(gtki.Dialog)

	conf := acc.session.GetConfig()
	accName := builder.getObj("account-name").(gtki.Label)
	accName.SetText(conf.Account)

	contactJID := builder.getObj("jid").(gtki.Label)
	contactJID.SetText(jid)

	nickNameEntry := builder.getObj("nickname").(gtki.Entry)
	//nickNameEntry.SetText(peer.Name)
	if peer, ok := r.ui.getPeer(acc, jid); ok {
		nickNameEntry.SetText(peer.Nickname)
	}

	requireEncryptionEntry := builder.getObj("require-encryption").(gtki.CheckButton)
	shouldEncryptTo := conf.ShouldEncryptTo(jid)
	requireEncryptionEntry.SetActive(shouldEncryptTo)

	currentGroups := builder.getObj("current-groups").(gtki.ListStore)
	currentGroups.Clear()

	for n := range peer.Groups {
		currentGroups.SetValue(currentGroups.Append(), 0, n)
	}

	existingGroups := builder.getObj("groups-menu").(gtki.Menu)
	allGroups := r.allGroupNames()
	for _, gr := range allGroups {
		menu, err := g.gtk.MenuItemNewWithLabel(gr)
		if err != nil {
			continue
		}

		menu.SetVisible(true)
		menu.Connect("activate", func(m gtki.MenuItem) {
			currentGroups.SetValue(currentGroups.Append(), 0, m.GetLabel())
		})

		existingGroups.Add(menu)
	}

	if len(allGroups) > 0 {
		sep, err := g.gtk.SeparatorMenuItemNew()
		if err != nil {
			return
		}

		sep.SetVisible(true)
		existingGroups.Add(sep)
	}

	addMenuItem := builder.getObj("addGroup").(gtki.MenuItem)
	existingGroups.Add(addMenuItem)

	currentGroupsView := builder.getObj("groups-view").(gtki.TreeView)

	builder.ConnectSignals(map[string]interface{}{
		"on-add-new-group": func() {
			r.addGroupDialog(currentGroups)
		},
		"on-remove-group": func() {
			path, err := currentGroupsView.GetCursor()
			if err != nil {
				return
			}

			iter, err2 := currentGroups.GetIter(path)
			if err2 != nil {
				return
			}

			currentGroups.Remove(iter)
		},
		"on-cancel": dialog.Destroy,
		"on-save": func() {
			defer dialog.Destroy()

			groups := toArray(currentGroups)
			nickname, _ := nickNameEntry.GetText()

			err := r.updatePeer(acc, jid, nickname, groups, requireEncryptionEntry.GetActive() != shouldEncryptTo, requireEncryptionEntry.GetActive())
			if err != nil {
				log.Println(err)
			}
		},
	})

	defaultBtn := builder.getObj("btn-save").(gtki.Button)
	defaultBtn.GrabDefault()
	dialog.SetTransientFor(r.ui.window)
	dialog.ShowAll()
}

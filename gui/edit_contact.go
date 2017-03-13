package gui

import (
	"fmt"
	"log"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3adapter/gtki"
)

type editContactDialog struct {
	builder                 *builder
	dialog                  gtki.Dialog
	accountName             gtki.Label
	contactJID              gtki.Label
	nickname                gtki.Entry
	requireEncryption       gtki.CheckButton
	currentGroups           gtki.ListStore
	existingGroups          gtki.Menu
	addGroup                gtki.MenuItem
	removeGroup             gtki.Button
	currentGroupsView       gtki.TreeView
	fingerprintsInformation gtki.Label
	fingerprintsGrid        gtki.Grid
	save                    gtki.Button
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
		"fingerprintsInformation", &ecd.fingerprintsInformation,
		"fingerprintsGrid", &ecd.fingerprintsGrid,
		"btn-save", &ecd.save,
	)

}

func (ecd *editContactDialog) closeDialog() {
	ecd.dialog.Destroy()
}

func (ecd *editContactDialog) enableRemoveButton() {
	ts, _ := ecd.currentGroupsView.GetSelection()
	_, _, ok := ts.GetSelected()
	ecd.removeGroup.SetSensitive(ok)
}

func (ecd *editContactDialog) removeSelectedGroup() {
	ts, _ := ecd.currentGroupsView.GetSelection()
	if _, iter, ok := ts.GetSelected(); ok {
		ecd.currentGroups.Remove(iter)
	}
}

func (ecd *editContactDialog) openAddGroupDialog() {
	groupList := ecd.currentGroups
	builder := newBuilder("GroupDetails")
	dialog := builder.getObj("dialog").(gtki.Dialog)

	nameEntry := builder.getObj("group-name").(gtki.Entry)

	defaultBtn := builder.getObj("btn-ok").(gtki.Button)
	defaultBtn.GrabDefault()

	dialog.SetTransientFor(ecd.dialog)
	dialog.ShowAll()

	response := dialog.Run()
	defer dialog.Destroy()

	if gtki.ResponseType(response) != gtki.RESPONSE_OK {
		return
	}

	groupName, _ := nameEntry.GetText()
	groupList.SetValue(groupList.Append(), 0, groupName)
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

func (ecd *editContactDialog) showFingerprintsForPeer(jid string, account *account) {
	info := ecd.fingerprintsInformation
	grid := ecd.fingerprintsGrid

	fprs := []*config.Fingerprint{}
	p, ok := account.session.GetConfig().GetPeer(jid)
	if ok {
		fprs = p.Fingerprints
	}

	if len(fprs) == 0 {
		info.SetText(fmt.Sprintf(i18n.Local("There are no known fingerprints for %s"), jid))
	} else {
		info.SetText(fmt.Sprintf(i18n.Local("These are the fingerprints known for %s:"), jid))
	}

	for ix, fpr := range fprs {
		flabel, _ := g.gtk.LabelNew(config.FormatFingerprint(fpr.Fingerprint))
		flabel.SetSelectable(true)
		trusted := i18n.Local("not trusted")
		if fpr.Trusted {
			trusted = i18n.Local("trusted")
		}

		ftrusted, _ := g.gtk.LabelNew(trusted)
		ftrusted.SetSelectable(true)

		grid.Attach(flabel, 0, ix, 1, 1)
		grid.Attach(ftrusted, 1, ix, 1, 1)
	}
}

func (r *roster) openEditContactDialog(jid string, acc *account) {
	assertInUIThread()
	peer, ok := r.ui.accountManager.getPeer(acc, jid)
	if !ok {
		log.Printf("Couldn't find peer %s in account %v", jid, acc)
		return
	}

	conf := acc.session.GetConfig()

	ecd := &editContactDialog{}
	ecd.init()
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
	ecd.showFingerprintsForPeer(jid, acc)

	//TODO: Move to editContactDialog
	ecd.builder.ConnectSignals(map[string]interface{}{
		"on-add-new-group":           ecd.openAddGroupDialog,
		"on-group-selection-changed": ecd.enableRemoveButton,
		"on-remove-group":            ecd.removeSelectedGroup,
		"on-cancel":                  ecd.closeDialog,

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

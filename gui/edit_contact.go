package gui

import (
	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type editContactDialog struct {
	builder                 *builder
	dialog                  gtki.Dialog      `gtk-widget:"dialog"`
	accountName             gtki.Label       `gtk-widget:"account-name"`
	contactJID              gtki.Label       `gtk-widget:"jid"`
	nickname                gtki.Entry       `gtk-widget:"nickname"`
	requireEncryption       gtki.CheckButton `gtk-widget:"require-encryption"`
	currentGroups           gtki.ListStore   `gtk-widget:"current-groups"`
	existingGroups          gtki.Menu        `gtk-widget:"groups-menu"`
	addGroup                gtki.MenuItem    `gtk-widget:"addGroup"`
	removeGroup             gtki.Button      `gtk-widget:"remove-btn"`
	currentGroupsView       gtki.TreeView    `gtk-widget:"groups-view"`
	fingerprintsInformation gtki.Label       `gtk-widget:"fingerprintsInformation"`
	fingerprintsGrid        gtki.Grid        `gtk-widget:"fingerprintsGrid"`
	save                    gtki.Button      `gtk-widget:"btn-save"`
}

func (ecd *editContactDialog) init() {
	ecd.builder = newBuilder("PeerDetails")
	panicOnDevError(ecd.builder.bindObjects(ecd))
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
	_ = groupList.SetValue(groupList.Append(), 0, groupName)
}

func (ecd *editContactDialog) addCurrentGroup(g string) {
	_ = ecd.currentGroups.SetValue(ecd.currentGroups.Append(), 0, g)
}

func (ecd *editContactDialog) initAllGroups(allGroups []string) {
	for _, gr := range allGroups {
		menu, _ := g.gtk.MenuItemNewWithLabel(gr)
		menu.SetVisible(true)
		_, _ = menu.Connect("activate", func(m gtki.MenuItem) {
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

func (ecd *editContactDialog) showFingerprintsForPeer(peer jid.WithoutResource, account *account) {
	info := ecd.fingerprintsInformation
	grid := ecd.fingerprintsGrid

	fprs := []*config.Fingerprint{}
	p, ok := account.session.GetConfig().GetPeer(peer.String())
	if ok {
		fprs = p.Fingerprints
	}

	if len(fprs) == 0 {
		info.SetText(i18n.Localf("There are no known fingerprints for %s", peer))
	} else {
		info.SetText(i18n.Localf("These are the fingerprints known for %s:", peer))
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

func (r *roster) openEditContactDialog(peer jid.WithoutResource, acc *account) {
	assertInUIThread()
	p, ok := r.ui.accountManager.getPeer(acc, peer)
	if !ok {
		acc.log.WithField("peer", peer).Warn("Couldn't find peer in account")
		return
	}

	conf := acc.session.GetConfig()

	ecd := &editContactDialog{}
	ecd.init()
	ecd.accountName.SetText(conf.Account)
	ecd.contactJID.SetText(peer.String())

	//nickNameEntry.SetText(peer.Name)
	if p, ok := r.ui.getPeer(acc, peer); ok {
		ecd.nickname.SetText(p.Nickname)
	}

	shouldEncryptTo := conf.ShouldEncryptTo(peer.String())
	ecd.requireEncryption.SetActive(shouldEncryptTo)

	ecd.initCurrentGroups(sortedGroupNames(p.Groups))
	ecd.initAllGroups(r.getGroupNamesFor(acc))

	ecd.existingGroups.Add(ecd.addGroup)
	ecd.removeGroup.SetSensitive(false)
	ecd.showFingerprintsForPeer(peer, acc)

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

			err := r.updatePeer(acc, peer, nickname, groups, ecd.requireEncryption.GetActive() != shouldEncryptTo, ecd.requireEncryption.GetActive())
			if err != nil {
				acc.log.WithError(err).Warn("Something went wrong")
			}
		},
	})

	ecd.save.GrabDefault()
	ecd.dialog.SetTransientFor(r.ui.window)
	ecd.dialog.ShowAll()
}

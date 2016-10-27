package gui

import (
	"fmt"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3adapter/gtki"
)

func (u *gtkUI) showFingerprintsForPeer(jid string, account *account) {
	builder := newBuilder("PeerFingerprints")
	dialog := builder.getObj("dialog").(gtki.Dialog)
	info := builder.getObj("information").(gtki.Label)
	grid := builder.getObj("grid").(gtki.Grid)

	info.SetSelectable(true)

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

	builder.ConnectSignals(map[string]interface{}{
		"on_close_signal": func() {
			dialog.Destroy()
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

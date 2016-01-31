package gui

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/config/importer"
	"github.com/twstrike/coyim/i18n"
)

func valAt(s *gtk.ListStore, iter *gtk.TreeIter, col int) interface{} {
	gv, _ := s.GetValue(iter, col)
	vv, _ := gv.GoValue()
	return vv
}

type applicationAndAccount struct {
	app string
	acc string
}

func (u *gtkUI) doActualImportOf(choices map[applicationAndAccount]bool, potential map[string][]*config.ApplicationConfig) {
	for k, v := range choices {
		if !v {
			continue
		}

		for _, accs := range potential[k.app] {
			for _, a := range accs.Accounts {
				if a.Account == k.acc {
					log.Printf("[import] Doing import of %s from %s", k.acc, k.app)
					accountToImport := a
					u.config.WhenLoaded(func(conf *config.ApplicationConfig) {
						_, exists := conf.GetAccount(k.acc)
						if exists {
							// TODO: view message
							log.Printf("[import] Can't import account %s since you already have an account configured with the same name. Remove that account and import again if you really want to overwrite it.", k.acc)
							return
						}

						if conf.RawLogFile == "" {
							conf.RawLogFile = accs.RawLogFile
						}
						if len(conf.NotifyCommand) == 0 {
							conf.NotifyCommand = accs.NotifyCommand
						}
						if conf.IdleSecondsBeforeNotification == 0 {
							conf.IdleSecondsBeforeNotification = accs.IdleSecondsBeforeNotification
						}
						if !conf.Bell {
							conf.Bell = accs.Bell
						}

						u.addAndSaveAccountConfig(accountToImport)
					})
				}
			}
		}
	}
}

func (u *gtkUI) runImporter() {
	importSettings := make(map[applicationAndAccount]bool)
	allImports := importer.TryImportAll()

	builder := builderForDefinition("Importer")

	win, _ := builder.GetObject("importerWindow")
	w := win.(*gtk.Dialog)

	store, _ := builder.GetObject("importAccountsStore")
	s := store.(*gtk.ListStore)

	for appName, v := range allImports {
		for _, vv := range v {
			for _, a := range vv.Accounts {
				it := s.Append()
				s.SetValue(it, 0, appName)
				s.SetValue(it, 1, a.Account)
				s.SetValue(it, 2, false)
			}
		}
	}

	rend, _ := builder.GetObject("import-this-account-renderer")
	rr := rend.(*gtk.CellRendererToggle)

	rr.Connect("toggled", func(_ interface{}, path string) {
		iter, _ := s.GetIterFromString(path)
		current, _ := valAt(s, iter, 2).(bool)
		app, _ := valAt(s, iter, 0).(string)
		acc, _ := valAt(s, iter, 1).(string)

		importSettings[applicationAndAccount{app, acc}] = !current

		s.SetValue(iter, 2, !current)
	})

	w.Connect("response", func(_ interface{}, rid int) {
		if gtk.ResponseType(rid) == gtk.RESPONSE_OK {
			u.doActualImportOf(importSettings, allImports)
		}
		w.Destroy()
	})

	u.connectShortcutsChildWindow(&w.Window)
	doInUIThread(func() {
		w.SetTransientFor(u.window)
		w.ShowAll()
	})
}

func (u *gtkUI) importFingerprintsFor(account *config.Account, file string) (int, bool) {
	fprs, ok := importer.ImportFingerprintsFromPidginStyle(file, func(string) bool { return true })
	if !ok {
		return 0, false
	}

	num := 0
	for _, kfprs := range fprs {
		for _, kfpr := range kfprs {
			log.Printf("Importing fingerprint %X for %s", kfpr.Fingerprint, kfpr.UserID)
			fpr := account.EnsurePeer(kfpr.UserID).EnsureHasFingerprint(kfpr.Fingerprint)
			num = num + 1
			if !kfpr.Untrusted {
				fpr.Trusted = true
			}
		}
	}

	return num, true
}

func (u *gtkUI) importKeysFor(account *config.Account, file string) (int, bool) {
	keys, ok := importer.ImportKeysFromPidginStyle(file, func(string) bool { return true })
	if !ok {
		return 0, false
	}

	newKeys := [][]byte{}
	for _, kk := range keys {
		newKeys = append(newKeys, kk)
	}
	account.PrivateKeys = newKeys

	return len(newKeys), true
}

func (u *gtkUI) exportFingerprintsFor(account *config.Account, file string) bool {
	return true
}

func (u *gtkUI) exportKeysFor(account *config.Account, file string) bool {
	return true
}

func (u *gtkUI) importKeysForDialog(account *config.Account, w *gtk.Dialog) {
	dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
		i18n.Local("Import private keys"),
		&w.Window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		i18n.Local("_Cancel"),
		gtk.RESPONSE_CANCEL,
		i18n.Local("_Import"),
		gtk.RESPONSE_OK,
	)

	if gtk.ResponseType(dialog.Run()) == gtk.RESPONSE_OK {
		num, ok := u.importKeysFor(account, dialog.GetFilename())
		if ok {
			u.notify(i18n.Local("Keys imported"), fmt.Sprintf(i18n.Local("%d key(s) were imported correctly."), num))
		} else {
			u.notify(i18n.Local("Failure importing keys"), fmt.Sprintf(i18n.Local("Couldn't import any keys from %s."), dialog.GetFilename()))
		}
	}
	dialog.Destroy()
}

func (u *gtkUI) exportKeysForDialog(account *config.Account, w *gtk.Dialog) {
	dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
		i18n.Local("Export private keys"),
		&w.Window,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		i18n.Local("_Cancel"),
		gtk.RESPONSE_CANCEL,
		i18n.Local("_Export"),
		gtk.RESPONSE_OK,
	)

	dialog.SetCurrentName("otr.private_key")

	if gtk.ResponseType(dialog.Run()) == gtk.RESPONSE_OK {
		u.exportKeysFor(account, dialog.GetFilename())
	}
	dialog.Destroy()
}

func (u *gtkUI) importFingerprintsForDialog(account *config.Account, w *gtk.Dialog) {
	dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
		i18n.Local("Import fingerprints"),
		&w.Window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		i18n.Local("_Cancel"),
		gtk.RESPONSE_CANCEL,
		i18n.Local("_Import"),
		gtk.RESPONSE_OK,
	)

	if gtk.ResponseType(dialog.Run()) == gtk.RESPONSE_OK {
		num, ok := u.importFingerprintsFor(account, dialog.GetFilename())
		if ok {
			u.notify(i18n.Local("Fingerprints imported"), fmt.Sprintf(i18n.Local("%d fingerprint(s) were imported correctly."), num))
		} else {
			u.notify(i18n.Local("Failure importing fingerprints"), fmt.Sprintf(i18n.Local("Couldn't import any fingerprints from %s."), dialog.GetFilename()))
		}
	}
	dialog.Destroy()
}

func (u *gtkUI) exportFingerprintsForDialog(account *config.Account, w *gtk.Dialog) {
	dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
		i18n.Local("Export fingerprints"),
		&w.Window,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		i18n.Local("_Cancel"),
		gtk.RESPONSE_CANCEL,
		i18n.Local("_Export"),
		gtk.RESPONSE_OK,
	)

	dialog.SetCurrentName("otr.fingerprints")

	if gtk.ResponseType(dialog.Run()) == gtk.RESPONSE_OK {
		u.exportKeysFor(account, dialog.GetFilename())
	}
	dialog.Destroy()
}

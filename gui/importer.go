package gui

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/config/importer"
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
		if v {
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
}

func (u *gtkUI) runImporter() {
	importSettings := make(map[applicationAndAccount]bool)
	allImports := importer.TryImportAll()

	builder, _ := loadBuilderWith("Importer")

	win, _ := builder.GetObject("importerWindow")
	w, _ := win.(*gtk.Dialog)

	store, _ := builder.GetObject("importAccountsStore")
	s, _ := store.(*gtk.ListStore)

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
	rr, _ := rend.(*gtk.CellRendererToggle)

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
	glib.IdleAdd(func() {
		w.SetTransientFor(u.window)
		w.ShowAll()
	})
}

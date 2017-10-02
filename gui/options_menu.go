package gui

import (
	"github.com/coyim/coyim/config"
	"github.com/coyim/gotk3adapter/gtki"
)

type optionsMenu struct {
	encryptConfig gtki.CheckMenuItem
}

func (v *optionsMenu) setFromConfig(c *config.ApplicationConfig) {
	doInUIThread(func() {
		v.encryptConfig.SetActive(c.HasEncryptedStorage())
	})
}

func (u *gtkUI) toggleEncryptedConfig() {
	if u.config != nil {
		val := u.optionsMenu.encryptConfig.GetActive()
		if u.config.SetShouldSaveFileEncrypted(val) {
			if val {
				u.captureInitialMasterPassword(u.saveConfigOnly, func() {
					u.config.SetShouldSaveFileEncrypted(false)
					u.saveConfigOnly()
				})
			} else {
				u.saveConfigOnly()
			}
		}
	}
}

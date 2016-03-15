package gui

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/config"
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
				u.captureInitialMasterPassword(u.saveConfigOnly)
			} else {
				u.saveConfigOnly()
			}
		}
	}
}

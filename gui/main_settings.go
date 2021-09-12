package gui

import "github.com/coyim/coyim/gui/settings"

type mainSettings struct {
	displaySettings  *displaySettings
	keyboardSettings *keyboardSettings
	settings         *settings.Settings
}

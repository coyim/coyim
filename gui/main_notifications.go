package gui

import "time"

type mainNotifications struct {
	deNotify    *desktopNotifications
	actionTimes map[string]time.Time
}

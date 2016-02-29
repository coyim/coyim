package gui

func initSignals() {
	accountChangedSignal, _ = g.glib.SignalNew("coyim-account-changed")
	enableWindow, _ = g.glib.SignalNew("enable")
	disableWindow, _ = g.glib.SignalNew("disable")
}

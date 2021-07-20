package gui

func initSignals(iui *inUIThread) {
	accountChangedSignal, _ = iui.g.glib.SignalNew("coyim-account-changed")
	enableWindow, _ = iui.g.glib.SignalNew("enable")
	disableWindow, _ = iui.g.glib.SignalNew("disable")
}

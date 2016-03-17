package gui

func assertInUIThread() {
	if g.glib.MainDepth() == 0 {
		panic("This function has to be called from the UI thread")
	}
}

func doInUIThread(f func()) {
	g.glib.IdleAdd(f)
}

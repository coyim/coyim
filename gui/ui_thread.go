package gui

//TODO: could this use a compiling flag to generate a noop function when released?
func assertInUIThread() {
	if g.glib.MainDepth() == 0 {
		panic("This function has to be called from the UI thread")
	}
}

//GTK process events in glib event loop (see [1]). In order to keep the UI
//responsive, it is a good practice to not block long running tasks in a signal's
//callback (you dont want a button to keep looking pressed for a couple of seconds).
//doInUIThread schedule the function to run in the next
//1 - https://developer.gnome.org/glib/unstable/glib-The-Main-Event-Loop.html
//TODO: Try other patterns and expose them as API. Example: http://www.mono-project.com/docs/gui/gtksharp/responsive-applications/
func doInUIThread(f func()) {
	_, _ = g.glib.IdleAdd(f)
}

package gui

import (
	"bytes"

	"github.com/coyim/gotk3adapter/gtki"
)

type xmlConsoleView struct {
	gtki.Dialog `gtk-widget:"XMLConsole"`

	buf gtki.TextBuffer `gtk-widget:"consoleContent"`
	log *bytes.Buffer
}

func newXMLConsoleView(log *bytes.Buffer) *xmlConsoleView {
	builder := newBuilder("XMLConsole")
	d := &xmlConsoleView{
		log: log,
	}

	panicOnDevError(builder.bindObjects(d))

	d.refresh()
	builder.ConnectSignals(map[string]interface{}{
		"on_refresh": d.refresh,
		"on_close":   d.Destroy,
	})

	return d
}

func (c *xmlConsoleView) refresh() {
	c.buf.Delete(c.buf.GetStartIter(), c.buf.GetEndIter())
	if c.log != nil {
		c.buf.Insert(c.buf.GetEndIter(), c.log.String())
	}
}

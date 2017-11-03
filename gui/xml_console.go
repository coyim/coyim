package gui

import (
	"bytes"

	"github.com/coyim/gotk3adapter/gtki"
)

type xmlConsoleView struct {
	gtki.Dialog

	buf gtki.TextBuffer
	log *bytes.Buffer
}

func newXMLConsoleView(log *bytes.Buffer) *xmlConsoleView {
	builder := newBuilder("XMLConsole")
	d := &xmlConsoleView{
		Dialog: builder.getObj("XMLConsole").(gtki.Dialog),
		buf:    builder.getObj("consoleContent").(gtki.TextBuffer),
		log:    log,
	}

	d.refresh()
	builder.ConnectSignals(map[string]interface{}{
		"on_refresh_signal": d.refresh,
		"on_close_signal":   d.Destroy,
	})

	return d
}

func (c *xmlConsoleView) refresh() {
	c.buf.Delete(c.buf.GetStartIter(), c.buf.GetEndIter())
	if c.log != nil {
		c.buf.Insert(c.buf.GetEndIter(), c.log.String())
	}
}

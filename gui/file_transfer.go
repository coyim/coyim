package gui

import (
	"log"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp/utils"
	"github.com/twstrike/gotk3adapter/gtki"
)

func (u *gtkUI) startAllListenersFor(ev events.FileTransfer) {
	go func() {
		err, ok := <-ev.ErrorOccurred
		if ok {
			log.Printf("File transfer of file %s failed with %v", ev.Name, err)
			close(ev.CancelTransfer)
		}
	}()

	go func() {
		_, ok := <-ev.TransferFinished
		if ok {
			log.Printf("File transfer of file %s finished with success", ev.Name)
			close(ev.CancelTransfer)
		}
	}()

	go func() {
		for upd := range ev.Update {
			log.Printf("File transfer of file %s: %d/%d done", ev.Name, upd, ev.Size)
		}
	}()
}

func (u *gtkUI) handleFileTransfer(ev events.FileTransfer) {
	dialogID := "FileTransferAskToReceive"
	builder := newBuilder(dialogID)
	dialogOb := builder.getObj(dialogID)

	d := dialogOb.(gtki.MessageDialog)
	d.SetDefaultResponse(gtki.RESPONSE_YES)
	d.SetTransientFor(u.window)

	message := i18n.Localf("%s wants to send you a file - do you want to receive it?", utils.RemoveResourceFromJid(ev.Peer))
	secondary := i18n.Localf("File name: %s", ev.Name)
	if ev.Description != "" {
		secondary = i18n.Localf("%s\nDescription: %s", secondary, ev.Description)
	}
	if ev.DateLastModified != "" {
		secondary = i18n.Localf("%s\nLast modified: %s", secondary, ev.DateLastModified)
	}
	if ev.Size != 0 {
		secondary = i18n.Localf("%s\nSize: %d bytes", secondary, ev.Size)
	}

	d.SetProperty("text", message)
	d.SetProperty("secondary_text", secondary)

	responseType := gtki.ResponseType(d.Run())
	result := responseType == gtki.RESPONSE_YES
	d.Destroy()
	fname := "bogus"
	if result {
		u.startAllListenersFor(ev)
		ev.Answer <- &fname
	} else {
		ev.Answer <- nil
	}

	// for more fancy use, we should allow people a choice in where to save it, but that's for later.
}

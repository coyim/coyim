package gui

import (
	"log"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp/utils"
	"github.com/twstrike/gotk3adapter/gtki"
)

// OK, so from the user interface, for now, we need a few things:
//  - A way of choosing where the file should be put
//  - A way of displaying errors that happened during the transfer
//  - A way for the user to cancel the transfer
//  - A way to notify the user when the transfer is done
//  - A way to update the user interface about progress
//  In general, hopefully these methods are completely independent of transport. Once we get to encrypted transfer we might want to highlight that (and say something about the file being transmitted in the clear otherwise)

// Actual user interface:
//   First - ask if you want the file
//   Second - choose where to put the file using standard file chooser/saver
//   Third - one status bar per file with percentages etc.
//       This will get a checkbox and a message when done
//       Or it will get an error message when failed
//       There will be a cancel button there, that will cancel the file receipt

func (u *gtkUI) startAllListenersFor(ev events.FileTransfer, cv conversationView) {
	go func() {
		err, ok := <-ev.ErrorOccurred
		if ok {
			cv.failFileTransfer()
			log.Printf("File transfer of file %s failed with %v", ev.Name, err)
			close(ev.CancelTransfer)
		}
	}()

	go func() {
		_, ok := <-ev.TransferFinished
		if ok {
			cv.successFileTransfer()
			log.Printf("File transfer of file %s finished with success", ev.Name)
			close(ev.CancelTransfer)
		}
	}()

	go func() {
		for upd := range ev.Update {
			cv.startFileTransfer(float64((upd*100)/ev.Size) / 100)
			log.Printf("File transfer of file %s: %d/%d done", ev.Name, upd, ev.Size)

			if cv.isFileTransferCanceled() {
				log.Printf("File transfer of file canceled")
				ev.CancelTransfer <- true
				return
			}
		}
	}()
}

func (u *gtkUI) handleFileTransfer(ev events.FileTransfer) {
	dialogID := "FileTransferAskToReceive"
	builder := newBuilder(dialogID)
	dialogOb := builder.getObj(dialogID)
	account := u.findAccountForSession(ev.Session)

	d := dialogOb.(gtki.MessageDialog)
	d.SetDefaultResponse(gtki.RESPONSE_YES)
	d.SetTransientFor(u.window)

	message := i18n.Localf("%s wants to send you a file: do you want to receive it?", utils.RemoveResourceFromJid(ev.Peer))
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

	var name string

	if result {
		fdialog, _ := g.gtk.FileChooserDialogNewWith2Buttons(
			i18n.Local("Choose where to save file"),
			u.window,
			gtki.FILE_CHOOSER_ACTION_SAVE,
			i18n.Local("_Cancel"),
			gtki.RESPONSE_CANCEL,
			i18n.Local("_Save"),
			gtki.RESPONSE_OK,
		)

		fdialog.SetCurrentName(ev.Name)

		if gtki.ResponseType(fdialog.Run()) == gtki.RESPONSE_OK {
			name = fdialog.GetFilename()
		}
		fdialog.Destroy()
	}

	if result && name != "" {
		cv, _ := u.roster.openConversationView(account, ev.Peer, true)

		cv.showFileTransferNotification()

		u.startAllListenersFor(ev, cv)
		ev.Answer <- &name
	} else {
		ev.Answer <- nil
	}
}

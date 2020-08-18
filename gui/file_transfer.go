package gui

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

// OK, so from the user interface, for now, we need a few things:
//  - A way of choosing where the file should be put
//  - A way of displaying errors that happened during the transfer
//  - A way for the user to cancel the transfer
//  - A way to notify the user when the transfer is done
//  - A way to update the user interface about progress
//  In general, hopefully these methods are completely independent of transport.
// Once we get to encrypted transfer we might want to highlight that (and say
// something about the file being transmitted in the clear otherwise)

// Actual user interface:
//   First - ask if you want the file
//   Second - choose where to put the file using standard file chooser/saver
//   Third - one status bar per file with percentages etc.
//       This will get a checkbox and a message when done
//       Or it will get an error message when failed
//       There will be a cancel button there, that will cancel the file receipt

func (u *gtkUI) startAllListenersForFile(ctl *data.FileTransferControl, cv conversationView, file *fileNotification, name, verbing string) {
	go ctl.WaitForError(func(err error) {
		file.fail()
		cv.Log().WithError(err).WithFields(log.Fields{
			"verbing": verbing,
			"file":    name,
		}).Warn("file transfer failed")
	})

	go ctl.WaitForFinish(func(finished bool) {
		if finished {
			file.succeed()
			cv.Log().WithFields(log.Fields{
				"verbing": verbing,
				"file":    name,
			}).Info("file transfer finished with success")
			return
		}
		file.decline()
		cv.Log().WithFields(log.Fields{
			"verbing": verbing,
			"file":    name,
		}).Info("file transfer declined by receiver")
	})

	go ctl.WaitForUpdate(func(upd, total int64) {
		file.progress = float64((upd*100)/total) / 100
		cv.updateFileTransfer(file)
		cv.Log().WithFields(log.Fields{
			"verbing": verbing,
			"file":    name,
			"update":  upd,
			"total":   total,
		}).Info("file transfer partially done")
	})
}

func (u *gtkUI) askIfContinueUnencrypted(a *account, peer jid.Any) bool {
	dialogID := "FileTransferContinueWithoutEncryption"
	builder := newBuilder(dialogID)
	dialogOb := builder.getObj(dialogID)

	d := dialogOb.(gtki.MessageDialog)
	d.SetDefaultResponse(gtki.RESPONSE_NO)
	d.SetTransientFor(u.window)

	message := i18n.Localf("The transfer to %s can't be done encrypted and securely.", peer.NoResource())
	secondary := i18n.Localf("Do you want to continue anyway? This means an adversary or server administrator can potentially see the content of the file.")

	_ = d.SetProperty("text", message)
	_ = d.SetProperty("secondary_text", secondary)

	responseType := gtki.ResponseType(d.Run())
	result := responseType == gtki.RESPONSE_YES
	d.Destroy()

	return result
}

func (u *gtkUI) handleFileTransfer(ev events.FileTransfer, a *account) {
	dialogID := "FileTransferAskToReceive"
	builder := newBuilder(dialogID)
	dialogOb := builder.getObj(dialogID)

	d := dialogOb.(gtki.MessageDialog)
	d.SetDefaultResponse(gtki.RESPONSE_YES)
	d.SetTransientFor(u.window)

	var message, secondary string

	enc := ev.Encrypted

	if ev.IsDirectory {
		if enc {
			message = i18n.Localf("%s wants to send you a directory - this transfer will be encrypted and secure - do you want to receive it?", ev.Peer.NoResource())
			secondary = i18n.Localf("Directory name: %s", ev.Name)
		} else {
			message = i18n.Localf("%s wants to send you a directory - the transfer will NOT be encrypted or secure - do you want to receive it?", ev.Peer.NoResource())
			secondary = i18n.Localf("Directory name: %s", ev.Name)
		}
	} else {
		if enc {
			message = i18n.Localf("%s wants to send you a file - this transfer will be encrypted and secure - do you want to receive it?", ev.Peer.NoResource())
			secondary = i18n.Localf("File name: %s", ev.Name)
		} else {
			message = i18n.Localf("%s wants to send you a file - the transfer will NOT be encrypted or secure - do you want to receive it?", ev.Peer.NoResource())
			secondary = i18n.Localf("File name: %s", ev.Name)
		}
	}

	secondary = i18n.Localf("%s\nEncrypted: %v", secondary, enc)

	if ev.Description != "" {
		secondary = i18n.Localf("%s\nDescription: %s", secondary, ev.Description)
	}
	if ev.DateLastModified != "" {
		secondary = i18n.Localf("%s\nLast modified: %s", secondary, ev.DateLastModified)
	}
	if ev.Size != 0 {
		secondary = i18n.Localf("%s\nSize: %d bytes", secondary, ev.Size)
	}

	_ = d.SetProperty("text", message)
	_ = d.SetProperty("secondary_text", secondary)

	responseType := gtki.ResponseType(d.Run())
	result := responseType == gtki.RESPONSE_YES
	d.Destroy()

	var name string

	if result {
		label := "Choose where to save file"
		action := gtki.FILE_CHOOSER_ACTION_SAVE
		if ev.IsDirectory {
			label = "Choose where to save folder"
			action = gtki.FILE_CHOOSER_ACTION_CREATE_FOLDER
		}

		fdialog, _ := g.gtk.FileChooserDialogNewWith2Buttons(
			i18n.Local(label),
			u.window,
			action,
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
		fileName := resizeFileName(ev.Name)
		cv := u.openConversationView(a, ev.Peer, true)
		f := createNewFileTransferWithDefaults(fileName, ev.IsDirectory, false, true, ev.Control, cv)
		f.setEncryptionInformation(enc, cv.hasVerifiedKey())
		f.updateLabel()
		u.startAllListenersForFile(ev.Control, cv, f, ev.Name, "Receiving")
		ev.Answer <- &name
	} else {
		ev.Answer <- nil
	}
}

func createNewFileTransferWithDefaults(fileName string, dir bool, sending, receiving bool, ctl *data.FileTransferControl, cv conversationView) *fileNotification {
	f := cv.newFileTransfer(fileName, dir, sending, receiving)
	f.afterCancel(func() {
		cv.updateFileTransferNotificationCounts()
		ctl.Cancel()
	})
	f.afterDeclined(cv.updateFileTransferNotificationCounts)
	f.afterFail(cv.updateFileTransferNotificationCounts)
	f.afterSucceed(cv.updateFileTransferNotificationCounts)
	return f
}

func (account *account) sendThingTo(peer jid.Any, u *gtkUI, name string, dir bool, ctl *data.FileTransferControl, verifiedKey bool, encDecision chan bool) {
	nm := resizeFileName(filepath.Base(name))
	cv := u.openConversationView(account, peer, true)
	f := createNewFileTransferWithDefaults(nm, dir, true, false, ctl, cv)
	go func() {
		encDec := <-encDecision
		f.setEncryptionInformation(encDec, verifiedKey)
		f.updateLabel()
	}()
	u.startAllListenersForFile(ctl, cv, f, nm, "Sending")
}

func (u *gtkUI) ensureJidHasResource(account *account, p jid.Any) jid.Any {
	if v, ok := p.(jid.WithResource); ok {
		return v
	}
	if peer, ok := u.getPeer(account, p.NoResource()); ok {
		res := peer.ResourceToUseFallback()
		account.Log().WithFields(log.Fields{
			"peer":     p,
			"resource": res,
		}).Debug("adding resource to peer")
		if res != jid.NewResource("") {
			return p.WithResource(res)
		}
	}

	account.Log().WithField("peer", p).Error("couldn't find a valid resource for peer")

	return p
}

func (account *account) sendFileTo(peer jid.Any, u *gtkUI, cp *conversationPane) {
	peer = u.ensureJidHasResource(account, peer)
	if file, ok := chooseFileToSend(u.window); ok {
		encDecision := make(chan bool, 1)
		ctl := account.session.SendFileTo(peer, file, func() bool {
			result := make(chan bool)
			doInUIThread(func() {
				result <- u.askIfContinueUnencrypted(account, peer)
			})
			return <-result
		}, func(enc bool) {
			encDecision <- enc
		})
		account.sendThingTo(peer, u, file, false, ctl, cp.hasVerifiedKey(), encDecision)
	}
}

func (account *account) sendDirectoryTo(peer jid.Any, u *gtkUI, cp *conversationPane) {
	peer = u.ensureJidHasResource(account, peer)
	if dir, ok := chooseDirToSend(u.window); ok {
		encDecision := make(chan bool, 1)
		ctl := account.session.SendDirTo(peer, dir, func() bool {
			result := make(chan bool)
			doInUIThread(func() {
				result <- u.askIfContinueUnencrypted(account, peer)
			})
			return <-result
		}, func(enc bool) {
			encDecision <- enc
		})
		account.sendThingTo(peer, u, dir, true, ctl, cp.hasVerifiedKey(), encDecision)
	}
}

func chooseItemToSend(w gtki.Window, action gtki.FileChooserAction, title string) (string, bool) {
	dialog, _ := g.gtk.FileChooserDialogNewWith2Buttons(
		i18n.Local(title),
		w,
		action,
		i18n.Local("_Cancel"),
		gtki.RESPONSE_CANCEL,
		i18n.Local("Send"),
		gtki.RESPONSE_OK,
	)
	defer dialog.Destroy()

	if gtki.ResponseType(dialog.Run()) == gtki.RESPONSE_OK {
		return dialog.GetFilename(), true
	}
	return "", false
}

func chooseFileToSend(w gtki.Window) (string, bool) {
	return chooseItemToSend(w, gtki.FILE_CHOOSER_ACTION_OPEN, "Choose file to send")
}

func chooseDirToSend(w gtki.Window) (string, bool) {
	return chooseItemToSend(w, gtki.FILE_CHOOSER_ACTION_SELECT_FOLDER, "Choose directory to send")
}

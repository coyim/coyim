package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

// In these notifications we will use the convention that:
// - "secure transfer" means that you are sending or receiving something encrypted from/to a peer that is verified
// - "encrypted transfer" means transfer to/from a peer that is not verified
// - "insecure transfer" is unencrypted

type fileNotification struct {
	area                      gtki.Box
	label                     gtki.Label
	image                     gtki.Image
	name                      string
	progress                  float64
	state                     string
	directory                 bool
	sending                   bool
	receiving                 bool
	encrypted                 bool
	verifiedPeer              bool
	haveEncryptionInformation bool
	afterCancelHook           func()
	afterFailHook             func()
	afterSucceedHook          func()
	afterDeclinedHook         func()
}

type fileTransferNotification struct {
	area          gtki.Box
	image         gtki.Image
	label         gtki.Label
	box           gtki.Box
	progressBar   gtki.ProgressBar
	button        gtki.Button
	labelButton   gtki.Label
	totalProgress float64
	files         []*fileNotification
	count         int
	canceled      bool
}

func resizeFileName(name string) string {
	var fileName string

	if len(name) > 20 {
		fileName = name[:20] + "..."
		return fileName
	}

	return name
}

func (file *fileNotification) afterCancel(f func()) {
	file.afterCancelHook = f
}

func (file *fileNotification) afterDeclined(f func()) {
	file.afterDeclinedHook = f
}

func (file *fileNotification) afterFail(f func()) {
	file.afterFailHook = f
}

func (file *fileNotification) afterSucceed(f func()) {
	file.afterSucceedHook = f
}

func (file *fileNotification) destroy() {
	file.cancel()
}

func (file *fileNotification) update(fileName string, prov gtki.CssProvider) {
	updateWithStyle(file.label, prov)
	file.label.SetLabel(fileName)
	file.image.Hide()
}

func (b *builder) fileTransferNotifInit() *fileTransferNotification {
	fileTransferNotif := &fileTransferNotification{}

	b.getItems(
		"file-transfer", &fileTransferNotif.area,
		"image-file-transfer", &fileTransferNotif.image,
		"label-file-transfer", &fileTransferNotif.label,
		"info-file-transfer", &fileTransferNotif.box,
		"bar-file-transfer", &fileTransferNotif.progressBar,
		"button-file-transfer", &fileTransferNotif.button,
		"button-label-file-transfer", &fileTransferNotif.labelButton,
	)

	return fileTransferNotif
}

func (f *fileNotification) setEncryptionInformation(encrypted, verifiedPeer bool) {
	f.encrypted = encrypted
	f.verifiedPeer = verifiedPeer
	f.haveEncryptionInformation = true
}

func (conv *conversationPane) newFileTransfer(fileName string, dir, send, receive bool) *fileNotification {
	if !conv.fileTransferNotif.area.IsVisible() {
		prov := providerWithCSS("box { background-color: #fff9f3;  color: #000000; border: 3px; }")
		updateWithStyle(conv.fileTransferNotif.area, prov)

		conv.fileTransferNotif.progressBar.SetFraction(0.0)
		conv.fileTransferNotif.canceled = false
	}

	info := conv.createFileTransferNotification(fileName, dir, send, receive)
	info.updateLabel()

	conv.fileTransferNotif.area.SetVisible(true)

	countSending := 0
	countReceiving := 0

	label := "Transfer started"

	for _, f := range conv.fileTransferNotif.files {
		if f.sending {
			countSending++
		}
		if f.receiving {
			countReceiving++
		}
	}

	cc := i18n.Local("Cancel")

	if countSending > 0 && countReceiving == 0 {
		doInUIThread(func() {
			conv.updateFileTransferNotification(label, cc, "filetransfer_send.svg")
		})
	} else if countSending == 0 && countReceiving > 0 {
		doInUIThread(func() {
			conv.updateFileTransferNotification(label, cc, "filetransfer_receive.svg")
		})
	} else if countSending > 0 && countReceiving > 0 {
		doInUIThread(func() {
			conv.updateFileTransferNotification(label, cc, "filetransfer_receive_send.svg")
		})
	}

	return info
}

func (f *fileNotification) updateLabel() {
	var label string
	switch {
	case f.sending && !f.haveEncryptionInformation:
		label = i18n.Localf("Sending: %s", f.name)
	case !f.sending && !f.haveEncryptionInformation:
		label = i18n.Localf("Receiving: %s", f.name)
	case f.sending && f.encrypted && f.verifiedPeer:
		label = i18n.Localf("Sending securely: %s", f.name)
	case f.sending && f.encrypted:
		label = i18n.Localf("Sending encrypted: %s", f.name)
	case f.sending:
		label = i18n.Localf("Sending insecurely: %s", f.name)
	case f.encrypted && f.verifiedPeer:
		label = i18n.Localf("Receiving securely: %s", f.name)
	case f.encrypted:
		label = i18n.Localf("Receiving encrypted: %s", f.name)
	default:
		label = i18n.Localf("Receiving insecurely: %s", f.name)
	}

	doInUIThread(func() {
		f.label.SetLabel(label)
	})
}

func (conv *conversationPane) createFileTransferNotification(fileName string, dir, send, receive bool) *fileNotification {
	b := newBuilder("FileTransferNotification")

	file := &fileNotification{directory: dir, sending: send, receiving: receive, state: stateInProgress}

	b.getItems(
		"area-file-transfer-info", &file.area,
		"name-file-transfer-info", &file.label,
		"image-file-transfer-info", &file.image,
	)

	b.ConnectSignals(map[string]interface{}{
		"on_destroy_single_file_transfer": file.destroy,
	})

	file.name = fileName

	file.updateLabel()

	conv.fileTransferNotif.count++
	conv.fileTransferNotif.canceled = false
	conv.fileTransferNotif.totalProgress = 0

	conv.fileTransferNotif.box.Add(file.area)
	file.area.ShowAll()

	conv.fileTransferNotif.files = append(conv.fileTransferNotif.files, file)

	return file
}

func (conv *conversationPane) updateFileTransferNotification(label, buttonLabel, image string) {
	if buttonLabel == i18n.Local("Close") {
		prov := providerWithCSS("label { margin-right: 3px;  margin-left: 3px; }")
		updateWithStyle(conv.fileTransferNotif.labelButton, prov)
	}
	conv.account.log.WithField("label", label).Info("Updating notification")

	conv.fileTransferNotif.label.SetLabel(label)
	conv.fileTransferNotif.labelButton.SetLabel(buttonLabel)
	setImageFromFile(conv.fileTransferNotif.image, image)
}

const stateInProgress = "in-progress"
const stateSuccess = "success"
const stateFailed = "failed"
const stateCanceled = "canceled"
const stateDeclined = "declined"

func (conv *conversationPane) updateFileTransfer(file *fileNotification) {
	conv.fileTransferNotif.totalProgress = 0
	count := 0
	haveSuccess := false
	for _, f := range conv.fileTransferNotif.files {
		switch f.state {
		case stateInProgress:
			count++
			conv.fileTransferNotif.totalProgress += f.progress
		case stateSuccess:
			haveSuccess = true
		}
	}

	var upd float64
	if count == 0 {
		if haveSuccess {
			upd = 100
		} else {
			upd = conv.fileTransferNotif.totalProgress
		}
	} else {
		upd = conv.fileTransferNotif.totalProgress / float64(count)
	}

	doInUIThread(func() {
		conv.fileTransferNotif.progressBar.SetFraction(upd)
	})
}

func fileTransferCalculateStates(countCompleted, countCanceled, countFailed, countDeclined, countDirs, countDirsCompleted, countTotal int, canceledBefore bool) (label, image string, canceled bool) {
	verb := "successful"
	image = "success.svg"
	canceled = canceledBefore
	if countCanceled+countFailed+countDeclined == countTotal {
		image = "failure.svg"
		canceled = true
		verb = "failed"
		if countCanceled > (countFailed + countDeclined) {
			verb = "canceled"
		} else if countDeclined > (countCanceled + countFailed) {
			verb = "declined"
		}
	}

	switch {
	case countDirsCompleted > 0 && countCompleted != countDirsCompleted:
		label = i18n.Local("File and directory transfer(s) " + verb)
	case countDirsCompleted > 0:
		label = i18n.Local("Directory transfer(s) " + verb)
	default:
		label = i18n.Local("File transfer(s) " + verb)
	}
	return
}

func (conv *conversationPane) updateFileTransferNotificationCounts() {
	countCompleted := 0
	countCanceled := 0
	countFailed := 0
	countTotal := 0
	countDirs := 0
	countDirsCompleted := 0
	countDeclined := 0
	for _, f := range conv.fileTransferNotif.files {
		switch f.state {
		case stateInProgress:
		case stateSuccess:
			countCompleted++
			if f.directory {
				countDirsCompleted++
			}
		case stateCanceled:
			countCompleted++
			countCanceled++
		case stateFailed:
			countCompleted++
			countFailed++
		case stateDeclined:
			countCompleted++
			countDeclined++
		}
		if f.directory {
			countDirs++
		}
		countTotal++
	}

	conv.fileTransferNotif.count = countTotal - countCompleted
	if countCompleted == countTotal {
		label, image, c := fileTransferCalculateStates(countCompleted, countCanceled, countFailed, countDeclined, countDirs, countDirsCompleted, countTotal, conv.fileTransferNotif.canceled)
		conv.fileTransferNotif.canceled = c
		doInUIThread(func() {
			conv.updateFileTransferNotification(label, i18n.Local("Close"), image)
		})
	}
}

func (conv *conversationPane) isFileTransferNotifCanceled() bool {
	return conv.fileTransferNotif.canceled
}

func canceledProvider() gtki.CssProvider {
	return providerWithCSS("label { color: #cc3636; }")
}

func successProvider() gtki.CssProvider {
	return providerWithCSS("label { color: #89AF8F; }")
}

func (file *fileNotification) decline() {
	if file.state != stateInProgress {
		return
	}
	file.state = stateDeclined
	file.progress = 0
	file.update(i18n.Localf("Declined: %s", file.name), canceledProvider())
	hook := file.afterDeclinedHook
	file.afterDeclinedHook = nil
	if hook != nil {
		hook()
	}
}

func (file *fileNotification) cancel() {
	if file.state != stateInProgress {
		return
	}
	file.state = stateCanceled
	file.progress = 0
	file.update(i18n.Localf("Canceled: %s", file.name), canceledProvider())
	hook := file.afterCancelHook
	file.afterCancelHook = nil
	if hook != nil {
		hook()
	}
}

func (file *fileNotification) fail() {
	if file.state != stateInProgress {
		return
	}
	file.state = stateFailed
	file.progress = 0
	file.update(i18n.Localf("Failed: %s", file.name), canceledProvider())
	hook := file.afterFailHook
	file.afterFailHook = nil
	if hook != nil {
		hook()
	}
}

func (file *fileNotification) succeed() {
	if file.state != stateInProgress {
		return
	}
	file.state = stateSuccess

	var label string

	switch {
	case file.sending && file.encrypted && file.verifiedPeer:
		label = i18n.Localf("Sent securely: %s", file.name)
	case file.sending && file.encrypted:
		label = i18n.Localf("Sent encrypted: %s", file.name)
	case file.sending:
		label = i18n.Localf("Sent insecurely: %s", file.name)
	case file.encrypted && file.verifiedPeer:
		label = i18n.Localf("Received securely: %s", file.name)
	case file.encrypted:
		label = i18n.Localf("Received encrypted: %s", file.name)
	default:
		label = i18n.Localf("Received insecurely: %s", file.name)
	}

	file.update(label, successProvider())

	hook := file.afterSucceedHook
	file.afterSucceedHook = nil
	if hook != nil {
		hook()
	}
}

func (conv *conversationPane) onDestroyFileTransferNotif() {
	label := conv.fileTransferNotif.labelButton.GetLabel()
	if label == i18n.Local("Cancel") {
		for _, f := range conv.fileTransferNotif.files {
			f.cancel()
		}
	} else {
		conv.fileTransferNotif.canceled = false
		conv.fileTransferNotif.area.SetVisible(false)
		conv.fileTransferNotif.progressBar.SetFraction(0.0)
		for i := range conv.fileTransferNotif.files {
			conv.fileTransferNotif.files[i].area.Destroy()
		}
		conv.fileTransferNotif.files = conv.fileTransferNotif.files[:0]
		conv.fileTransferNotif.count = 0
	}
}

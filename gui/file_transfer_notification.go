package gui

import (
	"log"

	"github.com/coyim/gotk3adapter/gtki"
)

type fileNotification struct {
	area             gtki.Box
	label            gtki.Label
	image            gtki.Image
	name             string
	progress         float64
	state            string
	directory        bool
	afterCancelHook  func()
	afterFailHook    func()
	afterSucceedHook func()
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

// TODO: change ui when sending and receiving. No
// more file infobars, but something else

func resizeFileName(name string) string {
	var fileName string

	if len(name) > 20 {
		fileName = name[:21] + "..."
		return fileName
	}

	return name
}

func (file *fileNotification) afterCancel(f func()) {
	file.afterCancelHook = f
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

type StyleContextable interface {
	GetStyleContext() (gtki.StyleContext, error)
}

func (conv *conversationPane) newFileTransfer(fileName, purpose string, dir bool) *fileNotification {
	if !conv.fileTransferNotif.area.IsVisible() {
		prov := providerWithCss("box { background-color: #fff9f3;  color: #000000; border: 3px; }")
		updateWithStyle(conv.fileTransferNotif.area, prov)

		conv.fileTransferNotif.progressBar.SetFraction(0.0)
		conv.fileTransferNotif.canceled = false
	}

	info := conv.createFileTransferNotification(fileName, purpose, dir)
	conv.fileTransferNotif.area.SetVisible(true)
	return info
}

func (conv *conversationPane) createFileTransferNotification(fileName, purpose string, dir bool) *fileNotification {
	b := newBuilder("FileTransferNotification")

	file := &fileNotification{directory: dir, state: stateInProgress}

	b.getItems(
		"area-file-transfer-info", &file.area,
		"name-file-transfer-info", &file.label,
		"image-file-transfer-info", &file.image,
	)

	b.ConnectSignals(map[string]interface{}{
		"on_destroy_single_file_transfer": file.destroy,
	})

	label := "File transfer started"
	if dir {
		label = "Directory transfer started"
	}

	file.name = fileName
	if purpose == "send" {
		fileName = "Sending: " + fileName
		conv.updateFileTransferNotification(label, "Cancel", "filetransfer_send.svg")
	} else {
		fileName = "Receiving: " + fileName
		conv.updateFileTransferNotification(label, "Cancel", "filetransfer_receive.svg")
	}

	file.label.SetLabel(fileName)
	conv.fileTransferNotif.count++
	conv.fileTransferNotif.canceled = false
	conv.fileTransferNotif.totalProgress = 0

	conv.fileTransferNotif.box.Add(file.area)
	file.area.ShowAll()

	conv.fileTransferNotif.files = append(conv.fileTransferNotif.files, file)

	return file
}

func (conv *conversationPane) updateFileTransferNotification(label, buttonLabel, image string) {
	if buttonLabel == "Close" {
		prov := providerWithCss("label { margin-right: 3px;  margin-left: 3px; }")
		updateWithStyle(conv.fileTransferNotif.labelButton, prov)
	}
	log.Printf(label)

	conv.fileTransferNotif.label.SetLabel(label)
	conv.fileTransferNotif.labelButton.SetLabel(buttonLabel)
	setImageFromFile(conv.fileTransferNotif.image, image)
}

const stateInProgress = "in-progress"
const stateSuccess = "success"
const stateFailed = "failed"
const stateCanceled = "canceled"

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

func fileTransferCalculateStates(countCompleted, countCanceled, countFailed, countDirs, countTotal int, canceledBefore bool) (label, image string, canceled bool) {
	verb := "successful"
	image = "success.svg"
	canceled = canceledBefore
	if countCanceled+countFailed == countTotal {
		image = "failure.svg"
		canceled = true
		verb = "failed"
		if countCanceled > countFailed {
			verb = "canceled"
		}
	}

	label = "File transfer " + verb
	if countDirs > 0 {
		label = "Directory transfer " + verb
	}
	return
}

func (conv *conversationPane) updateFileTransferNotificationCounts() {
	countCompleted := 0
	countCanceled := 0
	countFailed := 0
	countTotal := 0
	countDirs := 0
	for _, f := range conv.fileTransferNotif.files {
		switch f.state {
		case stateInProgress:
		case stateSuccess:
			countCompleted++
		case stateCanceled:
			countCompleted++
			countCanceled++
		case stateFailed:
			countCompleted++
			countFailed++
		}
		if f.directory {
			countDirs++
		}
		countTotal++
	}
	conv.fileTransferNotif.count = countTotal - countCompleted
	if countCompleted == countTotal {
		label, image, c := fileTransferCalculateStates(countCompleted, countCanceled, countFailed, countDirs, countTotal, conv.fileTransferNotif.canceled)
		conv.fileTransferNotif.canceled = c
		doInUIThread(func() {
			conv.updateFileTransferNotification(label, "Close", image)
		})
	}
}

func (conv *conversationPane) isFileTransferNotifCanceled() bool {
	return conv.fileTransferNotif.canceled
}

func canceledProvider() gtki.CssProvider {
	return providerWithCss("label { color: #cc3636; }")
}

func successProvider() gtki.CssProvider {
	return providerWithCss("label { color: #89AF8F; }")
}

func providerWithCss(s string) gtki.CssProvider {
	p, _ := g.gtk.CssProviderNew()
	_ = p.LoadFromData(s)
	return p
}

func updateWithStyle(l StyleContextable, p gtki.CssProvider) {
	sc, _ := l.GetStyleContext()
	sc.AddProvider(p, 9999)
}

func (file *fileNotification) cancel() {
	if file.state != stateInProgress {
		return
	}
	file.state = stateCanceled
	file.progress = 0
	file.update("Canceled: "+file.name, canceledProvider())
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
	file.update("Failed: "+file.name, canceledProvider())
	hook := file.afterFailHook
	file.afterFailHook = nil
	if hook != nil {
		hook()
	}
}

func (file *fileNotification) succeed(purpose string) {
	if file.state != stateInProgress {
		return
	}
	file.state = stateSuccess
	if purpose == "send" {
		file.update("Sent: "+file.name, successProvider())
	} else {
		file.update("Received: "+file.name, successProvider())
	}
	hook := file.afterSucceedHook
	file.afterSucceedHook = nil
	if hook != nil {
		hook()
	}
}

func (conv *conversationPane) onDestroyFileTransferNotif() {
	label := conv.fileTransferNotif.labelButton.GetLabel()
	if label == "Cancel" {
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

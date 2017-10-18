package gui

import (
	"fmt"
	"log"

	"github.com/coyim/gotk3adapter/gtki"
)

type fileNotification struct {
	area      gtki.Box
	label     gtki.Label
	image     gtki.Image
	name      string
	progress  float64
	success   bool
	failed    bool
	canceled  bool
	completed bool
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
	count         int
	canceled      bool
	files         []*fileNotification
}

// TODO: there is still some issues around more than
// two transfers on cancel and failed
func any(vs []*fileNotification, f func(*fileNotification) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func all(vs []*fileNotification, f func(*fileNotification) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func resizeFileName(name string) string {
	var fileName string

	if len(name) > 20 {
		fileName = name[:21] + "..."
		return fileName
	}

	return name
}

func (file *fileNotification) destroy() {
	file.canceled = true
}

func (file *fileNotification) update(fileName string) {
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

func (conv *conversationPane) showFileTransferInfo(fileName string) *fileNotification {
	b := newBuilder("FileTransferNotification")

	file := &fileNotification{}

	b.getItems(
		"area-file-transfer-info", &file.area,
		"name-file-transfer-info", &file.label,
		"image-file-transfer-info", &file.image,
	)

	b.ConnectSignals(map[string]interface{}{
		"on_destroy_single_file_transfer": file.destroy,
	})

	label := "File transfer started"
	conv.updateFileTransferNotification(label, "Cancel", "filetransfer.svg")

	file.name = fileName
	fileName = "Receiving: " + fileName
	file.label.SetLabel(fileName)
	conv.fileTransferNotif.count++
	conv.fileTransferNotif.canceled = false

	conv.fileTransferNotif.box.Add(file.area)
	file.area.ShowAll()

	conv.fileTransferNotif.files = append(conv.fileTransferNotif.files, file)

	return file
}

func (conv *conversationPane) showFileTransferNotification(fileName string) *fileNotification {
	prov, _ := g.gtk.CssProviderNew()

	css := fmt.Sprintf(`
	box { background-color: #fff9f3;
	      color: #000000;
	      border: 3px;
	     }
	`)
	_ = prov.LoadFromData(css)

	styleContext, _ := conv.fileTransferNotif.area.GetStyleContext()
	styleContext.AddProvider(prov, 9999)

	label := "File transfer started"
	conv.updateFileTransferNotification(label, "Cancel", "filetransfer.svg")
	conv.fileTransferNotif.progressBar.SetFraction(0.0)
	conv.fileTransferNotif.canceled = false

	info := conv.showFileTransferInfo(fileName)

	conv.fileTransferNotif.area.SetVisible(true)

	return info
}

func (conv *conversationPane) getFileTransferNotification() bool {
	if !conv.fileTransferNotif.area.IsVisible() {
		return false
	}
	return true
}

func (conv *conversationPane) updateFileTransferNotification(label, buttonLabel, image string) {
	if buttonLabel == "Close" {
		prov, _ := g.gtk.CssProviderNew()

		css := fmt.Sprintf(`
	                  label { margin-right: 3px;
	                          margin-left: 3px;
	                  }
	                `)
		_ = prov.LoadFromData(css)

		styleContext, _ := conv.fileTransferNotif.labelButton.GetStyleContext()
		styleContext.AddProvider(prov, 9999)
	}
	log.Printf(label)

	conv.fileTransferNotif.label.SetLabel(label)
	conv.fileTransferNotif.labelButton.SetLabel(buttonLabel)
	setImageFromFile(conv.fileTransferNotif.image, image)
}

func (conv *conversationPane) startFileTransfer(file *fileNotification) {
	conv.fileTransferNotif.totalProgress = 0.0
	for i := range conv.fileTransferNotif.files {
		fmt.Println("THE PROGRESS %d", conv.fileTransferNotif.files[i].progress)
		conv.fileTransferNotif.totalProgress += conv.fileTransferNotif.files[i].progress
	}

	upd := conv.fileTransferNotif.totalProgress / float64(conv.fileTransferNotif.count)
	conv.fileTransferNotif.progressBar.SetFraction(upd)
}

func (conv *conversationPane) successFileTransfer(file *fileNotification) {
	prov, _ := g.gtk.CssProviderNew()

	css := fmt.Sprintf(`
	label { color: #89AF8F;
	      }
	`)
	_ = prov.LoadFromData(css)

	styleContext, _ := file.label.GetStyleContext()
	styleContext.AddProvider(prov, 9999)

	fileName := "Received: " + file.name
	file.update(fileName)
	file.success = true
	file.completed = true

	if all(conv.fileTransferNotif.files, func(f *fileNotification) bool {
		return f.completed
	}) {
		if any(conv.fileTransferNotif.files, func(f *fileNotification) bool {
			return f.success
		}) {
			label := "File transfer successful"
			conv.updateFileTransferNotification(label, "Close", "success.svg")
		}
	}
}

func (conv *conversationPane) failFileTransfer(file *fileNotification) {
	prov, _ := g.gtk.CssProviderNew()

	css := fmt.Sprintf(`
	label { color: #cc3636;
	     }
	`)
	_ = prov.LoadFromData(css)

	styleContext, _ := file.label.GetStyleContext()
	styleContext.AddProvider(prov, 9999)

	fileName := "Failed: " + file.name
	file.update(fileName)
	file.failed = true
	file.completed = true
	file.progress = 0.0
	conv.fileTransferNotif.count--

	if all(conv.fileTransferNotif.files, func(f *fileNotification) bool {
		return f.completed
	}) {
		if all(conv.fileTransferNotif.files, func(f *fileNotification) bool {
			return f.canceled || f.failed
		}) {
			label := "File transfer failed"
			conv.updateFileTransferNotification(label, "Close", "failure.svg")
		}
	}
}

func (conv *conversationPane) cancelFileTransfer(file *fileNotification) {
	prov, _ := g.gtk.CssProviderNew()

	css := fmt.Sprintf(`
	label { color: #cc3636;
	      }
	`)
	_ = prov.LoadFromData(css)

	styleContext, _ := file.label.GetStyleContext()
	styleContext.AddProvider(prov, 9999)

	fileName := "Canceled: " + file.name
	file.update(fileName)
	file.canceled = true
	file.completed = true
	file.progress = 0.0
	conv.fileTransferNotif.count--

	if all(conv.fileTransferNotif.files, func(f *fileNotification) bool {
		return f.completed
	}) {
		if all(conv.fileTransferNotif.files, func(f *fileNotification) bool {
			return f.canceled || f.failed
		}) {
			conv.fileTransferNotif.canceled = true
			label := "File transfer canceled"
			conv.updateFileTransferNotification(label, "Close", "failure.svg")
		}
	}
}

func (conv *conversationPane) isFileTransferNotifCanceled() bool {
	return conv.fileTransferNotif.canceled
}

func (conv *conversationPane) onDestroyFileTransferNotif() {
	label := conv.fileTransferNotif.labelButton.GetLabel()
	if label == "Cancel" {
		conv.fileTransferNotif.canceled = true

		prov, _ := g.gtk.CssProviderNew()

		css := fmt.Sprintf(`
	        label { color: #cc3636;
	             }
	        `)
		_ = prov.LoadFromData(css)

		label := "File transfer canceled"
		conv.updateFileTransferNotification(label, "Close", "failure.svg")

		files := conv.fileTransferNotif.files
		for i, f := range files {
			if f.success {
				break
			}
			styleContext, _ := files[i].label.GetStyleContext()
			styleContext.AddProvider(prov, 9999)
			files[i].update("Canceled: " + f.name)
		}
	} else {
		conv.fileTransferNotif.canceled = false
		conv.fileTransferNotif.area.SetVisible(false)
		conv.fileTransferNotif.progressBar.SetFraction(0.0)
		for i := range conv.fileTransferNotif.files {
			conv.fileTransferNotif.files[i].area.Destroy()
		}
		conv.fileTransferNotif.files = conv.fileTransferNotif.files[:0]
	}
}

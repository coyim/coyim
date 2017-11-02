package data

type transferUpdate struct {
	current, total int64
}

// FileTransferControl supplies the capabilities to control the file transfer
type FileTransferControl struct {
	cancelTransfer   chan bool           // one time use
	errorOccurred    chan error          // one time use
	update           chan transferUpdate // will be called many times
	transferFinished chan bool           // one time use
}

func newFileTransferControl(c chan bool, e chan error, u chan transferUpdate, t chan bool) *FileTransferControl {
	return &FileTransferControl{cancelTransfer: c, errorOccurred: e, update: u, transferFinished: t}
}

func CreateFileTransferControl() *FileTransferControl {
	return newFileTransferControl(make(chan bool), make(chan error), make(chan transferUpdate, 1000), make(chan bool))
}

func (ctl *FileTransferControl) WaitForFinish(k func()) {
	_, ok := <-ctl.transferFinished
	if ok {
		k()
		close(ctl.cancelTransfer)
	}
}

func (ctl *FileTransferControl) WaitForError(k func(error)) {
	e, ok := <-ctl.errorOccurred
	if ok {
		k(e)
		close(ctl.cancelTransfer)
	}
}

func (ctl *FileTransferControl) WaitForCancel(k func()) {
	if cancel, ok := <-ctl.cancelTransfer; ok && cancel {
		k()
		ctl.CloseAll()
	}
}

func (ctl *FileTransferControl) WaitForUpdate(k func(int64, int64)) {
	for upd := range ctl.update {
		k(upd.current, upd.total)
	}
}

func (ctl *FileTransferControl) Cancel() {
	ctl.cancelTransfer <- true
}

func (ctl *FileTransferControl) ReportError(e error) {
	ctl.closeTransferFinished()
	ctl.closeUpdate()
	ctl.sendAndCloseErrorOccurred(e)
}

func (ctl *FileTransferControl) ReportErrorNonblocking(e error) {
	go ctl.ReportError(e)
}

func (ctl *FileTransferControl) ReportFinished() {
	ctl.closeErrorOccurred()
	ctl.closeUpdate()
	ctl.sendAndCloseTransferFinished(true)
}

func (ctl *FileTransferControl) SendUpdate(current, total int64) {
	one := ctl.update
	if one != nil {
		one <- transferUpdate{current, total}
	}
}

func (ctl *FileTransferControl) CloseAll() {
	ctl.closeTransferFinished()
	ctl.closeUpdate()
	ctl.closeErrorOccurred()
}

func (ctl *FileTransferControl) closeTransferFinished() {
	c := ctl.transferFinished
	ctl.transferFinished = nil
	if c != nil {
		close(c)
	}
}

func (ctl *FileTransferControl) closeUpdate() {
	c := ctl.update
	ctl.update = nil
	if c != nil {
		close(c)
	}
}

func (ctl *FileTransferControl) closeErrorOccurred() {
	c := ctl.errorOccurred
	ctl.errorOccurred = nil
	if c != nil {
		close(c)
	}
}

func (ctl *FileTransferControl) sendAndCloseErrorOccurred(e error) {
	c := ctl.errorOccurred
	ctl.errorOccurred = nil
	if c != nil {
		c <- e
		close(c)
	}
}

func (ctl *FileTransferControl) sendAndCloseTransferFinished(v bool) {
	c := ctl.transferFinished
	ctl.transferFinished = nil
	if c != nil {
		c <- v
		close(c)
	}
}

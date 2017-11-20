package data

import "sync"

type transferUpdate struct {
	current, total int64
}

// FileTransferControl supplies the capabilities to control the file transfer
type FileTransferControl struct {
	sync.RWMutex                         // all updates to the channels need to be protected by this mutex
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

func (ctl *FileTransferControl) closeCancel() {
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.cancelTransfer != nil {
		close(ctl.cancelTransfer)
		ctl.cancelTransfer = nil
	}
}

func (ctl *FileTransferControl) WaitForFinish(k func()) {
	_, ok := <-ctl.transferFinished
	if ok {
		k()
		ctl.closeCancel()
	}
}

func (ctl *FileTransferControl) WaitForError(k func(error)) {
	e, ok := <-ctl.errorOccurred
	if ok {
		k(e)
		ctl.closeCancel()
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
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.cancelTransfer != nil {
		ctl.cancelTransfer <- true
	}
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
	ctl.RLock()
	defer ctl.RUnlock()
	if ctl.update != nil {
		ctl.update <- transferUpdate{current, total}
	}
}

func (ctl *FileTransferControl) CloseAll() {
	ctl.closeTransferFinished()
	ctl.closeUpdate()
	ctl.closeErrorOccurred()
}

func (ctl *FileTransferControl) closeTransferFinished() {
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.transferFinished != nil {
		close(ctl.transferFinished)
		ctl.transferFinished = nil
	}
}

func (ctl *FileTransferControl) closeUpdate() {
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.update != nil {
		close(ctl.update)
		ctl.update = nil
	}
}

func (ctl *FileTransferControl) closeErrorOccurred() {
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.errorOccurred != nil {
		close(ctl.errorOccurred)
		ctl.errorOccurred = nil
	}
}

func (ctl *FileTransferControl) sendAndCloseErrorOccurred(e error) {
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.errorOccurred != nil {
		ctl.errorOccurred <- e
		close(ctl.errorOccurred)
		ctl.errorOccurred = nil
	}
}

func (ctl *FileTransferControl) sendAndCloseTransferFinished(v bool) {
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.transferFinished != nil {
		ctl.transferFinished <- v
		close(ctl.transferFinished)
		ctl.transferFinished = nil
	}
}

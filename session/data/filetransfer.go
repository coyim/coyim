package data

import "sync"

type transferUpdate struct {
	current, total int64
}

// FileTransferControl supplies the capabilities to control the file transfer
type FileTransferControl struct {
	sync.RWMutex                                 // all updates to the channels need to be protected by this mutex
	cancelTransfer           chan bool           // one time use
	errorOccurred            chan error          // one time use
	update                   chan transferUpdate // will be called many times
	transferFinished         chan bool           // one time use
	OnEncryptionNotSupported func() bool         // called if encryption is not supported - the return value is true if we should continue without encryption
	EncryptionDecision       func(bool)          // called when we have decided to use encryption or not
}

func newFileTransferControl(c chan bool, e chan error, u chan transferUpdate, t chan bool) *FileTransferControl {
	return &FileTransferControl{cancelTransfer: c, errorOccurred: e, update: u, transferFinished: t}
}

// CreateFileTransferControl will return a new control object for file transfers
func CreateFileTransferControl(onNoEnc func() bool, encDecision func(bool)) *FileTransferControl {
	ctl := newFileTransferControl(make(chan bool), make(chan error), make(chan transferUpdate, 1000), make(chan bool))
	ctl.OnEncryptionNotSupported = onNoEnc
	ctl.EncryptionDecision = encDecision
	return ctl
}

func (ctl *FileTransferControl) closeCancel() {
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.cancelTransfer != nil {
		close(ctl.cancelTransfer)
		ctl.cancelTransfer = nil
	}
}

// WaitForFinish will wait for a transfer to be finished and calls the function after it's done - the argument to the
// function will be false if the transfer was declined and true if it finished
func (ctl *FileTransferControl) WaitForFinish(k func(bool)) {
	notDeclined, ok := <-ctl.transferFinished
	if ok {
		k(notDeclined)
		ctl.closeCancel()
	}
}

// WaitForError will wait for an error to occur and then call the function when it happens
func (ctl *FileTransferControl) WaitForError(k func(error)) {
	e, ok := <-ctl.errorOccurred
	if ok {
		k(e)
		ctl.closeCancel()
	}
}

// WaitForCancel will wait for a cancel event and call the given function when it happens
func (ctl *FileTransferControl) WaitForCancel(k func()) {
	if cancel, ok := <-ctl.cancelTransfer; ok && cancel {
		k()
		ctl.CloseAll()
	}
}

// WaitForUpdate will call the function on every update event
func (ctl *FileTransferControl) WaitForUpdate(k func(int64, int64)) {
	for upd := range ctl.update {
		k(upd.current, upd.total)
	}
}

// Cancel will cancel the file transfer
func (ctl *FileTransferControl) Cancel() {
	ctl.Lock()
	defer ctl.Unlock()
	if ctl.cancelTransfer != nil {
		ctl.cancelTransfer <- true
	}
}

// ReportError will report an error
func (ctl *FileTransferControl) ReportError(e error) {
	ctl.closeTransferFinished()
	ctl.closeUpdate()
	ctl.sendAndCloseErrorOccurred(e)
}

// ReportErrorNonblocking will report an error to the file transfer in a non-blocking manner
func (ctl *FileTransferControl) ReportErrorNonblocking(e error) {
	go ctl.ReportError(e)
}

// ReportFinished will be called when the file transfer is finished
func (ctl *FileTransferControl) ReportFinished() {
	ctl.closeErrorOccurred()
	ctl.closeUpdate()
	ctl.sendAndCloseTransferFinished(true)
}

// ReportDeclined will be called if the file transfer is declined
func (ctl *FileTransferControl) ReportDeclined() {
	ctl.closeErrorOccurred()
	ctl.closeUpdate()
	ctl.sendAndCloseTransferFinished(false)
}

// SendUpdate sends update information
func (ctl *FileTransferControl) SendUpdate(current, total int64) {
	ctl.RLock()
	defer ctl.RUnlock()
	if ctl.update != nil {
		ctl.update <- transferUpdate{current, total}
	}
}

// CloseAll closes all channels
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

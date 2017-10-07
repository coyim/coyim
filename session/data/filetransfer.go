package data

// FileTransferControl supplies the capabilities to control the file transfer
type FileTransferControl struct {
	CancelTransfer   chan bool  // one time use
	ErrorOccurred    chan error // one time use
	Update           chan int64 // will be called many times
	TransferFinished chan bool  // one time use
}

func (ctl FileTransferControl) ReportError(e error) {
	close(ctl.TransferFinished)
	close(ctl.Update)
	ctl.ErrorOccurred <- e
	close(ctl.ErrorOccurred)
}

func (ctl FileTransferControl) ReportFinished() {
	close(ctl.ErrorOccurred)
	close(ctl.Update)
	ctl.TransferFinished <- true
	close(ctl.TransferFinished)
}

func (ctl FileTransferControl) ReportErrorNonblocking(e error) {
	go ctl.ReportError(e)
}

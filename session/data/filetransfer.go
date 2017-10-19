package data

// FileTransferControl supplies the capabilities to control the file transfer
type FileTransferControl struct {
	CancelTransfer   chan bool  // one time use
	ErrorOccurred    chan error // one time use
	Update           chan int64 // will be called many times
	TransferFinished chan bool  // one time use
}

func (ctl FileTransferControl) ReportError(e error) {
	one := ctl.TransferFinished
	ctl.TransferFinished = nil
	if one != nil {
		close(one)
	}

	two := ctl.Update
	ctl.Update = nil
	if two != nil {
		close(two)
	}

	three := ctl.ErrorOccurred
	ctl.ErrorOccurred = nil
	if three != nil {
		three <- e
		close(three)
	}
}

func (ctl FileTransferControl) ReportFinished() {
	one := ctl.ErrorOccurred
	ctl.ErrorOccurred = nil
	if one != nil {
		close(one)
	}

	two := ctl.Update
	ctl.Update = nil
	if two != nil {
		close(two)
	}

	three := ctl.TransferFinished
	ctl.TransferFinished = nil
	if three != nil {
		three <- true
		close(three)
	}
}

func (ctl FileTransferControl) SendUpdate(v int64) {
	one := ctl.Update
	if one != nil {
		one <- v
	}
}

func (ctl FileTransferControl) ReportErrorNonblocking(e error) {
	go ctl.ReportError(e)
}

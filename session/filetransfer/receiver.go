package filetransfer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

// This file and the object inside it encapsulate everything necessary for receiving data - it should work for both the bytestream and IBB methods

// The interface should be simple:
// - It will optionally use encryption
// - It will have a method for checking if things are done
// - It will have one method "AddData" that adds a []byte.
// - It will NOT support base64
// - It will do the reporting of how much we have received
// - It will only write to the temp file. The real file handling and dir handling will be done outside

// Internally it will use a buffer, and a lock, and then a separate goroutine that pushes things to the file

type receiver struct {
	sync.Mutex

	buffer []byte

	ctx *recvContext

	writeChannel chan []byte
	newData      *sync.Cond

	hadError bool
	err      error
	done     bool

	toSendAtFinish   []byte
	fileNameAtFinish string
}

func (ctx *recvContext) createReceiver() *receiver {
	initialBuffer := make([]byte, 0, 4096)
	r := &receiver{
		buffer:       initialBuffer,
		ctx:          ctx,
		writeChannel: make(chan []byte, 10),
	}

	r.newData = sync.NewCond(r)

	go r.processReceive()
	go r.readAndRun()

	return r
}

func (r *receiver) Read(b []byte) (int, error) {
	neededLen := len(b)
	r.Lock()
	defer r.Unlock()
	for {
		if r.hadError {
			return 0, errors.New("error happened somewhere else")
		}

		if r.done {
			return 0, io.EOF
		}

		if len(r.buffer) > 0 {
			toCopy := neededLen
			if len(r.buffer) < toCopy {
				toCopy = len(r.buffer)
			}
			copy(b, r.buffer[:toCopy])
			r.buffer = r.buffer[toCopy:]
			return toCopy, nil
		}

		r.newData.Wait()
	}
}

func (r *receiver) saveError(e error) {
	r.Lock()
	r.err = e
	r.hadError = true
	r.newData.Broadcast()
	r.Unlock()
}

func (r *receiver) readAndRun() {
	ff, err := r.ctx.openDestinationTempFile()
	if err != nil {
		r.ctx.s.Warn(fmt.Sprintf("Failed to open temporary file: %v", err))
		removeInflightRecv(r.ctx.sid)
		r.saveError(err)
		return
	}

	totalWritten := int64(0)
	writes := 0

	reporting := func(v int) error {
		totalWritten += int64(v)
		writes++
		r.ctx.control.SendUpdate(totalWritten, r.ctx.size)
		return nil
	}

	rr, afterFinish := r.ctx.enc.wrapForReceiving(r)

	_, err = io.CopyN(io.MultiWriter(ff, &reportingWriter{report: reporting}), rr, r.ctx.size)
	if err != nil {
		r.ctx.s.Warn(fmt.Sprintf("Had error when trying to write to file: %v", err))
		r.ctx.control.ReportError(errors.New("Error writing to file"))
		closeAndIgnore(ff)
		_ = os.Remove(ff.Name())
		removeInflightRecv(r.ctx.sid)
		r.saveError(err)
		return
	}

	fstat, _ := ff.Stat()

	if totalWritten != r.ctx.size || fstat.Size() != totalWritten {
		r.ctx.s.Warn(fmt.Sprintf("Expected size of file to be %d, but was %d - this probably means the transfer was cancelled", r.ctx.size, fstat.Size()))
		err = errors.New("Incorrect final size of file - this implies the transfer was cancelled")
		r.ctx.control.ReportError(err)
		closeAndIgnore(ff)
		_ = os.Remove(ff.Name())
		removeInflightRecv(r.ctx.sid)
		r.saveError(err)
		return
	}

	toSend, err := afterFinish()
	if err != nil {
		r.ctx.s.Warn(fmt.Sprintf("Couldn't verify integrity of sent file: %v", err))
		r.ctx.control.ReportError(errors.New("Couldn't verify integrity of sent file"))
		closeAndIgnore(ff)
		_ = os.Remove(ff.Name())
		removeInflightRecv(r.ctx.sid)
		r.saveError(err)
		return
	}

	r.Lock()
	r.toSendAtFinish = toSend
	r.fileNameAtFinish = ff.Name()
	r.done = true
	r.newData.Broadcast()
	r.Unlock()

}

func (r *receiver) processReceive() {
	for {
		data := <-r.writeChannel
		r.Lock()
		r.buffer = append(r.buffer, data...)
		r.newData.Broadcast()
		r.Unlock()
	}
}

func (r *receiver) cancel() {
	r.saveError(errLocalCancel)
}

// addData adds the data for processing
// It will return any potential errors found during the process
func (r *receiver) Write(d []byte) (int, error) {
	if r.err != nil {
		e := r.err
		r.err = nil
		return 0, e
	}

	copyd := []byte{}
	copyd = append(copyd, d...)

	r.writeChannel <- copyd

	return len(d), nil
}

func (r *receiver) wait() ([]byte, string, bool, error) {
	r.Lock()
	defer r.Unlock()
	for {
		if r.hadError {
			return nil, "", false, r.err
		}
		if r.done {
			return r.toSendAtFinish, r.fileNameAtFinish, true, nil
		}
		r.newData.Wait()
	}
}

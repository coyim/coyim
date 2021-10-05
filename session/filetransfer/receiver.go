package filetransfer

import (
	"errors"
	"io"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
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

// The public interface of a receiver is quite small:
// type receiver interface{
//   io.Reader
//   io.Writer
// }
// - wait() ([]byte, string, bool, error)
// - cancel()

// Also, the io.Reader implementation seems to be incidental to the actual implementation
// It would be great to hide these implementation details

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

var ioCopy = io.Copy
var ioCopyN = io.CopyN
var ioMultiWriter = io.MultiWriter
var ioPipe = io.Pipe
var ioReadFull = io.ReadFull
var ioTeeReader = io.TeeReader

func (r *receiver) readAndRun() {
	ff, err := r.ctx.openDestinationTempFile()
	if err != nil {
		r.ctx.s.Log().WithError(err).Warn("Failed to open temporary file")
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

	_, err = ioCopyN(ioMultiWriter(ff, &reportingWriter{report: reporting}), rr, r.ctx.size)
	if err != nil {
		r.ctx.s.Log().WithError(err).Warn("Had error when trying to write to file")
		r.ctx.control.ReportError(errors.New("Error writing to file"))
		closeAndIgnore(ff)
		_ = os.Remove(ff.Name())
		removeInflightRecv(r.ctx.sid)
		r.saveError(err)
		return
	}

	fstat, _ := ff.Stat()

	// TODO: this seems like it can only happen if the encryption doesn't generate
	// the correct amount of data. If the transfer is cancelled from our
	// side, the previous error handling will trigger.
	// If the other side cancels, that will leave the whole code hanging, so we
	// likely will never get here.
	// I've tried to test this code, and I simply can't get it to trigger. It might
	// be dead.
	if totalWritten != r.ctx.size || fstat.Size() != totalWritten {
		r.ctx.s.Log().WithFields(log.Fields{
			"expected": r.ctx.size,
			"actual":   fstat.Size(),
		}).Warn("Unexpected file size - this probably means the transfer was cancelled")
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
		r.ctx.s.Log().WithError(err).Warn("Couldn't verify integrity of sent file")
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

// Write adds the data for processing
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

// wait will wait until the receipt of information has finished
// or an error has been encountered.
// if the transfer went well, it will return a true signalling all
// is ok, and no error. it will also optionally return
// some bytes to send back to the sender, as a "receipt"
// it wil return the file name where the received data
// was stored
// if an error happens, this will be returned instead.
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

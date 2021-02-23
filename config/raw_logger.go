package config

import (
	"bytes"
	"io"
	"sync"
)

// This rawLogger is a bit weird. For it to work properly, both
// _this_ instance, and the _other_ instance needs to share the same
// mutex.

type rawLogger struct {
	// the actual real writer where things will end up when flush is called
	out io.Writer
	// prefix will contain the prefix that each line will be preceded by
	prefix []byte
	// lock synchronizes the writes on this logger, but also the paired logger
	// in the other field
	lock *sync.Mutex
	// other is another raw logger which will be flushed when writes happen on this
	// logger
	other *rawLogger
	// buf contains the internal data waiting to be written
	buf []byte
}

func (r *rawLogger) Write(data []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.other.flush(); err != nil {
		return 0, nil
	}

	origLen := len(data)
	for len(data) > 0 {
		if newLine := bytes.IndexByte(data, '\n'); newLine >= 0 {
			r.buf = append(r.buf, data[:newLine]...)
			data = data[newLine+1:]
		} else {
			r.buf = append(r.buf, data...)
			data = nil
		}
	}

	return origLen, nil
}

func (r *rawLogger) flush() error {
	newLine := []byte{'\n'}

	if len(r.buf) == 0 {
		return nil
	}

	if _, err := r.out.Write(r.prefix); err != nil {
		return err
	}
	if _, err := r.out.Write(r.buf); err != nil {
		return err
	}
	if _, err := r.out.Write(newLine); err != nil {
		return err
	}
	r.buf = r.buf[:0]
	return nil
}

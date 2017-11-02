package filetransfer

// reportingWriter implements the Writer interface, which allows it to be used either together with io.MultiWriter or io.TeeReader to allow for updates during a process
// the report function can also cancel a process by returning an error
type reportingWriter struct {
	report func(int) error
}

func (rw *reportingWriter) Write(p []byte) (n int, err error) {
	v := len(p)
	e := rw.report(v)
	return v, e
}

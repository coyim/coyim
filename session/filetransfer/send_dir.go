package filetransfer

import (
	"io/ioutil"

	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
)

type dirSendContext struct {
	s                access.Session
	peer             string
	dir              string
	sid              string
	weWantToCancel   bool
	theyWantToCancel bool
	totalSent        int64
	control          *sdata.FileTransferControl
}

func (ctx *dirSendContext) startPackingDirectory() (<-chan string, <-chan error) {
	result := make(chan string)
	errorResult := make(chan error)

	go func() {
		tmpFile, e := ioutil.TempFile("", "coyim-packing")
		if e != nil {
			errorResult <- e
			return
		}
		defer tmpFile.Close()
		e = pack(ctx.dir, tmpFile)
		if e != nil {
			errorResult <- e
			return
		}
		result <- tmpFile.Name()
	}()

	return result, errorResult
}

const dirTransferProfile = "http://jabber.org/protocol/si/profile/directory-transfer"

func (ctx *dirSendContext) initSend() {
	result, errorResult := ctx.startPackingDirectory()

	_, err := discoverSupport(ctx.s, ctx.peer)
	if err != nil {
		ctx.control.ReportErrorNonblocking(err)
		return
	}

	go ctx.listenForCancellation()

	select {
	case tmpFile := <-result:
		ctx.offerSend(tmpFile)
	case e := <-errorResult:
		ctx.control.ReportErrorNonblocking(e)
	}
}

func (ctx *dirSendContext) listenForCancellation() {
	// TODO: fix
}

func (ctx *dirSendContext) offerSend(file string) {
	// TODO: fix
}

// InitSendDir starts the process of sending a directory to a peer
func InitSendDir(s access.Session, peer string, dir string) *sdata.FileTransferControl {
	ctx := &dirSendContext{
		s:       s,
		peer:    peer,
		dir:     dir,
		control: sdata.CreateFileTransferControl(),
	}
	go ctx.initSend()
	return ctx.control
}

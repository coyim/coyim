package filetransfer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
)

const dirTransferProfile = "http://jabber.org/protocol/si/profile/directory-transfer"

type dirSendContext struct {
	s                access.Session
	peer             string
	dir              string
	sid              string
	weWantToCancel   bool
	theyWantToCancel bool
	totalSent        int64
	control          *sdata.FileTransferControl
	fallback         *sendContext
	onErrorHook      func(*dirSendContext, error)
}

func (ctx *dirSendContext) onError(e error) {
	ctx.control.ReportErrorNonblocking(e)
	if ctx.onErrorHook != nil {
		ctx.onErrorHook(ctx, e)
	}
}

func (ctx *dirSendContext) startPackingDirectory() (<-chan string, <-chan error) {
	result := make(chan string)
	errorResult := make(chan error)

	go func() {
		tmpFile, e := ioutil.TempFile("", fmt.Sprintf("%s-directory-", filepath.Base(ctx.dir)))
		if e != nil {
			errorResult <- e
			return
		}
		e = pack(ctx.dir, tmpFile)
		if e != nil {
			errorResult <- e
			tmpFile.Close()
			return
		}
		newName := fmt.Sprintf("%v.zip", tmpFile.Name())
		tmpFile.Close()
		os.Rename(tmpFile.Name(), newName)
		result <- newName
	}()

	return result, errorResult
}

func (ctx *dirSendContext) initSend() {
	result, errorResult := ctx.startPackingDirectory()

	supported, err := discoverSupport(ctx.s, ctx.peer)
	if err != nil {
		ctx.onError(err)
		return
	}

	go ctx.listenForCancellation()

	select {
	case tmpFile := <-result:
		ctx.offerSend(tmpFile, supported)
	case e := <-errorResult:
		ctx.onError(e)
	}
}

func (ctx *dirSendContext) listenForCancellation() {
	ctx.control.WaitForCancel(func() {
		ctx.weWantToCancel = true
		if ctx.fallback != nil {
			ctx.fallback.weWantToCancel = true
		}
	})
}

func (ctx *dirSendContext) offerSendDirectory(file string) error {
	// TODO: for now, only supported on CoyIM
	ctx.sid = genSid(ctx.s.Conn())
	return nil
}

func (ctx *dirSendContext) offerSendDirectoryFallback(file string) error {
	// This one is a fallback for sending to clients that don't support directory sending, but do support file sending. We will simply send the packaged .zip file to them.
	fctx := &sendContext{
		s:       ctx.s,
		peer:    ctx.peer,
		file:    file,
		control: ctx.control,
		onFinishHook: func(_ *sendContext) {
			os.Remove(file)
		},
		onErrorHook: func(_ *sendContext, _ error) {
			os.Remove(file)
		},
	}
	ctx.fallback = fctx
	return fctx.offerSend()
}

func (ctx *dirSendContext) offerSend(file string, availableProfiles map[string]bool) error {
	if availableProfiles[dirTransferProfile] {
		return ctx.offerSendDirectory(file)
	}
	return ctx.offerSendDirectoryFallback(file)
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

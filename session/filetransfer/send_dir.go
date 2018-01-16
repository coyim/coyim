package filetransfer

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

const dirTransferProfile = "http://jabber.org/protocol/si/profile/directory-transfer"

type dirSendContext struct {
	dir string
	sc  *sendContext
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

	supported, err := discoverSupport(ctx.sc.s, ctx.sc.peer)
	if err != nil {
		ctx.sc.onError(err)
		return
	}

	go ctx.listenForCancellation()

	select {
	case tmpFile := <-result:
		ctx.offerSend(tmpFile, supported)
	case e := <-errorResult:
		ctx.sc.onError(e)
	}
}

func (ctx *dirSendContext) listenForCancellation() {
	ctx.sc.listenForCancellation()
}

func (ctx *sendContext) sendSIData(profile, file string, isDir bool) data.SI {
	if ctx.enc != nil {
		profile = encryptedTransferProfile
	}

	// TODO: Add Date and Hash here later?
	si := data.SI{
		ID:      ctx.sid,
		Profile: profile,
		Feature: data.FeatureNegotation{
			Form: data.Form{
				Type: "form",
				Fields: []data.FormFieldX{
					{
						Var:     "stream-method",
						Type:    "list-single",
						Options: calculateAvailableSendOptions(ctx.s),
					},
				},
			},
		},
	}

	if ctx.enc != nil {
		si.EncryptedData = &data.EncryptedData{
			Name: filepath.Base(file),
			Size: ctx.size,
			EncryptionParameters: &data.EncryptionParameters{
				Type:          "aes128-ctr",
				IV:            base64.StdEncoding.EncodeToString(ctx.enc.iv),
				MAC:           "hmac-sha256",
				EncryptionKey: &data.EncryptionKeyParameter{Type: ctx.enc.keyType},
				MACKey:        &data.EncryptionKeyParameter{Type: ctx.enc.keyType},
			},
		}
		if isDir {
			si.EncryptedData.Type = "directory"
		} else {
			si.EncryptedData.Type = "file"
		}
	} else {
		si.File = &data.File{
			Name: filepath.Base(file),
			Size: ctx.size,
		}
	}

	return si
}

// we assume that ctx.sc.file points to a valid file, since it's generated in previous code. thus, we don't check for existance.
func (ctx *dirSendContext) offerSendDirectory() error {
	fstat, _ := os.Stat(ctx.sc.file)
	ctx.sc.sid = genSid(ctx.sc.s.Conn())
	ctx.sc.size = fstat.Size()

	toSend := ctx.sc.sendSIData(dirTransferProfile, ctx.dir, true)

	var siq data.SI
	nonblockIQ(ctx.sc.s, ctx.sc.peer, "set", toSend, &siq, func(*data.ClientIQ) {
		if !isValidSubmitForm(siq) {
			ctx.sc.onError(errors.New("Invalid data sent from peer for directory sending"))
			return
		}
		prof := siq.Feature.Form.Fields[0].Values[0]
		if f, ok := supportedSendingMechanisms[prof]; ok {
			notifyUserThatSendStarted(prof, ctx.sc.s, ctx.sc.file, ctx.sc.peer)
			addInflightSend(ctx.sc)
			f(ctx.sc)
			return
		}
		ctx.sc.onError(errors.New("Invalid sending mechanism sent from peer for directory sending"))
	}, func(_ *data.ClientIQ, e error) {
		ctx.sc.onError(e)
	})

	return nil
}

// This one is a fallback for sending to clients that don't support directory sending, but do support file sending. We will simply send the packaged .zip file to them.
func (ctx *dirSendContext) offerSendDirectoryFallback() error {
	return ctx.sc.offerSend()
}

func (ctx *dirSendContext) offerSend(file string, availableProfiles map[string]bool) error {
	ctx.sc.file = file
	if availableProfiles[dirTransferProfile] {
		return ctx.offerSendDirectory()
	}
	return ctx.offerSendDirectoryFallback()
}

// InitSendDir starts the process of sending a directory to a peer
func InitSendDir(s access.Session, peer jid.Any, dir string, encrypted bool) *sdata.FileTransferControl {
	ctx := &dirSendContext{
		sc: &sendContext{
			s:       s,
			enc:     generateEncryptionParameters(encrypted, func() []byte { return s.CreateSymmetricKeyFor(peer) }, "external"),
			peer:    peer.String(),
			control: sdata.CreateFileTransferControl(),
			onFinishHook: func(ctx *sendContext) {
				os.Remove(ctx.file)
			},
			onErrorHook: func(ctx *sendContext, _ error) {
				os.Remove(ctx.file)
			},
		},
		dir: dir,
	}
	go ctx.initSend()
	return ctx.sc.control
}

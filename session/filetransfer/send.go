package filetransfer

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
)

func registerSendFileTransferMethod(name string, dispatch func(access.Session, *sendContext)) {
	supportedSendingMechanisms[name] = dispatch
}

var supportedSendingMechanisms = map[string]func(access.Session, *sendContext){}

func (ctx *sendContext) discoverSupport(s access.Session) (profiles []string, err error) {
	if res, ok := s.Conn().DiscoveryFeatures(ctx.peer); ok {
		foundSI := false
		for _, feature := range res {
			if feature == "http://jabber.org/protocol/si" {
				foundSI = true
			} else if strings.HasPrefix(feature, "http://jabber.org/protocol/si/profile/") {
				profiles = append(profiles, feature)
			}
		}

		if !foundSI {
			return nil, errors.New("Peer doesn't support stream initiation")
		}

		if len(profiles) == 0 {
			return nil, errors.New("Peer doesn't support any stream initiation profiles")
		}

		return profiles, nil
	}
	return nil, errors.New("Problem discovering the features of the peer")
}

func genSid(c interfaces.Conn) string {
	var buf [8]byte
	if _, err := c.Rand().Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return fmt.Sprintf("sid%d", binary.LittleEndian.Uint64(buf[:]))
}

const fileTransferProfile = "http://jabber.org/protocol/si/profile/file-transfer"

func calculateAvailableSendOptions() []data.FormFieldOptionX {
	res := []data.FormFieldOptionX{}
	for k, _ := range supportedSendingMechanisms {
		res = append(res, data.FormFieldOptionX{Value: k})
	}
	return res
}

func (ctx *sendContext) offerSend(s access.Session) error {
	fstat, e := os.Stat(ctx.file)
	if e != nil {
		return e
	}
	ctx.sid = genSid(s.Conn())

	// TODO: Add Date and Hash here later?
	toSend := data.SI{
		ID:      ctx.sid,
		Profile: fileTransferProfile,
		File: &data.File{
			Name: filepath.Base(ctx.file),
			Size: fstat.Size(),
		},
		Feature: data.FeatureNegotation{
			Form: data.Form{
				Type: "form",
				Fields: []data.FormFieldX{
					{
						Var:     "stream-method",
						Type:    "list-single",
						Options: calculateAvailableSendOptions(),
					},
				},
			},
		},
	}

	res, _, e2 := s.Conn().SendIQ(ctx.peer, "set", toSend)
	if e2 != nil {
		return e2
	}
	go ctx.waitForResultToStartFileSend(s, res)

	return nil
}

type sendContext struct {
	peer             string
	file             string
	sid              string
	weWantToCancel   bool
	theyWantToCancel bool
	totalSent        int64
	control          sdata.FileTransferControl
}

func (ctx *sendContext) notifyUserThatSendStarted(s access.Session) {
	s.Info(fmt.Sprintf("Started sending of %v to %v", ctx.file, ctx.peer))
}

func isValidSubmitForm(siq data.SI) bool {
	return siq.Feature.Form.Type == "submit" &&
		len(siq.Feature.Form.Fields) == 1 &&
		siq.Feature.Form.Fields[0].Var == "stream-method" &&
		len(siq.Feature.Form.Fields[0].Values) == 1
}

func (ctx *sendContext) listenForCancellation() {
	if cancel, ok := <-ctx.control.CancelTransfer; ok && cancel {
		ctx.weWantToCancel = true
		close(ctx.control.TransferFinished)
		close(ctx.control.Update)
		close(ctx.control.ErrorOccurred)
	}
}

func (ctx *sendContext) waitForResultToStartFileSend(s access.Session, reply <-chan data.Stanza) {
	r, ok := <-reply
	if ok {
		switch ciq := r.Value.(type) {
		case *data.ClientIQ:
			if ciq.Type != "result" {
				ctx.control.ReportErrorNonblocking(errors.New("Received error from peer when offering to send file"))
				return
			}

			var siq data.SI
			if err := xml.NewDecoder(bytes.NewBuffer(ciq.Query)).Decode(&siq); err != nil {
				ctx.control.ReportErrorNonblocking(err)
				return
			}
			if !isValidSubmitForm(siq) {
				ctx.control.ReportErrorNonblocking(errors.New("Invalid data sent from peer for file sending"))
				return
			}
			prof := siq.Feature.Form.Fields[0].Values[0]
			if f, ok := supportedSendingMechanisms[prof]; ok {
				ctx.notifyUserThatSendStarted(s)
				addInflightSend(ctx)
				f(s, ctx)
				return
			}
			ctx.control.ReportErrorNonblocking(errors.New("Invalid sending mechanism sent from peer for file sending"))
			return
		default:
			ctx.control.ReportErrorNonblocking(errors.New("Invalid stanza type - this shouldn't happen"))
			return
		}
	}
	ctx.control.ReportErrorNonblocking(errors.New("No response received to offer of sending a file"))
}

func createNewFileTransferControl() sdata.FileTransferControl {
	return sdata.FileTransferControl{
		CancelTransfer:   make(chan bool),
		ErrorOccurred:    make(chan error),
		Update:           make(chan int64, 1000),
		TransferFinished: make(chan bool),
	}
}

// InitSend starts the process of sending a file to a peer
func InitSend(s access.Session, peer string, file string) sdata.FileTransferControl {
	ctx := &sendContext{
		peer:    peer,
		file:    file,
		control: createNewFileTransferControl(),
	}

	_, err := ctx.discoverSupport(s)
	if err != nil {
		ctx.control.ReportErrorNonblocking(err)
		return ctx.control
	}

	go ctx.listenForCancellation()
	ctx.offerSend(s)
	return ctx.control
}

var inflightSends struct {
	sync.RWMutex
	transfers map[string]*sendContext
}

func init() {
	inflightSends.transfers = make(map[string]*sendContext)
}

func addInflightSend(ctx *sendContext) {
	inflightSends.Lock()
	defer inflightSends.Unlock()
	inflightSends.transfers[ctx.sid] = ctx
}

func getInflightSend(id string) (result *sendContext, ok bool) {
	inflightSends.RLock()
	defer inflightSends.RUnlock()
	result, ok = inflightSends.transfers[id]
	return
}

func removeInflightSend(ctx *sendContext) {
	inflightSends.Lock()
	defer inflightSends.Unlock()
	delete(inflightSends.transfers, ctx.sid)
}

package filetransfer

import (
	"encoding/binary"
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

func registerSendFileTransferMethod(name string, dispatch func(*sendContext), isCurrentlyValid func(string, *sendContext) bool) {
	supportedSendingMechanisms[name] = dispatch
	isSendingMechanismCurrentlyValid[name] = isCurrentlyValid
}

var supportedSendingMechanisms = map[string]func(*sendContext){}
var isSendingMechanismCurrentlyValid = map[string]func(string, *sendContext) bool{}

func discoverSupport(s access.Session, p string) (profiles []string, err error) {
	if res, ok := s.Conn().DiscoveryFeatures(p); ok {
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

func (ctx *sendContext) calculateAvailableSendOptions() []data.FormFieldOptionX {
	res := []data.FormFieldOptionX{}
	for k, _ := range supportedSendingMechanisms {
		if isSendingMechanismCurrentlyValid[k](k, ctx) {
			res = append(res, data.FormFieldOptionX{Value: k})
		}
	}
	return res
}

func (ctx *sendContext) offerSend() error {
	fstat, e := os.Stat(ctx.file)
	if e != nil {
		return e
	}
	ctx.sid = genSid(ctx.s.Conn())

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
						Options: ctx.calculateAvailableSendOptions(),
					},
				},
			},
		},
	}

	var siq data.SI
	nonblockIQ(ctx.s, ctx.peer, "set", toSend, &siq, func(*data.ClientIQ) {
		if !isValidSubmitForm(siq) {
			ctx.control.ReportErrorNonblocking(errors.New("Invalid data sent from peer for file sending"))
			return
		}
		prof := siq.Feature.Form.Fields[0].Values[0]
		if f, ok := supportedSendingMechanisms[prof]; ok {
			ctx.notifyUserThatSendStarted(prof)
			addInflightSend(ctx)
			f(ctx)
			return
		}
		ctx.control.ReportErrorNonblocking(errors.New("Invalid sending mechanism sent from peer for file sending"))
	}, func(_ *data.ClientIQ, e error) {
		ctx.control.ReportErrorNonblocking(e)
	})

	return nil
}

type sendContext struct {
	s                access.Session
	peer             string
	file             string
	sid              string
	weWantToCancel   bool
	theyWantToCancel bool
	totalSent        int64
	control          *sdata.FileTransferControl
}

func (ctx *sendContext) notifyUserThatSendStarted(method string) {
	ctx.s.Info(fmt.Sprintf("Started sending of %v to %v using %v", ctx.file, ctx.peer, method))
}

func isValidSubmitForm(siq data.SI) bool {
	return siq.Feature.Form.Type == "submit" &&
		len(siq.Feature.Form.Fields) == 1 &&
		siq.Feature.Form.Fields[0].Var == "stream-method" &&
		len(siq.Feature.Form.Fields[0].Values) == 1
}

func (ctx *sendContext) listenForCancellation() {
	ctx.control.WaitForCancel(func() {
		ctx.weWantToCancel = true
	})
}

func (ctx *sendContext) initSend() {
	_, err := discoverSupport(ctx.s, ctx.peer)
	if err != nil {
		ctx.control.ReportErrorNonblocking(err)
		return
	}

	go ctx.listenForCancellation()
	ctx.offerSend()
}

// InitSend starts the process of sending a file to a peer
func InitSend(s access.Session, peer string, file string) *sdata.FileTransferControl {
	ctx := &sendContext{
		peer:    peer,
		file:    file,
		control: sdata.CreateFileTransferControl(),
		s:       s,
	}
	go ctx.initSend()
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

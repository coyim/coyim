package filetransfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"

	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

const fileTransferProfile = "http://jabber.org/protocol/si/profile/file-transfer"
const encryptedTransferProfile = "http://jabber.org/protocol/si/profile/encrypted-data-transfer"

func registerSendFileTransferMethod(name string, dispatch func(*sendContext), isCurrentlyValid func(string, hasConnectionAndConfig) bool) {
	supportedSendingMechanisms[name] = dispatch
	isSendingMechanismCurrentlyValid[name] = isCurrentlyValid
}

var supportedSendingMechanisms = map[string]func(*sendContext){}
var isSendingMechanismCurrentlyValid = map[string]func(string, hasConnectionAndConfig) bool{}

func discoverSupport(s hasConnection, p string) (profiles map[string]bool, err error) {
	profiles = make(map[string]bool)
	if res, ok := s.Conn().DiscoveryFeatures(p); ok {
		foundSI := false
		for _, feature := range res {
			if feature == "http://jabber.org/protocol/si" {
				foundSI = true
			} else if strings.HasPrefix(feature, "http://jabber.org/protocol/si/profile/") {
				profiles[feature] = true
			}
		}

		if !foundSI {
			return nil, errors.New("Peer doesn't support stream initiation")
		}

		if len(profiles) == 0 {
			return nil, errors.New("Peer doesn't support any stream initiation profiles")
		}

		return
	}
	return profiles, errors.New("Problem discovering the features of the peer")
}

func genSid(c interfaces.Conn) string {
	var buf [8]byte
	if _, err := c.Rand().Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return fmt.Sprintf("sid%d", binary.LittleEndian.Uint64(buf[:]))
}

func calculateAvailableSendOptions(s hasConnectionAndConfig) []data.FormFieldOptionX {
	res := []data.FormFieldOptionX{}
	for k := range supportedSendingMechanisms {
		if isSendingMechanismCurrentlyValid[k](k, s) {
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
	ctx.size = fstat.Size()

	toSend := ctx.sendSIData(fileTransferProfile, ctx.file, false)

	var siq data.SI
	nonblockIQ(ctx.s, ctx.peer, "set", toSend, &siq, func(*data.ClientIQ) {
		if !isValidSubmitForm(siq) {
			ctx.onError(errors.New("Invalid data sent from peer for file sending"))
			return
		}
		prof := siq.Feature.Form.Fields[0].Values[0]
		if f, ok := supportedSendingMechanisms[prof]; ok {
			notifyUserThatSendStarted(prof, ctx.s, ctx.file, ctx.peer)
			addInflightSend(ctx)
			f(ctx)
			return
		}
		ctx.onError(errors.New("Invalid sending mechanism sent from peer for file sending"))
	}, func(stanza *data.ClientIQ, e error) {
		if stanza.Error.Code == "403" {
			ctx.onDecline()
		} else {
			ctx.onError(e)
		}
	})

	return nil
}

type sendContext struct {
	s                hasConnectionAndConfigAndLog
	peer             string
	file             string
	sid              string
	size             int64
	totalSize        int64
	enc              *encryptionParameters
	weWantToCancel   bool
	theyWantToCancel bool
	totalSent        int64
	control          *sdata.FileTransferControl
	onFinishHook     func(*sendContext)
	onErrorHook      func(*sendContext, error)
	onUpdateHook     func(*sendContext, int64)
	onDeclineHook    func(*sendContext)
}

func (ctx *sendContext) onFinish() {
	ctx.control.ReportFinished()
	removeInflightSend(ctx)
	if ctx.onFinishHook != nil {
		ctx.onFinishHook(ctx)
	}
}
func (ctx *sendContext) onError(e error) {
	ctx.control.ReportErrorNonblocking(e)
	removeInflightSend(ctx)
	if ctx.onErrorHook != nil {
		ctx.onErrorHook(ctx, e)
	}
}
func (ctx *sendContext) onUpdate(v int) {
	ctx.totalSent += int64(v)
	ctx.control.SendUpdate(ctx.totalSent, ctx.totalSize)
	if ctx.onUpdateHook != nil {
		ctx.onUpdateHook(ctx, ctx.totalSent)
	}
}

func (ctx *sendContext) onDecline() {
	ctx.control.ReportDeclined()
	removeInflightSend(ctx)
	if ctx.onDeclineHook != nil {
		ctx.onDeclineHook(ctx)
	}
}

func notifyUserThatSendStarted(method string, s hasLog, file, peer string) {
	s.Log().WithFields(log.Fields{
		"file":   file,
		"peer":   peer,
		"method": method,
	}).Info("Started sending file to peer using method")
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
	vals, err := discoverSupport(ctx.s, ctx.peer)
	if err != nil {
		ctx.onError(err)
		return
	}

	if !vals[encryptedTransferProfile] || ctx.enc == nil {
		if ctx.control.OnEncryptionNotSupported != nil && ctx.control.OnEncryptionNotSupported() {
			ctx.enc = nil
		} else {
			ctx.onError(errors.New("will not send unencrypted"))
			return
		}
	}

	if ctx.control.EncryptionDecision != nil {
		ctx.control.EncryptionDecision(ctx.enc != nil)
	}

	go ctx.listenForCancellation()
	_ = ctx.offerSend()
}

// InitSend starts the process of sending a file to a peer
func InitSend(s hasConnectionAndConfigAndLogAndHasSymmetricKey, peer jid.Any, file string, onNoEnc func() bool, encDecision func(bool)) *sdata.FileTransferControl {
	ctx := &sendContext{
		peer:    peer.String(),
		file:    file,
		control: sdata.CreateFileTransferControl(onNoEnc, encDecision),
		s:       s,
		enc:     generateEncryptionParameters(true, func() []byte { return s.CreateSymmetricKeyFor(peer) }, "external"),
	}
	go ctx.initSend()
	return ctx.control
}

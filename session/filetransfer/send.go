package filetransfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
)

const fileTransferProfile = "http://jabber.org/protocol/si/profile/file-transfer"

func registerSendFileTransferMethod(name string, dispatch func(*sendContext), isCurrentlyValid func(string, access.Session) bool) {
	supportedSendingMechanisms[name] = dispatch
	isSendingMechanismCurrentlyValid[name] = isCurrentlyValid
}

var supportedSendingMechanisms = map[string]func(*sendContext){}
var isSendingMechanismCurrentlyValid = map[string]func(string, access.Session) bool{}

func discoverSupport(s access.Session, p string) (profiles map[string]bool, err error) {
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

func calculateAvailableSendOptions(s access.Session) []data.FormFieldOptionX {
	res := []data.FormFieldOptionX{}
	for k, _ := range supportedSendingMechanisms {
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

	toSend := sendSIData(ctx.sid, fileTransferProfile, ctx.file, ctx.size, ctx.s)

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
	}, func(_ *data.ClientIQ, e error) {
		ctx.onError(e)
	})

	return nil
}

type sendContext struct {
	s                access.Session
	peer             string
	file             string
	sid              string
	size             int64
	weWantToCancel   bool
	theyWantToCancel bool
	totalSent        int64
	control          *sdata.FileTransferControl
	onFinishHook     func(*sendContext)
	onErrorHook      func(*sendContext, error)
	onUpdateHook     func(*sendContext, int64)
}

func (ctx *sendContext) onFinish() {
	fmt.Printf("onFinish()\n")
	fmt.Printf("  onFinish()-reportFinished\n")
	ctx.control.ReportFinished()
	fmt.Printf("  onFinish()-after reportFinished\n")
	fmt.Printf("  onFinish()-remove inflight send\n")
	removeInflightSend(ctx)
	fmt.Printf("  onFinish()-after remove inflight send\n")
	if ctx.onFinishHook != nil {
		fmt.Printf("  onFinish()-onfinishhook\n")
		ctx.onFinishHook(ctx)
		fmt.Printf("  onFinish()-after onfinishhook\n")
	}
}
func (ctx *sendContext) onError(e error) {
	fmt.Printf("sendContext.onError(%#v)\n", e)
	ctx.control.ReportErrorNonblocking(e)
	removeInflightSend(ctx)
	if ctx.onErrorHook != nil {
		ctx.onErrorHook(ctx, e)
	}
}
func (ctx *sendContext) onUpdate(v int) {
	ctx.totalSent += int64(v)
	ctx.control.SendUpdate(ctx.totalSent, ctx.size)
	if ctx.onUpdateHook != nil {
		ctx.onUpdateHook(ctx, ctx.totalSent)
	}
}

func notifyUserThatSendStarted(method string, s access.Session, file, peer string) {
	s.Info(fmt.Sprintf("Started sending of %v to %v using %v", file, peer, method))
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
		ctx.onError(err)
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

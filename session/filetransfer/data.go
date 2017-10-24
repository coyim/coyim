package filetransfer

import (
	"sync"

	"github.com/coyim/coyim/xmpp/data"
)

var inflightRecvs struct {
	sync.RWMutex
	transfers map[string]*recvContext
}

var inflightSends struct {
	sync.RWMutex
	transfers map[string]*sendContext
}

func init() {
	inflightRecvs.transfers = make(map[string]*recvContext)
	inflightSends.transfers = make(map[string]*sendContext)
}

var iqErrorBadRequest = data.ErrorReply{
	Type:   "cancel",
	Code:   400,
	Error:  data.ErrorBadRequest{},
	Error2: data.ErrorNoValidStreams{},
}

var iqErrorForbidden = data.ErrorReply{
	Type:  "cancel",
	Code:  403,
	Error: data.ErrorForbidden{},
	Text:  "Offer Declined",
}

var iqErrorNotAcceptable = data.ErrorReply{
	Type:  "cancel",
	Error: data.ErrorNotAcceptable{},
}

var iqErrorItemNotFound = data.ErrorReply{
	Type:  "cancel",
	Error: data.ErrorItemNotFound{},
}

var iqErrorUnexpectedRequest = data.ErrorReply{
	Type:  "cancel",
	Error: data.ErrorUnexpectedRequest{},
}

var iqErrorIBBBadRequest = data.ErrorReply{
	Type:  "cancel",
	Error: data.ErrorBadRequest{},
}

func addInflightRecv(ctx *recvContext) {
	inflightRecvs.Lock()
	defer inflightRecvs.Unlock()
	inflightRecvs.transfers[ctx.sid] = ctx
}

func getInflightRecv(id string) (result *recvContext, ok bool) {
	inflightRecvs.RLock()
	defer inflightRecvs.RUnlock()
	result, ok = inflightRecvs.transfers[id]
	return
}

func setInflightRecvDestination(id, destination string) {
	inflightRecvs.RLock()
	defer inflightRecvs.RUnlock()
	inflightRecvs.transfers[id].destination = destination
}

func removeInflightRecv(id string) {
	inflightRecvs.Lock()
	defer inflightRecvs.Unlock()
	delete(inflightRecvs.transfers, id)
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

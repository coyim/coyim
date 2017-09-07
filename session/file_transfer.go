package session

import (
	"errors"
	"fmt"
	"sync"

	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp/data"
)

// TODO: implement http://jabber.org/protocol/bytestreams

var supportedFileTransferMethods = map[string]int{
	"http://jabber.org/protocol/bytestreams": 1,
	"http://jabber.org/protocol/ibb":         0,
}

var fileTransferCancelListeners = map[string]func(*session, inflightFileTransfer){
	"http://jabber.org/protocol/ibb":         fileTransferIbbWaitForCancel,
	"http://jabber.org/protocol/bytestreams": fileTransferBytestreamWaitForCancel,
}

type inflightFileTransferStatus struct {
	destination string
	opaque      interface{}
}

type inflightFileTransfer struct {
	id      string
	mime    string
	options []string
	date    string
	hash    string
	name    string
	size    int64
	desc    string
	rng     struct {
		length *int
		offset *int
	}
	peer            string
	status          *inflightFileTransferStatus
	cancelChannel   <-chan bool
	errorChannel    chan<- error
	updateChannel   chan<- int64
	finishedChannel chan<- bool
}

var inflightFileTransfers struct {
	sync.RWMutex
	transfers map[string]inflightFileTransfer
}

func init() {
	inflightFileTransfers.transfers = make(map[string]inflightFileTransfer)
}

func extractFileTransferOptions(f data.Form) ([]string, error) {
	if f.Type != "form" || len(f.Fields) != 1 || f.Fields[0].Var != "stream-method" || f.Fields[0].Type != "list-single" {
		return nil, fmt.Errorf("Invalid form for file transfer initiation: %#v", f)
	}
	var result []string
	for _, opt := range f.Fields[0].Options {
		result = append(result, opt.Value)
	}
	return result, nil
}

// chooseAppropriateFileTransferOptionFrom returns the file transfer option that has the highest score
// or not OK if no acceptable options are available
func chooseAppropriateFileTransferOptionFrom(options []string) (best string, ok bool) {
	bestScore := -1
	for _, opt := range options {
		score, has := supportedFileTransferMethods[opt]
		if has {
			ok = true
			if score > bestScore {
				bestScore = score
				best = opt
			}
		}

	}
	return
}

func getInflightFileTransfer(id string) (result inflightFileTransfer, ok bool) {
	inflightFileTransfers.RLock()
	defer inflightFileTransfers.RUnlock()
	result, ok = inflightFileTransfers.transfers[id]
	return
}

func setInflightFileTransferDestination(id, destination string) {
	inflightFileTransfers.RLock()
	defer inflightFileTransfers.RUnlock()
	inflightFileTransfers.transfers[id].status.destination = destination
}

func removeInflightFileTransfer(id string) {
	inflightFileTransfers.Lock()
	defer inflightFileTransfers.Unlock()
	delete(inflightFileTransfers.transfers, id)
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

func iqResultChosenStreamMethod(opt string) data.SI {
	return data.SI{
		File: &data.File{},
		Feature: data.FeatureNegotation{
			Form: data.Form{
				Type: "submit",
				Fields: []data.FormFieldX{
					{Var: "stream-method", Values: []string{opt}},
				},
			},
		},
	}
}

func (ift inflightFileTransfer) reportError(e error) {
	close(ift.finishedChannel)
	close(ift.updateChannel)
	ift.errorChannel <- e
	close(ift.errorChannel)
}

func (ift inflightFileTransfer) reportFinished() {
	close(ift.errorChannel)
	close(ift.updateChannel)
	ift.finishedChannel <- true
	close(ift.finishedChannel)
}

func waitForFileTransferUserAcceptance(s *session, stanza *data.ClientIQ, si data.SI, acceptResult <-chan *string, ift inflightFileTransfer) {
	result := <-acceptResult

	var error *data.ErrorReply
	if result != nil {
		opt, ok := chooseAppropriateFileTransferOptionFrom(ift.options)
		if ok {
			setInflightFileTransferDestination(si.ID, *result)
			s.sendIQResult(stanza, iqResultChosenStreamMethod(opt))
			go fileTransferCancelListeners[opt](s, ift)
			return
		}
		ift.reportError(errors.New("No mutually acceptable file transfer methods available"))
		error = &iqErrorBadRequest
	} else {
		error = &iqErrorForbidden
	}
	removeInflightFileTransfer(si.ID)
	s.sendIQError(stanza, *error)
}

func registerNewFileTransfer(si data.SI, options []string, stanza *data.ClientIQ, f *data.File, cc <-chan bool, ec chan<- error, uc chan<- int64, fc chan<- bool) inflightFileTransfer {
	ift := inflightFileTransfer{
		id:              si.ID,
		mime:            si.MIMEType,
		options:         options,
		date:            f.Date,
		hash:            f.Hash,
		name:            f.Name,
		size:            f.Size,
		desc:            f.Desc,
		peer:            stanza.From,
		status:          &inflightFileTransferStatus{},
		cancelChannel:   cc,
		errorChannel:    ec,
		updateChannel:   uc,
		finishedChannel: fc,
	}

	if f.Range != nil {
		ift.rng.length = f.Range.Length
		ift.rng.offset = f.Range.Offset
	}

	inflightFileTransfers.Lock()
	defer inflightFileTransfers.Unlock()
	inflightFileTransfers.transfers[si.ID] = ift
	return ift
}

func fileStreamInitIQ(s *session, stanza *data.ClientIQ, si data.SI) (ret interface{}, ignore bool) {
	var options []string
	var err error
	if options, err = extractFileTransferOptions(si.Feature.Form); err != nil {
		s.warn(fmt.Sprintf("Failed to parse stream initiation: %v", err))
		return nil, false
	}

	f := si.File

	cancelChannel := make(chan bool)
	errorChannel := make(chan error)
	updateChannel := make(chan int64, 1000)
	finishedChannel := make(chan bool)
	ift := registerNewFileTransfer(si, options, stanza, f, cancelChannel, errorChannel, updateChannel, finishedChannel)

	acceptResult := make(chan *string)
	go waitForFileTransferUserAcceptance(s, stanza, si, acceptResult, ift)

	s.publishEvent(events.FileTransfer{
		Session:          s,
		Peer:             stanza.From,
		Mime:             f.Hash,
		DateLastModified: f.Date,
		Name:             f.Name,
		Size:             f.Size,
		Description:      f.Desc,
		Answer:           acceptResult,
		CancelTransfer:   cancelChannel,
		ErrorOccurred:    errorChannel,
		Update:           updateChannel,
		TransferFinished: finishedChannel,
	})

	return nil, true
}

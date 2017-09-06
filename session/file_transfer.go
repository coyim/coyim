package session

import (
	"fmt"
	"sync"

	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp/data"
)

var supportedFileTransferMethods = map[string]int{
	"http://jabber.org/protocol/bytestreams": 0,
	"http://jabber.org/protocol/ibb":         100, //TODO: this should never be the case in real life, but we use it while developing ibb and bytestreams
}

type inflightFileTransferStatus struct {
	state  string
	opaque interface{}
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
	status *inflightFileTransferStatus
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

func setInflightFileTransferState(id, state string) {
	inflightFileTransfers.RLock()
	defer inflightFileTransfers.RUnlock()
	inflightFileTransfers.transfers[id].status.state = state
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

func waitForFileTransferUserAcceptance(s *session, stanza *data.ClientIQ, si data.SI, acceptResult chan bool, options []string) {
	result := <-acceptResult
	close(acceptResult)
	var error *data.ErrorReply

	if result {
		opt, ok := chooseAppropriateFileTransferOptionFrom(options)
		if ok {
			setInflightFileTransferState(si.ID, "accepted")
			s.sendIQResult(stanza, iqResultChosenStreamMethod(opt))
			return
		}
		error = &iqErrorBadRequest
	} else {
		error = &iqErrorForbidden
	}
	removeInflightFileTransfer(si.ID)
	s.sendIQError(stanza, *error)
}

func registerNewFileTransfer(si data.SI, options []string, f *data.File) {
	ift := inflightFileTransfer{
		id:      si.ID,
		mime:    si.MIMEType,
		options: options,
		date:    f.Date,
		hash:    f.Hash,
		name:    f.Name,
		size:    f.Size,
		desc:    f.Desc,
		status:  &inflightFileTransferStatus{},
	}

	if f.Range != nil {
		ift.rng.length = f.Range.Length
		ift.rng.offset = f.Range.Offset
	}

	inflightFileTransfers.Lock()
	defer inflightFileTransfers.Unlock()
	inflightFileTransfers.transfers[si.ID] = ift
}

func fileStreamInitIQ(s *session, stanza *data.ClientIQ, si data.SI) (ret interface{}, ignore bool) {
	var options []string
	var err error
	if options, err = extractFileTransferOptions(si.Feature.Form); err != nil {
		s.warn(fmt.Sprintf("Failed to parse stream initiation: %v", err))
		return nil, false
	}

	f := si.File
	registerNewFileTransfer(si, options, f)

	acceptResult := make(chan bool, 1)
	go waitForFileTransferUserAcceptance(s, stanza, si, acceptResult, options)

	s.publishEvent(events.FileTransfer{
		Session:          s,
		Peer:             stanza.From,
		Mime:             f.Hash,
		DateLastModified: f.Date,
		Name:             f.Name,
		Size:             f.Size,
		Description:      f.Desc,
		Answer:           acceptResult,
	})

	return nil, true
}

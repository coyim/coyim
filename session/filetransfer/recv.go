package filetransfer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp/data"
)

var supportedFileTransferMethods = map[string]int{
	"http://jabber.org/protocol/bytestreams": 1,
	"http://jabber.org/protocol/ibb":         0,
}

var fileTransferCancelListeners = map[string]func(access.Session, inflight){
	"http://jabber.org/protocol/ibb":         ibbWaitForCancel,
	"http://jabber.org/protocol/bytestreams": bytestreamWaitForCancel,
}

type inflightStatus struct {
	destination string
	opaque      interface{}
}

type inflight struct {
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
	status          *inflightStatus
	cancelChannel   <-chan bool
	errorChannel    chan<- error
	updateChannel   chan<- int64
	finishedChannel chan<- bool
}

var inflights struct {
	sync.RWMutex
	transfers map[string]inflight
}

func init() {
	inflights.transfers = make(map[string]inflight)
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

func getInflight(id string) (result inflight, ok bool) {
	inflights.RLock()
	defer inflights.RUnlock()
	result, ok = inflights.transfers[id]
	return
}

func setInflightDestination(id, destination string) {
	inflights.RLock()
	defer inflights.RUnlock()
	inflights.transfers[id].status.destination = destination
}

func removeInflight(id string) {
	inflights.Lock()
	defer inflights.Unlock()
	delete(inflights.transfers, id)
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

func (ift inflight) finalizeFileTransfer(tempName string) error {
	if err := os.Rename(tempName, ift.status.destination); err != nil {
		ift.reportError(errors.New("Couldn't save final file"))
		return err
	}

	ift.reportFinished()
	removeInflight(ift.id)

	return nil
}

func (ift inflight) openDestinationTempFile() (f *os.File, err error) {
	// By creating a temp file next to the place where the real file should be saved
	// we avoid problems on linux when trying to os.Rename later - if tmp filesystem is different
	// than the destination file system. It also serves as an early permissions check.
	f, err = ioutil.TempFile(filepath.Dir(ift.status.destination), filepath.Base(ift.status.destination))
	if err != nil {
		ift.status.opaque = nil
		ift.reportError(errors.New("Couldn't open local temporary file"))
		removeInflight(ift.id)
	}
	return
}

func (ift inflight) reportError(e error) {
	close(ift.finishedChannel)
	close(ift.updateChannel)
	ift.errorChannel <- e
	close(ift.errorChannel)
}

func (ift inflight) reportFinished() {
	close(ift.errorChannel)
	close(ift.updateChannel)
	ift.finishedChannel <- true
	close(ift.finishedChannel)
}

func waitForFileTransferUserAcceptance(s access.Session, stanza *data.ClientIQ, si data.SI, acceptResult <-chan *string, ift inflight) {
	result := <-acceptResult

	var error *data.ErrorReply
	if result != nil {
		opt, ok := chooseAppropriateFileTransferOptionFrom(ift.options)
		if ok {
			setInflightDestination(si.ID, *result)
			s.SendIQResult(stanza, iqResultChosenStreamMethod(opt))
			go fileTransferCancelListeners[opt](s, ift)
			return
		}
		ift.reportError(errors.New("No mutually acceptable file transfer methods available"))
		error = &iqErrorBadRequest
	} else {
		error = &iqErrorForbidden
	}
	removeInflight(si.ID)
	s.SendIQError(stanza, *error)
}

func registerNewFileTransfer(si data.SI, options []string, stanza *data.ClientIQ, f *data.File, cc <-chan bool, ec chan<- error, uc chan<- int64, fc chan<- bool) inflight {
	ift := inflight{
		id:              si.ID,
		mime:            si.MIMEType,
		options:         options,
		date:            f.Date,
		hash:            f.Hash,
		name:            f.Name,
		size:            f.Size,
		desc:            f.Desc,
		peer:            stanza.From,
		status:          &inflightStatus{},
		cancelChannel:   cc,
		errorChannel:    ec,
		updateChannel:   uc,
		finishedChannel: fc,
	}

	if f.Range != nil {
		ift.rng.length = f.Range.Length
		ift.rng.offset = f.Range.Offset
	}

	inflights.Lock()
	defer inflights.Unlock()
	inflights.transfers[si.ID] = ift
	return ift
}

// InitIQ is the hook function that will be called when we receive a file transfer stream initiation IQ
func InitIQ(s access.Session, stanza *data.ClientIQ, si data.SI) (ret interface{}, iqtype string, ignore bool) {
	var options []string
	var err error
	if options, err = extractFileTransferOptions(si.Feature.Form); err != nil {
		s.Warn(fmt.Sprintf("Failed to parse stream initiation: %v", err))
		return nil, "", false
	}

	f := si.File

	cancelChannel := make(chan bool)
	errorChannel := make(chan error)
	updateChannel := make(chan int64, 1000)
	finishedChannel := make(chan bool)
	ift := registerNewFileTransfer(si, options, stanza, f, cancelChannel, errorChannel, updateChannel, finishedChannel)

	acceptResult := make(chan *string)
	go waitForFileTransferUserAcceptance(s, stanza, si, acceptResult, ift)

	s.PublishEvent(events.FileTransfer{
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

	return nil, "", true
}

package filetransfer

import (
	"fmt"

	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
)

// InitIQ is the hook function that will be called when we receive a file or directory transfer stream initiation IQ
func InitIQ(s access.Session, stanza *data.ClientIQ, si data.SI) (ret interface{}, iqtype string, ignore bool) {
	isDir := si.Profile == dirTransferProfile

	var options []string
	var err error
	if options, err = extractFileTransferOptions(si.Feature.Form); err != nil {
		s.Warn(fmt.Sprintf("Failed to parse stream initiation: %v", err))
		return nil, "", false
	}

	f := si.File

	ctl := sdata.CreateFileTransferControl()
	ctx := registerNewFileTransfer(s, si, options, stanza, f, ctl, isDir)

	acceptResult := make(chan *string)
	go waitForFileTransferUserAcceptance(stanza, si, acceptResult, ctx)

	s.PublishEvent(events.FileTransfer{
		Session:          s,
		Peer:             stanza.From,
		Mime:             f.Hash,
		DateLastModified: f.Date,
		Name:             f.Name,
		Size:             f.Size,
		Description:      f.Desc,
		Answer:           acceptResult,
		Control:          ctl,
		IsDirectory:      isDir,
	})

	return nil, "", true
}

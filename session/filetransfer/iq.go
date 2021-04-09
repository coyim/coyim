package filetransfer

import (
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// InitIQ is the hook function that will be called when we receive a file or directory transfer stream initiation IQ
func InitIQ(s hasLogConnectionIQSymmetricKeyAndIsPublisher, stanza *data.ClientIQ, si data.SI) (ret interface{}, iqtype string, ignore bool) {
	peer, ok := jid.Parse(stanza.From).(jid.WithResource)
	if !ok {
		s.Log().WithField("from", stanza.From).Warn("Stanza sender doesn't contain resource - this shouldn't happen")
		return nil, "", false
	}

	isDir := false
	isEnc := false
	switch si.Profile {
	case dirTransferProfile:
		isDir = true
	case encryptedTransferProfile:
		isEnc = true
		isDir = si.EncryptedData.Type == "directory"
	}

	var options []string
	var err error
	if options, err = extractFileTransferOptions(si.Feature.Form); err != nil {
		s.Log().WithError(err).Warn("Failed to parse stream initiation")
		return nil, "", false
	}

	ctl := sdata.CreateFileTransferControl(nil, nil)

	ctx := registerNewFileTransfer(s, si, options, stanza, ctl, isDir, isEnc)

	acceptResult := make(chan *string)
	go waitForFileTransferUserAcceptance(stanza, si, acceptResult, ctx)

	s.PublishEvent(events.FileTransfer{
		Peer:             peer,
		Mime:             ctx.hash,
		DateLastModified: ctx.date,
		Name:             ctx.name,
		Size:             ctx.size,
		Description:      ctx.desc,
		Answer:           acceptResult,
		Control:          ctl,
		IsDirectory:      isDir,
		Encrypted:        isEnc,
	})

	return nil, "", true
}

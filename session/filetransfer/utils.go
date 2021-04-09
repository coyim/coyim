package filetransfer

import (
	"bytes"
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
)

var errChannelClosed = errors.New("channel closed")
var errNotResultIQ = errors.New("expected result IQ")
var errNotClientIQ = errors.New("expected Client IQ")

type hasConnection interface {
	Conn() xi.Conn
}

type hasConfig interface {
	GetConfig() *config.Account
}

type hasConnectionAndConfig interface {
	hasConnection
	hasConfig
}

type hasLog interface {
	Log() coylog.Logger
}

type hasConnectionAndConfigAndLog interface {
	hasConnectionAndConfig
	hasLog
}

type hasConfigAndLog interface {
	hasConfig
	hasLog
}

type canSendIQError interface {
	SendIQError(*data.ClientIQ, interface{})
}

type canSendIQResult interface {
	SendIQResult(*data.ClientIQ, interface{})
}

type canSendIQ interface {
	canSendIQError
	canSendIQResult
}

type canSendIQAndHasLog interface {
	canSendIQ
	hasLog
}

type canSendIQAndHasLogAndConnection interface {
	canSendIQAndHasLog
	hasConnection
}

type canSendIQErrorAndHasLog interface {
	hasLog
	canSendIQError
}

type canSendIQErrorHasConfigAndHasLog interface {
	hasConfigAndLog
	canSendIQError
}

type hasSymmetricKey interface {
	CreateSymmetricKeyFor(jid.Any) []byte
	GetAndWipeSymmetricKeyFor(jid.Any) []byte
}

type hasConnectionAndConfigAndLogAndHasSymmetricKey interface {
	hasConnectionAndConfigAndLog
	hasSymmetricKey
}

type publisher interface {
	PublishEvent(interface{})
}

type hasLogConnectionIQSymmetricKeyAndIsPublisher interface {
	canSendIQAndHasLogAndConnection
	hasSymmetricKey
	publisher
}

func basicIQ(s hasConnection, to, tp string, toSend, unpackInto interface{}, onSuccess func(*data.ClientIQ)) error {
	done := make(chan error, 1)

	nonblockIQ(s, to, tp, toSend, unpackInto, func(ciq *data.ClientIQ) {
		onSuccess(ciq)
		done <- nil
	}, func(_ *data.ClientIQ, ee error) {
		done <- ee
	})

	return <-done
}

func nonblockIQ(s hasConnection, to, tp string, toSend, unpackInto interface{}, onSuccess func(*data.ClientIQ), onError func(*data.ClientIQ, error)) {
	rp, _, err := s.Conn().SendIQ(to, tp, toSend)
	if err != nil {
		onError(nil, err)
		return
	}
	go func() {
		r, ok := <-rp
		if !ok {
			onError(nil, errChannelClosed)
			return
		}

		switch ciq := r.Value.(type) {
		case *data.ClientIQ:
			if ciq.Type != "result" {
				onError(ciq, errNotResultIQ)
				return
			}
			if unpackInto != nil {
				if err := xml.NewDecoder(bytes.NewBuffer(ciq.Query)).Decode(unpackInto); err != nil {
					onError(ciq, err)
					return
				}
			}

			onSuccess(ciq)
			return
		}
		onError(nil, errNotClientIQ)
	}()
}

package filetransfer

import (
	"bytes"
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
)

var channelClosedError = errors.New("Channel closed")
var notResultIQ = errors.New("Expected result IQ")
var notClientIQ = errors.New("Expected Client IQ")

func basicIQ(s access.Session, to, tp string, toSend, unpackInto interface{}, onSuccess func(*data.ClientIQ)) error {
	done := make(chan error)

	nonblockIQ(s, to, tp, toSend, unpackInto, func(ciq *data.ClientIQ) {
		onSuccess(ciq)
		done <- nil
	}, func(_ *data.ClientIQ, ee error) {
		done <- ee
	})

	return <-done
}

func nonblockIQ(s access.Session, to, tp string, toSend, unpackInto interface{}, onSuccess func(*data.ClientIQ), onError func(*data.ClientIQ, error)) {
	rp, _, err := s.Conn().SendIQ(to, tp, toSend)
	if err != nil {
		onError(nil, err)
		return
	}
	go func() {
		r, ok := <-rp
		if !ok {
			onError(nil, channelClosedError)
			return
		}

		switch ciq := r.Value.(type) {
		case *data.ClientIQ:
			if ciq.Type != "result" {
				onError(ciq, notResultIQ)
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
		onError(nil, notClientIQ)
	}()
}

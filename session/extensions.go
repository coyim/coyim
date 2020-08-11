package session

import (
	"bytes"
	"fmt"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
)

type extensionFunction func(access.Session, *data.ClientMessage, *data.Extension)

var knownExtensions = map[string]extensionFunction{}

func unknownExtension(s access.Session, stanza *data.ClientMessage, ext *data.Extension) {
	s.Log().WithField("extension", bytes.NewBuffer([]byte(ext.Body))).Info("Unknown extension")
}

func registerKnownExtension(fullName string, f extensionFunction) {
	knownExtensions[fullName] = f
}

func getExtensionHandler(namespace, local string) extensionFunction {
	f, ok := knownExtensions[fmt.Sprintf("%s %s", namespace, local)]
	if ok {
		return f
	}
	return unknownExtension
}

func (s *session) processExtension(stanza *data.ClientMessage, ext *data.Extension) {
	if nspace, local, ok := tryDecodeXML([]byte(ext.Body)); ok {
		getExtensionHandler(nspace, local)(s, stanza, ext)
	}
}

func (s *session) processExtensions(stanza *data.ClientMessage) {
	for _, ext := range stanza.Extensions {
		s.processExtension(stanza, ext)
	}
}

package session

import (
	"bytes"
	"encoding/xml"
)

func tryDecodeXML(data []byte) (nspace, local string, ok bool) {
	token, _ := xml.NewDecoder(bytes.NewBuffer(data)).Token()
	if token == nil {
		return "", "", false
	}

	startElem, ok := token.(xml.StartElement)
	if !ok {
		return "", "", false
	}

	return startElem.Name.Space, startElem.Name.Local, true
}

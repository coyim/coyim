package session

import (
	"bytes"
	"fmt"

	"github.com/twstrike/coyim/event"
)

var otrWhitespaceTagStart = []byte("\x20\x09\x20\x20\x09\x09\x09\x09\x20\x09\x20\x09\x20\x09\x20\x20")

var otrWhiteSpaceTagV1 = []byte("\x20\x09\x20\x09\x20\x20\x09\x20")
var otrWhiteSpaceTagV2 = []byte("\x20\x20\x09\x09\x20\x20\x09\x20")
var otrWhiteSpaceTagV3 = []byte("\x20\x20\x09\x09\x20\x20\x09\x09")

var otrWhitespaceTag = append(otrWhitespaceTagStart, otrWhiteSpaceTagV2...)

// TODO: this is taken care of in OTR3, it shouldn't be necessary here
func (s *Session) processWhitespaceTag(encrypted bool, out []byte, from string) {
	detectedOTRVersion := 0
	// We don't need to alert about tags encoded inside of messages that are
	// already encrypted with OTR
	whitespaceTagLength := len(otrWhitespaceTagStart) + len(otrWhiteSpaceTagV1)
	if !encrypted && len(out) >= whitespaceTagLength {
		whitespaceTag := out[len(out)-whitespaceTagLength:]
		if bytes.Equal(whitespaceTag[:len(otrWhitespaceTagStart)], otrWhitespaceTagStart) {
			if bytes.HasSuffix(whitespaceTag, otrWhiteSpaceTagV1) {
				s.info(fmt.Sprintf("%s appears to support OTRv1. You should encourage them to upgrade their OTR client!", from))
				detectedOTRVersion = 1
			}
			if bytes.HasSuffix(whitespaceTag, otrWhiteSpaceTagV2) {
				detectedOTRVersion = 2
			}
			if bytes.HasSuffix(whitespaceTag, otrWhiteSpaceTagV3) {
				detectedOTRVersion = 3
			}
		}
	}

	if s.Config.OTRAutoStartSession && detectedOTRVersion >= 2 {
		s.info(fmt.Sprintf("%s appears to support OTRv%d. We are attempting to start an OTR session with them.", from, detectedOTRVersion))
		s.Conn.Send(from, event.QueryMessage)
	} else if s.Config.OTRAutoStartSession && detectedOTRVersion == 1 {
		s.info(fmt.Sprintf("%s appears to support OTRv%d. You should encourage them to upgrade their OTR client!", from, detectedOTRVersion))
	}
}

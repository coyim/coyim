package session

import (
	"bytes"
	"fmt"

	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/ui"
)

func (s *Session) processWhitespaceTag(encrypted bool, out []byte, from string) {
	detectedOTRVersion := 0
	// We don't need to alert about tags encoded inside of messages that are
	// already encrypted with OTR
	whitespaceTagLength := len(ui.OTRWhitespaceTagStart) + len(ui.OTRWhiteSpaceTagV1)
	if !encrypted && len(out) >= whitespaceTagLength {
		whitespaceTag := out[len(out)-whitespaceTagLength:]
		if bytes.Equal(whitespaceTag[:len(ui.OTRWhitespaceTagStart)], ui.OTRWhitespaceTagStart) {
			if bytes.HasSuffix(whitespaceTag, ui.OTRWhiteSpaceTagV1) {
				s.info(fmt.Sprintf("%s appears to support OTRv1. You should encourage them to upgrade their OTR client!", from))
				detectedOTRVersion = 1
			}
			if bytes.HasSuffix(whitespaceTag, ui.OTRWhiteSpaceTagV2) {
				detectedOTRVersion = 2
			}
			if bytes.HasSuffix(whitespaceTag, ui.OTRWhiteSpaceTagV3) {
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

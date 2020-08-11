package session

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/coylog"
)

//TODO: error
func openLogFile(logFile string) io.Writer {
	if len(logFile) == 0 {
		return nil
	}

	log.WithField("file", logFile).Debug("Logging XMPP messages to file")

	rawLog, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.WithError(err).Warn("Failed to open log file.")
		//return nil, errors.New("Failed to open raw log file: " + err.Error())
		return nil
	}

	return rawLog
}

// Log is the implementation for session interface
func (s *session) Log() coylog.Logger {
	return s.log
}

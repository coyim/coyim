package session

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/session/events"
)

//TODO: error
func openLogFile(logFile string) io.Writer {
	if len(logFile) == 0 {
		return nil
	}

	log.Println("Logging XMPP messages to:", logFile)

	rawLog, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Println("Failed to open log file.", err)
		//return nil, errors.New("Failed to open raw log file: " + err.Error())
		return nil
	}

	return rawLog
}

func (s *session) info(m string) {
	s.publishEvent(events.Log{
		Level:   events.Info,
		Message: m,
	})
}

func (s *session) Warn(m string) {
	s.warn(m)
}

func (s *session) Info(m string) {
	s.info(m)
}

func (s *session) warn(m string) {
	s.publishEvent(events.Log{
		Level:   events.Warn,
		Message: m,
	})
}

func (s *session) alert(m string) {
	s.publishEvent(events.Log{
		Level:   events.Alert,
		Message: m,
	})
}

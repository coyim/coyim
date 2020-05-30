package session

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type logger struct {
	*log.Logger
}

// LogToDebugLog return the log for debugging purposes.
func LogToDebugLog() io.Writer {
	return &logger{}
}

func newLogger() io.Writer {
	l := log.New()
	l.SetOutput(os.Stderr)

	return &logger{l}
}

func (l *logger) Write(m []byte) (int, error) {
	if l.Logger != nil {
		l.Print(string(m))
	} else {
		log.Print(string(m))
	}
	return len(m), nil
}

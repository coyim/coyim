package session

import (
	"io"
	"log"
	"os"
)

type logger struct {
	*log.Logger
}

func LogToDebugLog() io.Writer {
	return &logger{}
}

func newLogger() io.Writer {
	return &logger{
		log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l *logger) Write(m []byte) (int, error) {
	if l.Logger != nil {
		l.Print(string(m))
	} else {
		log.Print(string(m))
	}
	return len(m), nil
}

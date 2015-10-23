package session

import (
	"io"
	"log"
	"os"
)

type logger struct {
	*log.Logger
}

func newLogger() io.Writer {
	return &logger{
		log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l *logger) Write(m []byte) (int, error) {
	l.Print(string(m))
	return len(m), nil
}

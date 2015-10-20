package session

import "log"

type logger struct{}

func (logger) Write(m []byte) (int, error) {
	log.Print(string(m))
	return len(m), nil
}

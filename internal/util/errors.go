package util

import "github.com/sirupsen/logrus"

// LogIgnoredError will log the error given if not nil. It will log to the standard logger if none given.
func LogIgnoredError(err error, log logrus.FieldLogger, msg string) {
    if err != nil {
		if log == nil {
			log = logrus.StandardLogger()
		}
        log.WithError(err).Debug(msg)
    }
}

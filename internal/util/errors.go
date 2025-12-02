package util

import "github.com/sirupsen/logrus"

func LogIgnoredError(err error, log logrus.FieldLogger, msg string) {
    if err != nil {
		if log == nil {
			log = logrus.StandardLogger()
		}
        log.WithError(err).Debug(msg)
    }
}

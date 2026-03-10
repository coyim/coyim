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

// OrErr will log the error if it is not nil, returning the value given
func OrErr[T any](v T, err error, log logrus.FieldLogger, msg string) T {
	if err != nil {
		if log == nil {
			log = logrus.StandardLogger()
		}
		log.WithError(err).Debug(msg)
	}
	return v
}

// Must will panic if an error is given, otherwise returning the value
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err) // or log.Fatal(err)
	}
	return v
}

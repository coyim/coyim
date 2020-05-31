package coylog

import (
	"context"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

// Logger represents any of the *logrus.Logger and *logrus.Entry types
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Log(level log.Level, args ...interface{})
	Logf(level log.Level, format string, args ...interface{})
	Logln(level log.Level, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Traceln(args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Warningln(args ...interface{})
	Warnln(args ...interface{})
	WithContext(ctx context.Context) *log.Entry
	WithError(err error) *log.Entry
	WithField(key string, value interface{}) *log.Entry
	WithFields(fields log.Fields) *log.Entry
	WithTime(t time.Time) *log.Entry
	Writer() *io.PipeWriter
	WriterLevel(level log.Level) *io.PipeWriter
}

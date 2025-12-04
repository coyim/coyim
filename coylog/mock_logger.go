package coylog

import (
	"context"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of the Logger interface
type MockLogger struct {
	mock.Mock
}

// Debug implements the Logger interface
func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

// Debugf implements the Logger interface
func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

// Debugln implements the Logger interface
func (m *MockLogger) Debugln(args ...interface{}) {
	m.Called(args...)
}

// Error implements the Logger interface
func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

// Errorf implements the Logger interface
func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

// Errorln implements the Logger interface
func (m *MockLogger) Errorln(args ...interface{}) {
	m.Called(args...)
}

// Fatal implements the Logger interface
func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

// Fatalf implements the Logger interface
func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args)
}

// Fatalln implements the Logger interface
func (m *MockLogger) Fatalln(args ...interface{}) {
	m.Called(args...)
}

// Info implements the Logger interface
func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

// Infof implements the Logger interface
func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

// Infoln implements the Logger interface
func (m *MockLogger) Infoln(args ...interface{}) {
	m.Called(args...)
}

// Log implements the Logger interface
func (m *MockLogger) Log(level log.Level, args ...interface{}) {
	m.Called(level, args)
}

// Logf implements the Logger interface
func (m *MockLogger) Logf(level log.Level, format string, args ...interface{}) {
	m.Called(level, format, args)
}

// Logln implements the Logger interface
func (m *MockLogger) Logln(level log.Level, args ...interface{}) {
	m.Called(level, args)
}

// Panic implements the Logger interface
func (m *MockLogger) Panic(args ...interface{}) {
	m.Called(args...)
}

// Panicf implements the Logger interface
func (m *MockLogger) Panicf(format string, args ...interface{}) {
	m.Called(format, args)
}

// Panicln implements the Logger interface
func (m *MockLogger) Panicln(args ...interface{}) {
	m.Called(args...)
}

// Print implements the Logger interface
func (m *MockLogger) Print(args ...interface{}) {
	m.Called(args...)
}

// Printf implements the Logger interface
func (m *MockLogger) Printf(format string, args ...interface{}) {
	m.Called(format, args)
}

// Println implements the Logger interface
func (m *MockLogger) Println(args ...interface{}) {
	m.Called(args...)
}

// Trace implements the Logger interface
func (m *MockLogger) Trace(args ...interface{}) {
	m.Called(args...)
}

// Tracef implements the Logger interface
func (m *MockLogger) Tracef(format string, args ...interface{}) {
	m.Called(format, args)
}

// Traceln implements the Logger interface
func (m *MockLogger) Traceln(args ...interface{}) {
	m.Called(args...)
}

// Warn implements the Logger interface
func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args...)
}

// Warnf implements the Logger interface
func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

// Warning implements the Logger interface
func (m *MockLogger) Warning(args ...interface{}) {
	m.Called(args...)
}

// Warningf implements the Logger interface
func (m *MockLogger) Warningf(format string, args ...interface{}) {
	m.Called(format, args)
}

// Warningln implements the Logger interface
func (m *MockLogger) Warningln(args ...interface{}) {
	m.Called(args...)
}

// Warnln implements the Logger interface
func (m *MockLogger) Warnln(args ...interface{}) {
	m.Called(args...)
}

// WithContext implements the Logger interface
func (m *MockLogger) WithContext(ctx context.Context) *log.Entry {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*log.Entry)
}

// WithError implements the Logger interface
func (m *MockLogger) WithError(err error) *log.Entry {
	args := m.Called(err)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*log.Entry)
}

// WithField implements the Logger interface
func (m *MockLogger) WithField(key string, value interface{}) *log.Entry {
	args := m.Called(key, value)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*log.Entry)
}

// WithFields implements the Logger interface
func (m *MockLogger) WithFields(fields log.Fields) *log.Entry {
	args := m.Called(fields)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*log.Entry)
}

// WithTime implements the Logger interface
func (m *MockLogger) WithTime(t time.Time) *log.Entry {
	args := m.Called(t)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*log.Entry)
}

// Writer implements the Logger interface
func (m *MockLogger) Writer() *io.PipeWriter {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*io.PipeWriter)
}

// WriterLevel implements the Logger interface
func (m *MockLogger) WriterLevel(level log.Level) *io.PipeWriter {
	args := m.Called(level)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*io.PipeWriter)
}

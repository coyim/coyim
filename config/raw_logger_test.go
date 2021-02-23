package config

import (
	"errors"
	"strings"
	"sync"

	. "gopkg.in/check.v1"
)

type RawLoggerSuite struct{}

var _ = Suite(&RawLoggerSuite{})

func (s *RawLoggerSuite) Test_rawLogger_flush_emptyBuffer(c *C) {
	r := &rawLogger{}
	e := r.flush()
	c.Assert(e, IsNil)
}

type multipleResultsWriter struct {
	args         [][]byte
	resultsLen   []int
	resultsError []error
}

func (mew *multipleResultsWriter) Write(p []byte) (int, error) {
	mew.args = append(mew.args, p)
	rl, re := mew.resultsLen[0], mew.resultsError[0]
	mew.resultsLen = mew.resultsLen[1:]
	mew.resultsError = mew.resultsError[1:]
	return rl, re
}

func (s *RawLoggerSuite) Test_rawLogger_flush_simpleOutput(c *C) {
	b := strings.Builder{}
	r := &rawLogger{buf: []byte("Something simple"), prefix: []byte("a prefix"), out: &b}
	e := r.flush()
	c.Assert(e, IsNil)
	c.Assert(b.String(), Equals, "a prefixSomething simple\n")
	c.Assert(r.buf, DeepEquals, []byte{})
}

func (s *RawLoggerSuite) Test_rawLogger_flush_firstError(c *C) {
	e1 := errors.New("first error")
	b := &multipleResultsWriter{
		resultsLen:   []int{0, 0, 0},
		resultsError: []error{e1, nil, nil},
	}
	r := &rawLogger{buf: []byte("Something simple"), prefix: []byte("a prefix"), out: b}
	e := r.flush()
	c.Assert(e, Equals, e1)
	c.Assert(b.args, DeepEquals, [][]byte{[]byte("a prefix")})
	c.Assert(r.buf, DeepEquals, []byte("Something simple"))
}

func (s *RawLoggerSuite) Test_rawLogger_flush_secondError(c *C) {
	e1 := errors.New("first error")
	b := &multipleResultsWriter{
		resultsLen:   []int{0, 0, 0},
		resultsError: []error{nil, e1, nil},
	}
	r := &rawLogger{buf: []byte("Something simple"), prefix: []byte("a prefix"), out: b}
	e := r.flush()
	c.Assert(e, Equals, e1)
	c.Assert(b.args, DeepEquals, [][]byte{[]byte("a prefix"), []byte("Something simple")})
	c.Assert(r.buf, DeepEquals, []byte("Something simple"))
}

func (s *RawLoggerSuite) Test_rawLogger_flush_thirdError(c *C) {
	e1 := errors.New("first error")
	b := &multipleResultsWriter{
		resultsLen:   []int{0, 0, 0},
		resultsError: []error{nil, nil, e1},
	}
	r := &rawLogger{buf: []byte("Something simple"), prefix: []byte("a prefix"), out: b}
	e := r.flush()
	c.Assert(e, Equals, e1)
	c.Assert(b.args, DeepEquals, [][]byte{[]byte("a prefix"), []byte("Something simple"), []byte("\n")})
	c.Assert(r.buf, DeepEquals, []byte("Something simple"))
}

func (s *RawLoggerSuite) Test_rawLogger_Write_returnsEarlyIfFlushingOtherFails(c *C) {
	lock := new(sync.Mutex)
	e1 := errors.New("first error")
	b := &multipleResultsWriter{
		resultsLen:   []int{0, 0, 0},
		resultsError: []error{e1, nil, nil},
	}
	other := &rawLogger{buf: []byte("something"), prefix: nil, lock: lock, out: b}
	main := &rawLogger{buf: []byte{}, prefix: nil, lock: lock, other: other}
	other.other = main

	num, e := main.Write([]byte{})
	c.Assert(num, Equals, 0)
	c.Assert(e, IsNil)
}

func (s *RawLoggerSuite) Test_rawLogger_Write_storesTheDataSubmittedAndReturnsTheLength(c *C) {
	lock := new(sync.Mutex)
	b := &multipleResultsWriter{
		resultsLen:   []int{0, 0, 0},
		resultsError: []error{nil, nil, nil},
	}
	other := &rawLogger{buf: []byte("something"), prefix: nil, lock: lock, out: b}
	main := &rawLogger{buf: []byte{}, prefix: nil, lock: lock, other: other}
	other.other = main

	num, e := main.Write([]byte("some data\nwith new \nlines"))
	c.Assert(num, Equals, 25)
	c.Assert(e, IsNil)
	c.Assert(string(main.buf), Equals, "some datawith new lines")
}

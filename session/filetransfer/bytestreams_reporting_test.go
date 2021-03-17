package filetransfer

import (
	"errors"

	. "gopkg.in/check.v1"
)

type BytestreamsReportingSuite struct{}

var _ = Suite(&BytestreamsReportingSuite{})

func (s *BytestreamsReportingSuite) Test_reportingWriter_Write_reports(c *C) {
	calledWith := 0

	rr := &reportingWriter{
		report: func(v int) error {
			calledWith = v
			return errors.New("marker one")
		},
	}

	n, e := rr.Write([]byte{1, 2, 3})
	c.Assert(n, Equals, 3)
	c.Assert(e, ErrorMatches, "marker one")
	c.Assert(calledWith, Equals, 3)
}

type mockReadCloser struct {
	fr func([]byte) (int, error)
	fc func() error
}

func (m *mockReadCloser) Read(p []byte) (int, error) {
	return m.fr(p)
}

func (m *mockReadCloser) Close() error {
	return m.fc()
}

func (s *BytestreamsReportingSuite) Test_reportingReader_Read_reports(c *C) {
	m := &mockReadCloser{}
	m.fr = func([]byte) (int, error) {
		return 55, errors.New("final marker")
	}

	calledWith := 0
	rr := &reportingReader{}
	rr.r = m
	rr.report = func(v int) error {
		calledWith = v
		return errors.New("another marker")
	}

	n, e := rr.Read(nil)
	c.Assert(calledWith, Equals, 55)
	c.Assert(n, Equals, 55)
	c.Assert(e, ErrorMatches, "final marker")
}

func (s *BytestreamsReportingSuite) Test_reportingReader_Close_closes(c *C) {
	m := &mockReadCloser{}
	called := false
	m.fc = func() error {
		called = true
		return errors.New("one final marker")
	}

	rr := &reportingReader{}
	rr.r = m

	e := rr.Close()
	c.Assert(called, Equals, true)
	c.Assert(e, ErrorMatches, "one final marker")
}

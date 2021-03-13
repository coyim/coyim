package session

import (
	"time"

	. "gopkg.in/check.v1"
)

type ReconnectSuite struct{}

var _ = Suite(&ReconnectSuite{})

func (s *ReconnectSuite) Test_calculateRandomDelay_works(c *C) {
	del := calculateRandomDelay()
	c.Assert(del > time.Duration(0), Equals, true)
	c.Assert(del < time.Duration(18)*time.Second, Equals, true)
}

func (s *ReconnectSuite) Test_randomDelayChannel_works(c *C) {
	cc := randomDelayChannel()
	c.Assert(cc, Not(IsNil))
}

func (s *ReconnectSuite) Test_checkReconnect_whileConnected(c *C) {
	orgRandomDelayChannel := randomDelayChannel
	defer func() {
		randomDelayChannel = orgRandomDelayChannel
	}()

	cc := make(chan time.Time, 1)
	calledNum := 0
	done := make(chan bool)
	randomDelayChannel = func() <-chan time.Time {
		calledNum++
		if calledNum == 2 {
			done <- true
		}
		return cc
	}

	called := false

	mc := &mockConnector{
		connect: func() {
			called = true
		},
	}

	sess := &session{
		connector: mc,
	}

	cc <- time.Time{}
	go checkReconnect(sess)
	close(cc)
	<-done
	c.Assert(calledNum, Equals, 2)
	c.Assert(called, Equals, false)
}

func (s *ReconnectSuite) Test_checkReconnect_whileDisconnected(c *C) {
	orgRandomDelayChannel := randomDelayChannel
	defer func() {
		randomDelayChannel = orgRandomDelayChannel
	}()

	cc := make(chan time.Time, 1)
	calledNum := 0
	done := make(chan bool)
	randomDelayChannel = func() <-chan time.Time {
		calledNum++
		if calledNum == 2 {
			done <- true
		}
		return cc
	}

	called := false

	mc := &mockConnector{
		connect: func() {
			called = true
		},
	}

	sess := &session{
		connStatus: DISCONNECTED,
		connector:  mc,
	}

	cc <- time.Time{}
	go checkReconnect(sess)
	close(cc)
	<-done
	c.Assert(calledNum, Equals, 2)
	c.Assert(called, Equals, false)
}

func (s *ReconnectSuite) Test_checkReconnect_whileDisconnectedAndWantToBeOnline(c *C) {
	orgRandomDelayChannel := randomDelayChannel
	defer func() {
		randomDelayChannel = orgRandomDelayChannel
	}()

	cc := make(chan time.Time, 1)
	calledNum := 0
	done := make(chan bool)
	randomDelayChannel = func() <-chan time.Time {
		calledNum++
		if calledNum == 2 {
			done <- true
		}
		return cc
	}

	called := false

	mc := &mockConnector{
		connect: func() {
			called = true
		},
	}

	sess := &session{
		connStatus:     DISCONNECTED,
		wantToBeOnline: true,
		connector:      mc,
	}

	cc <- time.Time{}
	go checkReconnect(sess)
	close(cc)
	<-done
	c.Assert(calledNum, Equals, 2)
	c.Assert(called, Equals, true)
}

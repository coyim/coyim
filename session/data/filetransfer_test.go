package data

import (
	"errors"
	"sync"

	. "gopkg.in/check.v1"
)

type FiletransferSuite struct{}

var _ = Suite(&FiletransferSuite{})

func (s *FiletransferSuite) Test_newFileTransferControl_createsANewControl(c *C) {
	c1 := make(chan bool)
	c2 := make(chan error)
	c3 := make(chan transferUpdate)
	c4 := make(chan bool)
	v := newFileTransferControl(c1, c2, c3, c4)

	c.Assert(v.cancelTransfer, Equals, c1)
	c.Assert(v.errorOccurred, Equals, c2)
	c.Assert(v.update, Equals, c3)
	c.Assert(v.transferFinished, Equals, c4)
}

func (s *FiletransferSuite) Test_CreateFileTransferControl_createsAControl(c *C) {
	c1Called := false
	c2Called := false
	c1 := func() bool {
		c1Called = true
		return false
	}
	c2 := func(bool) {
		c2Called = true
	}
	v := CreateFileTransferControl(c1, c2)

	c.Assert(v.cancelTransfer, Not(IsNil))

	_ = v.OnEncryptionNotSupported()
	c.Assert(c1Called, Equals, true)

	v.EncryptionDecision(false)
	c.Assert(c2Called, Equals, true)
}

func (s *FiletransferSuite) Test_closeCancel_doesNothingIfNoChannelExists(c *C) {
	ctl := &FileTransferControl{}
	ctl.closeCancel()
	c.Assert(ctl.cancelTransfer, IsNil)
}

func (s *FiletransferSuite) Test_closeCancel_closesChannelIfExists(c *C) {
	c1 := make(chan bool)
	ctl := &FileTransferControl{cancelTransfer: c1}
	ctl.closeCancel()
	c.Assert(ctl.cancelTransfer, IsNil)
	_, notClosed := <-c1
	c.Assert(notClosed, Equals, false)
}

func (s *FiletransferSuite) Test_WaitForFinish_doesNotCallFunctionIfChannelClosed(c *C) {
	c1 := make(chan bool)
	called := false
	k := func(v bool) {
		called = true
	}

	ctl := &FileTransferControl{transferFinished: c1}
	close(c1)
	ctl.WaitForFinish(k)

	c.Assert(called, Equals, false)
}

func (s *FiletransferSuite) Test_WaitForFinish_callsFunctionWithValueFromChannel(c *C) {
	c1 := make(chan bool)
	cancel := make(chan bool)
	called := false
	kValue := false
	k := func(v bool) {
		kValue = v
		called = true
	}

	ctl := &FileTransferControl{transferFinished: c1, cancelTransfer: cancel}
	go func() {
		c1 <- true
	}()
	ctl.WaitForFinish(k)

	c.Assert(called, Equals, true)
	c.Assert(kValue, Equals, true)
	c.Assert(ctl.cancelTransfer, IsNil)
}

func (s *FiletransferSuite) Test_WaitForError_doesNotCallFunctionIfChannelClosed(c *C) {
	c1 := make(chan error)
	called := false
	k := func(error) {
		called = true
	}

	ctl := &FileTransferControl{errorOccurred: c1}
	close(c1)
	ctl.WaitForError(k)

	c.Assert(called, Equals, false)
}

func (s *FiletransferSuite) Test_WaitForError_callsFunctionWithValueFromChannel(c *C) {
	c1 := make(chan error)
	cancel := make(chan bool)
	called := false
	var kValue error
	k := func(v error) {
		kValue = v
		called = true
	}

	ctl := &FileTransferControl{errorOccurred: c1, cancelTransfer: cancel}
	go func() {
		c1 <- errors.New("hello")
	}()
	ctl.WaitForError(k)

	c.Assert(called, Equals, true)
	c.Assert(kValue, ErrorMatches, "hello")
	c.Assert(ctl.cancelTransfer, IsNil)
}

func (s *FiletransferSuite) Test_WaitForCancel_doesntCallFunctionIfChannelClosed(c *C) {
	c1 := make(chan bool)
	called := false
	k := func() {
		called = true
	}

	ctl := &FileTransferControl{cancelTransfer: c1}
	close(c1)
	ctl.WaitForCancel(k)

	c.Assert(called, Equals, false)
}

func (s *FiletransferSuite) Test_WaitForCancel_doesntCallFunctionIfValueGivenIsFalse(c *C) {
	c1 := make(chan bool)
	called := false
	k := func() {
		called = true
	}

	ctl := &FileTransferControl{cancelTransfer: c1}
	go func() {
		c1 <- false
	}()

	ctl.WaitForCancel(k)

	c.Assert(called, Equals, false)
}

func (s *FiletransferSuite) Test_WaitForCancel_callsFunctionWhenAskedToCancel(c *C) {
	c1 := make(chan bool)
	called := false
	k := func() {
		called = true
	}

	ctl := &FileTransferControl{cancelTransfer: c1}
	go func() {
		c1 <- true
	}()

	ctl.WaitForCancel(k)

	c.Assert(called, Equals, true)
}

func (s *FiletransferSuite) Test_WaitForUpdate_callsFunctionWithEachUpdate(c *C) {
	c1 := make(chan transferUpdate)
	arg1 := []int64{}
	arg2 := []int64{}

	k := func(a1, a2 int64) {
		arg1 = append(arg1, a1)
		arg2 = append(arg2, a2)
	}

	ctl := &FileTransferControl{update: c1}
	go func() {
		c1 <- transferUpdate{current: 0, total: 55}
		c1 <- transferUpdate{current: 3, total: 10}
		c1 <- transferUpdate{current: 55, total: 1423}
		close(c1)
	}()

	ctl.WaitForUpdate(k)

	c.Assert(arg1, DeepEquals, []int64{0, 3, 55})
	c.Assert(arg2, DeepEquals, []int64{55, 10, 1423})
}

func (s *FiletransferSuite) Test_Cancel_doesntDoAnythingIfNoChannelIsOpen(c *C) {
	ctl := &FileTransferControl{}
	ctl.Cancel()
	c.Assert(ctl.cancelTransfer, IsNil)
}

func (s *FiletransferSuite) Test_Cancel_sendsCancelRequestOnChannel(c *C) {
	c1 := make(chan bool)
	ctl := &FileTransferControl{cancelTransfer: c1}

	valueReceived := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		valueReceived = <-c1
		wg.Done()
	}()

	ctl.Cancel()
	wg.Wait()
	c.Assert(valueReceived, Equals, true)
}

func (s *FiletransferSuite) Test_ReportError_closesChannelsAndSendsError(c *C) {
	c1 := make(chan bool)
	c2 := make(chan transferUpdate)
	c3 := make(chan error)

	ctl := &FileTransferControl{transferFinished: c1, update: c2, errorOccurred: c3}

	var valueReceived error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		valueReceived = <-c3
		wg.Done()
	}()

	ctl.ReportError(errors.New("oopsie daysie"))
	wg.Wait()

	c.Assert(ctl.transferFinished, IsNil)
	c.Assert(ctl.update, IsNil)
	c.Assert(ctl.errorOccurred, IsNil)
	c.Assert(valueReceived, ErrorMatches, "oopsie .*")
}

func (s *FiletransferSuite) Test_ReportErrorNonblocking_closesChannelsAndSendsError(c *C) {
	c1 := make(chan bool)
	c2 := make(chan transferUpdate)
	c3 := make(chan error)

	ctl := &FileTransferControl{transferFinished: c1, update: c2, errorOccurred: c3}

	var valueReceived error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		valueReceived = <-c3
		_, _ = <-c3
		wg.Done()
	}()

	ctl.ReportErrorNonblocking(errors.New("oopsie daysie"))
	wg.Wait()

	c.Assert(ctl.transferFinished, IsNil)
	c.Assert(ctl.update, IsNil)
	c.Assert(ctl.errorOccurred, IsNil)
	c.Assert(valueReceived, ErrorMatches, "oopsie .*")
}

func (s *FiletransferSuite) Test_ReportFinished_closesChannelsAndSendsResult(c *C) {
	tu := make(chan transferUpdate)
	eo := make(chan error)
	tf := make(chan bool)

	ctl := &FileTransferControl{transferFinished: tf, update: tu, errorOccurred: eo}

	var valueReceived bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		valueReceived = <-tf
		wg.Done()
	}()

	ctl.ReportFinished()
	wg.Wait()

	c.Assert(ctl.transferFinished, IsNil)
	c.Assert(ctl.update, IsNil)
	c.Assert(ctl.errorOccurred, IsNil)
	c.Assert(valueReceived, Equals, true)
}

func (s *FiletransferSuite) Test_ReportDeclined_closesChannelsAndSendsResult(c *C) {
	tu := make(chan transferUpdate)
	eo := make(chan error)
	tf := make(chan bool)

	ctl := &FileTransferControl{transferFinished: tf, update: tu, errorOccurred: eo}

	var valueReceived bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		valueReceived = <-tf
		wg.Done()
	}()

	ctl.ReportDeclined()
	wg.Wait()

	c.Assert(ctl.transferFinished, IsNil)
	c.Assert(ctl.update, IsNil)
	c.Assert(ctl.errorOccurred, IsNil)
	c.Assert(valueReceived, Equals, false)
}

func (s *FiletransferSuite) Test_SendUpdate_doesntDoAnythingIfNoChannelExists(c *C) {
	ctl := &FileTransferControl{}
	ctl.SendUpdate(0, 1)
	c.Assert(ctl.update, IsNil)
}

func (s *FiletransferSuite) Test_SendUpdate_sendsUpdates(c *C) {
	c1 := make(chan transferUpdate)
	ctl := &FileTransferControl{update: c1}

	resCurrent := []int64{}
	resTotal := []int64{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		v := <-c1
		resCurrent = append(resCurrent, v.current)
		resTotal = append(resTotal, v.total)

		v = <-c1
		resCurrent = append(resCurrent, v.current)
		resTotal = append(resTotal, v.total)
		wg.Done()
	}()

	ctl.SendUpdate(0, 1)
	ctl.SendUpdate(55, 1424)
	wg.Wait()
	c.Assert(resCurrent, DeepEquals, []int64{0, 55})
	c.Assert(resTotal, DeepEquals, []int64{1, 1424})
}

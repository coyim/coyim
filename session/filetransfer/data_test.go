package filetransfer

import (
	. "gopkg.in/check.v1"
)

type DataSuite struct{}

var _ = Suite(&DataSuite{})

func (s *DataSuite) Test_inflightRecv(c *C) {
	defer func() {
		removeInflightRecv("foo")
	}()

	vv := &recvContext{
		sid: "foo",
	}

	addInflightRecv(vv)

	setInflightRecvDestination("foo", "somewhere")

	val, ok := getInflightRecv("foo")

	c.Assert(ok, Equals, true)
	c.Assert(val, Equals, vv)
	c.Assert(vv.destination, Equals, "somewhere")
}

func (s *DataSuite) Test_inflightSend(c *C) {
	defer func() {
		removeInflightSend(&sendContext{sid: "baz"})
	}()

	vv := &sendContext{
		sid: "baz",
	}

	addInflightSend(vv)

	val, ok := getInflightSend("baz")

	c.Assert(ok, Equals, true)
	c.Assert(val, Equals, vv)
}

func (s *DataSuite) Test_inflightMAC(c *C) {
	defer func() {
		delete(inflightMACs.transfers, "quux")
	}()

	vv := &sendContext{
		sid: "quux",
	}

	addInflightMAC(vv)

	res := hasAndRemoveInflightMAC("quux")
	c.Assert(res, Equals, true)

	res = hasAndRemoveInflightMAC("quux")
	c.Assert(res, Equals, false)
}

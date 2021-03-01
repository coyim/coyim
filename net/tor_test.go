package net

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	. "gopkg.in/check.v1"
)

type TorSuite struct{}

var _ = Suite(&TorSuite{})

func (s *TorSuite) TestDetectTor(c *C) {
	host := "127.0.0.1"
	if isTails() {
		host = getLocalIP()
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:0", host))
	c.Assert(err, IsNil)

	_, port, err := net.SplitHostPort(ln.Addr().String())
	c.Assert(err, IsNil)

	tor := &defaultTorManager{
		torPorts: []string{port},
		torHost:  host,
	}

	torAddress := ln.Addr().String()
	c.Assert(tor.Address(), Equals, torAddress)

	_ = ln.Close()

	c.Assert(tor.Address(), Equals, torAddress)

	c.Assert(tor.Detect(), Equals, false)
	c.Assert(tor.Address(), Equals, "")
}

func (s *TorSuite) TestDetectTorConnectionRefused(c *C) {
	host := "127.0.0.1"
	if isTails() {
		host = getLocalIP()
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:0", host))
	c.Assert(err, IsNil)

	_, port, err := net.SplitHostPort(ln.Addr().String())
	c.Assert(err, IsNil)

	_ = ln.Close()

	tor := &defaultTorManager{
		torPorts: []string{port},
		torHost:  host,
	}

	c.Assert(tor.Detect(), Equals, false)
	c.Assert(tor.Address(), Equals, "")
}

func (s *TorSuite) Test_defaultTorManager_Detect_usesDefaultsIfNoneGiven(c *C) {
	origDetectTor := detectTor
	defer func() {
		detectTor = origDetectTor
	}()

	calledHost := ""
	calledPorts := []string{}

	detectTor = func(host string, ports []string) (string, bool) {
		calledHost = host
		calledPorts = ports
		return "", false
	}

	df := &defaultTorManager{}

	_ = df.Detect()
	c.Assert(calledHost, Equals, defaultTorHost)
	c.Assert(calledPorts, DeepEquals, defaultTorPorts)
}

func (s *TorSuite) Test_defaultTorManager_IsConnectionOverTor_works(c *C) {
	origHttpGet := httpGet
	defer func() {
		httpGet = origHttpGet
	}()

	called := false
	httpGet = func(*http.Client, string) (*http.Response, error) {
		called = true
		resp := &http.Response{}
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(`{"IsTor": true, "IP": "1.2.3.4"}`))
		return resp, nil
	}

	res := (&defaultTorManager{}).IsConnectionOverTor(nil)
	c.Assert(res, Equals, true)
	c.Assert(called, Equals, true)
}

func (s *TorSuite) Test_httpGet_works(c *C) {
	cl := &http.Client{Transport: &http.Transport{Dial: func(network, addr string) (net.Conn, error) {
		return nil, errors.New("sorry")
	}}}

	res, e := httpGet(cl, "http://hello.com")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "Get .* sorry")
}

func (s *TorSuite) Test_defaultTorManager_IsConnectionOverTor_failsIfDialerFails(c *C) {
	dialer := &mockDialer{
		returnErr: errors.New("oops"),
	}

	res := (&defaultTorManager{}).IsConnectionOverTor(dialer)
	c.Assert(res, Equals, false)
}

type errorReader struct {
	e error
}

func (r *errorReader) Read([]byte) (int, error) {
	return 0, r.e
}

func (s *TorSuite) Test_defaultTorManager_IsConnectionOverTor_failsWhenReadingBody(c *C) {
	origHttpGet := httpGet
	defer func() {
		httpGet = origHttpGet
	}()

	httpGet = func(*http.Client, string) (*http.Response, error) {
		resp := &http.Response{}
		resp.Body = ioutil.NopCloser(&errorReader{errors.New("feeh")})
		return resp, nil
	}

	res := (&defaultTorManager{}).IsConnectionOverTor(nil)
	c.Assert(res, Equals, false)
}

func (s *TorSuite) Test_defaultTorManager_IsConnectionOverTor_failsWhenGivenBadJSONResponse(c *C) {
	origHttpGet := httpGet
	defer func() {
		httpGet = origHttpGet
	}()

	httpGet = func(*http.Client, string) (*http.Response, error) {
		resp := &http.Response{}
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(`{"IsTor": true, "I`))
		return resp, nil
	}

	res := (&defaultTorManager{}).IsConnectionOverTor(nil)
	c.Assert(res, Equals, false)
}

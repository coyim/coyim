// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import "github.com/coyim/coyim/xmpp/data"

// inflight contains the details of a pending request to which we are awaiting
// a reply.
type inflight struct {
	// replyChan is the channel to which we'll send the reply.
	replyChan chan<- data.Stanza
	// to is the address to which we sent the request.
	to string
}

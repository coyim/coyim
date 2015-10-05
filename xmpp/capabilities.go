// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

type xep0115Less interface {
	xep0115Less(interface{}) bool
}

type xep0115Sorter struct{ s []xep0115Less }

func (s *xep0115Sorter) add(c xep0115Less)  { s.s = append(s.s, c) }
func (s *xep0115Sorter) Len() int           { return len(s.s) }
func (s *xep0115Sorter) Swap(i, j int)      { s.s[i], s.s[j] = s.s[j], s.s[i] }
func (s *xep0115Sorter) Less(i, j int) bool { return s.s[i].xep0115Less(s.s[j]) }

func (a *DiscoveryIdentity) xep0115Less(other interface{}) bool {
	b := other.(*DiscoveryIdentity)
	if a.Category != b.Category {
		return a.Category < b.Category
	}
	if a.Type != b.Type {
		return a.Type < b.Type
	}
	return a.Lang < b.Lang
}

func (a *DiscoveryFeature) xep0115Less(other interface{}) bool {
	b := other.(*DiscoveryFeature)
	return a.Var < b.Var
}

func (a *formField) xep0115Less(other interface{}) bool {
	b := other.(*formField)
	if a.Var == "FORM_TYPE" {
		return true
	} else if b.Var == "FORM_TYPE" {
		return false
	}
	return a.Var < b.Var
}

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import "github.com/twstrike/coyim/xmpp/data"

type xep0115Sorter struct{ s []interface{} }

func (s *xep0115Sorter) add(c interface{})  { s.s = append(s.s, c) }
func (s *xep0115Sorter) Len() int           { return len(s.s) }
func (s *xep0115Sorter) Swap(i, j int)      { s.s[i], s.s[j] = s.s[j], s.s[i] }
func (s *xep0115Sorter) Less(i, j int) bool { return xep0115Less(s.s[i], s.s[j]) }

func xep0115Less(a interface{}, other interface{}) bool {
	switch v := a.(type) {
	case *data.DiscoveryIdentity:
		return xep0115LessDI(v, other)
	case *data.DiscoveryFeature:
		return xep0115LessDF(v, other)
	case *data.FormFieldX:
		return xep0115LessFF(v, other)
	}
	return false
}

func xep0115LessDI(a *data.DiscoveryIdentity, other interface{}) bool {
	b := other.(*data.DiscoveryIdentity)
	if a.Category != b.Category {
		return a.Category < b.Category
	}
	if a.Type != b.Type {
		return a.Type < b.Type
	}
	return a.Lang < b.Lang
}

func xep0115LessDF(a *data.DiscoveryFeature, other interface{}) bool {
	b := other.(*data.DiscoveryFeature)
	return a.Var < b.Var
}

func xep0115LessFF(a *data.FormFieldX, other interface{}) bool {
	b := other.(*data.FormFieldX)
	if a.Var == "FORM_TYPE" {
		return true
	} else if b.Var == "FORM_TYPE" {
		return false
	}
	return a.Var < b.Var
}

// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type Tuple []interface{}

func NewTupleFromSlice(slice []interface{}) *Tuple {
	t := Tuple(slice)
	return &t
}

func (t *Tuple) Get(i int) interface{} {
	return (*t)[i]
}

func (t *Tuple) Len() int {
	return len(*t)
}

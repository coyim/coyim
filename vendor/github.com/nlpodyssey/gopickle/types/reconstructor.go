// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "fmt"

type Reconstructor struct{}

var _ Callable = &Reconstructor{}

func (r *Reconstructor) Call(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("Reconstructor: invalid arguments: %#v", args)
	}
	class := args[0]
	switch base := args[1].(type) {
	case PyNewable:
		return base.PyNew(class)
	default:
		return nil, fmt.Errorf(
			"Reconstructor: unprocessable arguments: %#v", args)
	}
}

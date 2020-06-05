// +build darwin

package otr3

import "math/big"

// Not supported on OS X at the moment

func mlockWordSlice(b []big.Word) (err error) {
	return nil
}

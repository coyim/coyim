// +build windows

package otr3

import (
	"fmt"
	"math/big"
	"math/bits"

	"golang.org/x/sys/windows"
)

func mlockWordSlice(b []big.Word) (err error) {
	/* #nosec G103*/
	if err := windows.VirtualLock(_getPtr(b), uintptr(len(b)*(bits.UintSize/8))); err != nil {
		return fmt.Errorf("<memcall> could not acquire lock on %p, limit reached? [Err: %s]", &b[0], err)
	}

	return nil
}

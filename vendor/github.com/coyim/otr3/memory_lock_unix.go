// +build !windows,!darwin

package otr3

import (
	"fmt"
	"math/big"
	"math/bits"
	"syscall"
	"unsafe"
)

func mlockWordSlice(b []big.Word) (err error) {
	if len(b) > 0 {
		/* #nosec G103*/
		_p0 := unsafe.Pointer(&b[0])
		_, _, e1 := syscall.Syscall(syscall.SYS_MLOCK, uintptr(_p0), uintptr(len(b)*(bits.UintSize/8)), 0)
		if e1 != 0 {
			return fmt.Errorf("got error %d", e1)
		}
	}
	return
}

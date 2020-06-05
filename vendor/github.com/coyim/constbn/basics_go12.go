// +build go1.12

package constbn

import "math/bits"

func mul31Lo(x, y base) base {
	_, lo := bits.Mul32(uint32(x), uint32(y))
	return base(lo) & mask31
}

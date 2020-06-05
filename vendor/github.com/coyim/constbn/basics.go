package constbn

type base uint32

// constant time primitive implementations. the ctl argument has to be base(0) or base(1)

func not(ctl base) base {
	return ctl ^ one
}

// mux returns x if ctl is true, y if it's false
func mux(ctl, x, y base) base {
	return y ^ (-ctl & (x ^ y))
}

func eq(x, y base) base {
	q := x ^ y
	return not((q | -q) >> 31)
}

func neq(x, y base) base {
	q := x ^ y
	return (q | -q) >> 31
}

// gt returns 1 if x > y, 0 otherwise
func gt(x, y base) base {
	z := y - x
	return (z ^ ((x ^ y) & (x ^ z))) >> 31
}

func ge(x, y base) base {
	return not(gt(y, x))
}

func lt(x, y base) base {
	return gt(y, x)
}

// func le(x, y base) base {
// 	return not(gt(x, y))
// }

func ccopy(ctl base, dst, src []base, len base) {
	for i := zero; i < len; i++ {
		x := src[i]
		y := dst[i]
		dst[i] = mux(ctl, x, y)
	}
}

const zero = base(0)
const one = base(1)

func bitLen(x base) base {
	k := neq(x, zero)

	c := gt(x, base(0xFFFF))
	x = mux(c, x>>16, x)
	k += c << 4

	c = gt(x, base(0x00FF))
	x = mux(c, x>>8, x)
	k += c << 3

	c = gt(x, base(0x000F))
	x = mux(c, x>>4, x)
	k += c << 2

	c = gt(x, base(0x0003))
	x = mux(c, x>>2, x)
	k += c << 1

	k += gt(x, base(0x0001))

	return k
}

// func min(x, y base) base {
// 	return mux(gt(x, y), y, x)
// }

// func max(x, y base) base {
// 	return mux(gt(x, y), x, y)
// }

func mul31(x, y base) uint64 {
	return uint64(x) * uint64(y)
}

func zeroes(len base) []base {
	return make([]base, len)
}

func zeroesBytes(len int) []byte {
	return make([]byte, len)
}

/*
 * Zeroize an integer. The announced bit length is set to the provided
 * value, and the corresponding words are set to 0. The ENCODED bit length
 * is expected here.
 */

func zeroize(x []base, bitlen base) {
	x[0] = bitlen
	toZero := (bitlen + 31) >> 5

	copy(x[1:], zeroes(toZero))
}

func zeroizeBytes(x []byte) {
	copy(x, zeroesBytes(len(x)))
}

const mask31 = base(0x7FFFFFFF)

func enc32be(dst []byte, x base) {
	dst[0] = byte(x >> 24)
	dst[1] = byte(x >> 16)
	dst[2] = byte(x >> 8)
	dst[3] = byte(x)
}

// func byteLen(a []base) base {
// 	return baseLen(a) << 2
// }

func baseLen(a []base) base {
	return (a[0] + 31) >> 5
}

// func byteLenWithHeader(a []base) base {
// 	return baseLenWithHeader(a) << 2
// }

func baseLenWithHeader(a []base) base {
	return (a[0] + 63) >> 5
}

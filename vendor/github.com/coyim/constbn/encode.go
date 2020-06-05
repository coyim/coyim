package constbn

/*
 * Encode an integer into its big-endian unsigned representation. The
 * output length in bytes is provided (parameter 'len'); if the length
 * is too short then the integer is appropriately truncated; if it is
 * too long then the extra bytes are set to 0.
 */

func simpleEncode(x []base) []byte {
	result := make([]byte, len(x)*5)
	encode(result, x)
	return result
}

func encode(dst []byte, x []base) {
	xlen := baseLen(x)
	if xlen == 0 {
		zeroizeBytes(dst)
		return
	}
	l := len(dst)
	buf := l
	k := one
	acc := zero
	accLen := uint(0)

	for l != 0 {
		w := zero
		if k <= xlen {
			w = x[k]
		}
		k++
		if accLen == 0 {
			acc = w
			accLen = 31
		} else {
			z := acc | (w << accLen)
			accLen--
			acc = w >> (31 - accLen)
			if l >= 4 {
				buf -= 4
				l -= 4
				enc32be(dst[buf:], z)
			} else {
				switch l {
				case 3:
					dst[buf-3] = byte(z >> 16)
					fallthrough
				case 2:
					dst[buf-2] = byte(z >> 8)
					fallthrough
				case 1:
					dst[buf-1] = byte(z)
				}
				return
			}
		}
	}
}

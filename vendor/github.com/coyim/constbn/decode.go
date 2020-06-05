package constbn

/*
 * Decode an integer from its big-endian unsigned representation. The
 * "true" bit length of the integer is computed and set in the encoded
 * announced bit length (x[0]), but all words of x[] corresponding to
 * the full 'len' bytes of the source are set.
 *
 * CT: value or length of x does not leak.
 */

func simpleDecode(src []byte) []base {
	result := make([]base, (len(src)/2)+2)
	decode(result, src)
	return result
}

func decode(x []base, src []byte) {
	u := len(src)
	v := 1
	acc := uint(0)
	accLen := uint(0)
	for u > 0 {
		u--
		b := src[u]
		acc |= uint(base(b) << accLen)
		accLen += 8
		if accLen >= 31 {
			x[v] = base(acc) & mask31
			v++
			accLen -= 31
			acc = uint(b) >> (8 - accLen)
		}
	}
	if accLen != 0 {
		x[v] = base(acc)
		v++
	}
	x[0] = bitLength(x[1:], v-1)
}

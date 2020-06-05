package constbn

/*
 * Compute the ENCODED actual bit length of an integer. The argument x
 * should point to the first (least significant) value word of the
 * integer. The len 'xlen' contains the number of 32-bit words to
 * access. The upper bit of each value word MUST be 0.
 * Returned value is ((k / 31) << 5) + (k % 31) if the bit length is k.
 *
 * CT: value or length of x does not leak.
 */
func bitLength(x []base, xlen int) base {
	tw := zero
	twk := zero
	for xlen > 0 {
		xlen--
		c := eq(tw, zero)
		w := x[xlen]
		tw = mux(c, w, tw)
		twk = mux(c, base(xlen), twk)
	}
	return (twk << 5) + bitLen(tw)
}

package constbn

/*
 * Convert a modular integer back from Montgomery representation. The
 * integer x[] MUST be lower than m[], but with the same announced bit
 * length. The "m0i" parameter is equal to -(1/m0) mod 2^32, where m0 is
 * the least significant value word of m[] (this works only if m[] is
 * an odd integer).
 */

func fromMonty(x []base, m []base, m0i base) {
	len := baseLen(m)

	for u := zero; u < len; u++ {
		f := mul31Lo(x[1], m0i)
		cc := uint64(0)
		for v := zero; v < len; v++ {
			z := uint64(x[v+1]) + mul31(f, m[v+1]) + cc
			cc = z >> 31
			if v != 0 {
				x[v] = base(z) & mask31
			}
		}
		x[len] = base(cc)
	}

	sub(x, m, not(sub(x, m, zero)))
}

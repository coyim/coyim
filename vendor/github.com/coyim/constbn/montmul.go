package constbn

/*
 * Compute a modular Montgomery multiplication. d[] is filled with the
 * value of x*y/R modulo m[] (where R is the Montgomery factor). The
 * array d[] MUST be distinct from x[], y[] and m[]. x[] and y[] MUST be
 * numerically lower than m[]. x[] and y[] MAY be the same array. The
 * "m0i" parameter is equal to -(1/m0) mod 2^31, where m0 is the least
 * significant value word of m[] (this works only if m[] is an odd
 * integer).
 */

func montmul(d, x, y, m []base, m0i base) {
	len := baseLen(m)
	len4 := len & ^base(3)
	zeroize(d, m[0])
	dh := zero

	for u := zero; u < len; u++ {
		xu := x[u+1]
		f := mul31Lo((d[1] + mul31Lo(x[u+1], y[1])), m0i)
		r := uint64(0)

		v := zero
		for ; v < len4; v += 4 {
			z := uint64(d[v+1]) + mul31(xu, y[v+1]) + mul31(f, m[v+1]) + r
			r = z >> 31
			d[v+0] = base(z) & mask31

			z = uint64(d[v+2]) + mul31(xu, y[v+2]) + mul31(f, m[v+2]) + r
			r = z >> 31
			d[v+1] = base(z) & mask31

			z = uint64(d[v+3]) + mul31(xu, y[v+3]) + mul31(f, m[v+3]) + r
			r = z >> 31
			d[v+2] = base(z) & mask31

			z = uint64(d[v+4]) + mul31(xu, y[v+4]) + mul31(f, m[v+4]) + r
			r = z >> 31
			d[v+3] = base(z) & mask31
		}

		for ; v < len; v++ {
			z := uint64(d[v+1]) + mul31(xu, y[v+1]) + mul31(f, m[v+1]) + r
			r = z >> 31
			d[v] = base(z) & mask31
		}

		dh += base(r)
		d[len] = dh & mask31
		dh >>= 31
	}

	d[0] = m[0]

	sub(d, m, neq(dh, 0)|not(sub(d, m, zero)))
}

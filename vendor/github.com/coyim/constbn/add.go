package constbn

/*
 * Add b[] to a[] and return the carry (0 or 1). If ctl is 0, then a[]
 * is unmodified, but the carry is still computed and returned. The
 * arrays a[] and b[] MUST have the same announced bit length.
 *
 * a[] and b[] MAY be the same array, but partial overlap is not allowed.
 */

func add(a, b []base, ctl base) base {
	cc := zero
	m := baseLenWithHeader(a)
	for u := one; u < m; u++ {
		aw := a[u]
		bw := b[u]
		naw := aw + bw + cc
		cc = naw >> 31
		a[u] = mux(ctl, naw&mask31, aw)
	}
	return cc
}

package constbn

/*
 * Multiply x[] by 2^31 and then add integer z, modulo m[]. This
 * function assumes that x[] and m[] have the same announced bit
 * length, the announced bit length of m[] matches its true
 * bit length.
 *
 * x[] and m[] MUST be distinct arrays. z MUST fit in 31 bits (upper
 * bit set to 0).
 *
 * CT: only the common announced bit length of x and m leaks, not
 * the values of x, z or m.
 */

func muladdSmall(x []base, z base, m []base) {
	mBitlen := m[0]

	if mBitlen == zero {
		return
	}

	if mBitlen <= base(31) {
		hi := x[1] >> 1
		lo := (x[1] << 31) | z
		x[1] = rem(hi, lo, m[1])
		return
	}

	mlen := baseLen(m)
	mblr := mBitlen & 31

	hi := x[mlen]
	var a0, a1, b0 base
	if mblr == zero {
		a0 = x[mlen]
		copy(x[2:], x[1:mlen])
		x[1] = z
		a1 = x[mlen]
		b0 = m[mlen]
	} else {
		a0 = ((x[mlen] << (31 - mblr)) | (x[mlen-1] >> mblr)) & mask31
		copy(x[2:], x[1:mlen])
		x[1] = z
		a1 = ((x[mlen] << (31 - mblr)) | (x[mlen-1] >> mblr)) & mask31
		b0 = ((m[mlen] << (31 - mblr)) | (m[mlen-1] >> mblr)) & mask31
	}
	g := div(a0>>1, a1|(a0<<31), b0)
	q := mux(eq(a0, b0), mask31, mux(eq(g, zero), zero, g-1))

	cc := zero
	tb := one

	for u := one; u <= mlen; u++ {
		mw := m[u]
		zl := mul31(mw, q) + uint64(cc)
		cc = base(zl >> 31)
		zw := base(zl) & mask31
		xw := x[u]
		nxw := xw - zw
		cc += nxw >> 31
		nxw &= mask31
		x[u] = nxw
		tb = mux(eq(nxw, mw), tb, gt(nxw, mw))
	}

	over := gt(cc, hi)
	under := ^over & (tb | lt(cc, hi))
	add(x, m, over)
	sub(x, m, under)
}

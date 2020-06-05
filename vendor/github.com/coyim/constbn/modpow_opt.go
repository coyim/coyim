package constbn

/*
 * Compute a modular exponentiation. x[] MUST be an integer modulo m[]
 * (same announced bit length, lower value). m[] MUST be odd. The
 * exponent is in big-endian unsigned notation, over 'elen' bytes. The
 * "m0i" parameter is equal to -(1/m0) mod 2^31, where m0 is the least
 * significant value word of m[] (this works only if m[] is an odd
 * integer). The tmp[] array is used for temporaries, and has size
 * 'twlen' words; it must be large enough to accommodate at least two
 * temporary values with the same size as m[] (including the leading
 * "bit length" word). If there is room for more temporaries, then this
 * function may use the extra room for window-based optimisation,
 * resulting in faster computations.
 *
 * Returned value is 1 on success, 0 on error. An error is reported if
 * the provided tmp[] array is too short.
 */

func simpleModpowOpt(x []base, e []byte, m []base) []base {
	l := len(x)
	if l < len(m) {
		l = len(m)
	}

	result := make([]base, l)
	copy(result, x)
	result[0] = m[0]
	m0i := ninv(m[1])
	tmp := make([]base, 5*baseLenWithHeader(m))
	_ = modpowOpt(result, e, m, m0i, tmp)
	return result
}

func modpowOpt(x []base, e []byte, m []base, m0i base, tmp []base) base {
	mwlen := baseLenWithHeader(m)
	mlen := mwlen
	mwlen += mwlen & 1

	twlen := base(len(tmp))
	if twlen < (mwlen << 1) {
		return zero
	}

	t1 := tmp
	t2 := tmp[mwlen:]

	winLen := uint(5)
	for ; winLen > 1; winLen-- {
		if ((one<<winLen)+1)*mwlen <= twlen {
			break
		}
	}

	toMonty(x, m)

	if winLen == 1 {
		copy(t2, x[:mlen])
	} else {
		bs := t2[mwlen:]
		copy(bs, x[:mlen])
		for u := base(2); u < one<<winLen; u++ {
			montmul(bs[mwlen:], bs, x, m, m0i)
			bs = bs[mwlen:]
		}
	}

	zeroize(x, m[0])
	x[baseLen(m)] = one
	muladdSmall(x, zero, m)

	acc := base(0)
	accLen := uint(0)
	elen := len(e)
	ep := 0

	for accLen > 0 || elen > 0 {
		k := winLen
		if accLen < winLen {
			if elen > 0 {
				acc = (acc << 8) | base(e[ep])
				ep++
				elen--
				accLen += 8
			} else {
				k = accLen
			}
		}

		bits := (acc >> (accLen - k)) & ((one << k) - one)
		accLen -= k

		for i := uint(0); i < k; i++ {
			montmul(t1, x, x, m, m0i)
			copy(x, t1[:mlen])
		}

		if winLen > 1 {
			zeroize(t2, m[0])
			bs := mwlen
			for u := one; u < (one << k); u++ {
				mask := -eq(u, bits)
				for v := one; v < mwlen; v++ {
					t2[v] |= mask & t2[bs+v]
				}
				bs += mwlen
			}
		}

		montmul(t1, x, t2, m, m0i)
		ccopy(neq(bits, zero), x, t1, mlen)
	}

	fromMonty(x, m, m0i)
	return one
}

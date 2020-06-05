package constbn

func div(hi, lo, d base) base {
	q, _ := divrem(hi, lo, d)
	return q
}

func rem(hi, lo, d base) base {
	_, r := divrem(hi, lo, d)
	return r
}

func divrem(hi, lo, d base) (quo, rem base) {
	q := zero
	ch := eq(hi, d)
	hi = mux(ch, zero, hi)

	for k := base(31); k > zero; k-- {
		j := 32 - k
		w := (hi << j) | (lo >> k)
		ctl := ge(w, d) | (hi >> k)
		hi2 := (w - d) >> j
		lo2 := lo - (d << k)
		hi = mux(ctl, hi2, hi)
		lo = mux(ctl, lo2, lo)
		q |= ctl << k
	}
	cf := ge(lo, d) | hi
	q |= cf
	r := mux(cf, lo-d, lo)

	return q, r
}

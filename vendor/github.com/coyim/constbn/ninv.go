package constbn

/*
 * Compute -(1/x) mod 2^31. If x is even, then this function returns 0.
 */

func ninv(x base) base {
	two := base(2)
	y := two - x
	y *= two - y*x
	y *= two - y*x
	y *= two - y*x
	y *= two - y*x
	return mux(x&one, -y, zero) & mask31
}

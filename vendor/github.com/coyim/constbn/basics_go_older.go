// +build !go1.12

package constbn

func mul31Lo(x, y base) base {
	return base(x*y) & mask31
}

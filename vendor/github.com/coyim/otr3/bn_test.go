package otr3

import (
	"math/big"
	"testing"
)

func Test_mod_returnsTheValueModAnotherValue(t *testing.T) {
	result := mod(big.NewInt(7), big.NewInt(3))
	assertDeepEquals(t, result, big.NewInt(1))
}

func Test_lt_returnsTrueIfTheLeftIsSmallerThanTheRight(t *testing.T) {
	assertDeepEquals(t, lt(big.NewInt(3), big.NewInt(7)), true)
	assertDeepEquals(t, lt(big.NewInt(6), big.NewInt(7)), true)
	assertDeepEquals(t, lt(big.NewInt(7), big.NewInt(7)), false)
	assertDeepEquals(t, lt(big.NewInt(8), big.NewInt(7)), false)
}

func Test_gt_returnsTrueIfTheLeftIsGreaterThanTheRight(t *testing.T) {
	assertDeepEquals(t, gt(big.NewInt(3), big.NewInt(3)), false)
	assertDeepEquals(t, gt(big.NewInt(4), big.NewInt(3)), true)
	assertDeepEquals(t, gt(big.NewInt(7), big.NewInt(3)), true)
}

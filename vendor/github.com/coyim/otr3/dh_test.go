package otr3

import (
	"math/big"
	"testing"
)

func Test_thatIsGroupElementDisallowsThingsLessThanTwo(t *testing.T) {
	assertEquals(t, isGroupElement(new(big.Int).SetInt64(0)), false)
	assertEquals(t, isGroupElement(new(big.Int).SetInt64(1)), false)
	assertEquals(t, isGroupElement(new(big.Int).SetInt64(2)), true)
	assertEquals(t, isGroupElement(new(big.Int).SetInt64(-1)), false)
}

func Test_thatIsGroupElementDisallowsThingsLargerThanTheModuloMinusTwo(t *testing.T) {
	assertEquals(t, isGroupElement(p), false)
	assertEquals(t, isGroupElement(new(big.Int).Add(p, new(big.Int).SetInt64(1))), false)
	assertEquals(t, isGroupElement(new(big.Int).Sub(p, new(big.Int).SetInt64(1))), false)
	assertEquals(t, isGroupElement(new(big.Int).Sub(p, new(big.Int).SetInt64(2))), true)
	assertEquals(t, isGroupElement(new(big.Int).Sub(p, new(big.Int).SetInt64(3))), true)
}

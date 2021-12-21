package gui

import (
	"math"

	. "gopkg.in/check.v1"
)

type ColorsSuite struct{}

var _ = Suite(&ColorsSuite{})

func colorValueCloseEnough(left, right colorValue, tol float64) bool {
	diff := math.Abs(float64(left - right))
	return diff < tol
}

func (s *ColorsSuite) Test_createColorValueFrom_createsAColorValue(c *C) {
	c.Assert(createColorValueFrom(0), Equals, colorValue(0))
	c.Assert(createColorValueFrom(255), Equals, colorValue(1.0))
	c.Assert(colorValueCloseEnough(createColorValueFrom(127), colorValue(0.5), 0.01), Equals, true)
	c.Assert(colorValueCloseEnough(createColorValueFrom(64), colorValue(0.25), 0.01), Equals, true)
	c.Assert(colorValueCloseEnough(createColorValueFrom(192), colorValue(0.75), 0.01), Equals, true)
}

func (s *ColorsSuite) Test_colorValue_toScaledValue_givesGoodResults(c *C) {
	c.Assert(colorValue(0).toScaledValue(), Equals, uint8(0))
	c.Assert(colorValue(1).toScaledValue(), Equals, uint8(255))
	c.Assert(colorValue(0.5).toScaledValue(), Equals, uint8(127))
	c.Assert(colorValue(0.20).toScaledValue(), Equals, uint8(51))
	c.Assert(colorValue(0.63).toScaledValue(), Equals, uint8(160))
}

func (s *ColorsSuite) Test_colorValueFromHex_returnsProperColorValues(c *C) {
	c.Assert(colorValueFromHex("0"), Equals, colorValue(0))
	c.Assert(colorValueCloseEnough(colorValueFromHex("1"), colorValue(0.0039), 0.001), Equals, true)
	c.Assert(colorValueFromHex("ff"), Equals, colorValue(1))
	c.Assert(colorValueFromHex("FF"), Equals, colorValue(1))
	c.Assert(colorValueCloseEnough(colorValueFromHex("a1"), colorValue(0.631), 0.01), Equals, true)
}

func (s *ColorsSuite) Test_colorValueFromHex_returnsZeroOnInabilityToParse(c *C) {
	c.Assert(colorValueFromHex("***"), Equals, colorValue(0))
	c.Assert(colorValueFromHex("1x"), Equals, colorValue(0))
}

func (s *ColorsSuite) Test_rgbFromHex_works(c *C) {
	c.Assert(rgbFromHex("#FF00FF"), DeepEquals, &rgb{colorValue(1), colorValue(0), colorValue(1)})
	c.Assert(rgbFromHex("ff00ff"), DeepEquals, rgbFromHex("#FF00FF"))
	c.Assert(rgbFromHex("#F0F"), DeepEquals, rgbFromHex("#FF00FF"))
	c.Assert(rgbFromHex("000000"), DeepEquals, &rgb{colorValue(0), colorValue(0), colorValue(0)})
	c.Assert(rgbFromHex("blarg"), IsNil)
}

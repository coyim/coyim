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

func (s *ColorsSuite) Test_rgbFrom_works(c *C) {
	c.Assert(rgbFrom(0, 0, 0), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(0),
		blue:  colorValue(0),
	})

	c.Assert(rgbFrom(255, 0, 0), DeepEquals, &rgb{
		red:   colorValue(1),
		green: colorValue(0),
		blue:  colorValue(0),
	})

	c.Assert(rgbFrom(0, 255, 0), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(1),
		blue:  colorValue(0),
	})

	c.Assert(rgbFrom(0, 0, 255), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(0),
		blue:  colorValue(1),
	})

	c.Assert(rgbFrom(255, 0, 255), DeepEquals, &rgb{
		red:   colorValue(1),
		green: colorValue(0),
		blue:  colorValue(1),
	})

	c.Assert(rgbFrom(255, 255, 255), DeepEquals, &rgb{
		red:   colorValue(1),
		green: colorValue(1),
		blue:  colorValue(1),
	})
}

func (s *ColorsSuite) Test_rgbaFrom_works(c *C) {
	c.Assert(rgbaFrom(0, 0, 0, 0), DeepEquals, &rgba{&rgb{
		red:   colorValue(0),
		green: colorValue(0),
		blue:  colorValue(0),
	}, 0.0})

	c.Assert(rgbaFrom(255, 0, 0, 0), DeepEquals, &rgba{&rgb{
		red:   colorValue(1),
		green: colorValue(0),
		blue:  colorValue(0),
	}, 0.0})

	c.Assert(rgbaFrom(0, 255, 0, 0.2), DeepEquals, &rgba{&rgb{
		red:   colorValue(0),
		green: colorValue(1),
		blue:  colorValue(0),
	}, 0.2})

	c.Assert(rgbaFrom(0, 0, 255, 0), DeepEquals, &rgba{&rgb{
		red:   colorValue(0),
		green: colorValue(0),
		blue:  colorValue(1),
	}, 0.0})

	c.Assert(rgbaFrom(255, 0, 255, 1), DeepEquals, &rgba{&rgb{
		red:   colorValue(1),
		green: colorValue(0),
		blue:  colorValue(1),
	}, 1.0})

	c.Assert(rgbaFrom(255, 255, 255, 0.76565665), DeepEquals, &rgba{&rgb{
		red:   colorValue(1),
		green: colorValue(1),
		blue:  colorValue(1),
	}, 0.76565665})
}

func (s *ColorsSuite) Test_rgbFromPercent_works(c *C) {
	c.Assert(rgbFromPercent(0, 0, 0), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(0),
		blue:  colorValue(0),
	})

	c.Assert(rgbFromPercent(1, 0, 0), DeepEquals, &rgb{
		red:   colorValue(1),
		green: colorValue(0),
		blue:  colorValue(0),
	})

	c.Assert(rgbFromPercent(0, 0.7, 0), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(0.7),
		blue:  colorValue(0),
	})

	c.Assert(rgbFromPercent(0, 0, 0.5), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(0),
		blue:  colorValue(0.5),
	})

	c.Assert(rgbFromPercent(0.2, 0, 0.25), DeepEquals, &rgb{
		red:   colorValue(0.2),
		green: colorValue(0),
		blue:  colorValue(0.25),
	})

	c.Assert(rgbFromPercent(1, 0.777, 1), DeepEquals, &rgb{
		red:   colorValue(1),
		green: colorValue(0.777),
		blue:  colorValue(1),
	})
}

type testColorWithGetters struct {
	rr, gg, bb float64
}

func (v *testColorWithGetters) GetRed() float64 {
	return v.rr
}

func (v *testColorWithGetters) GetGreen() float64 {
	return v.gg
}

func (v *testColorWithGetters) GetBlue() float64 {
	return v.bb
}

func (s *ColorsSuite) Test_rgbFromGetters_work(c *C) {
	c.Assert(rgbFromGetters(&testColorWithGetters{0, 0, 0}), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(0),
		blue:  colorValue(0),
	})

	c.Assert(rgbFromGetters(&testColorWithGetters{1, 0, 0}), DeepEquals, &rgb{
		red:   colorValue(1),
		green: colorValue(0),
		blue:  colorValue(0),
	})

	c.Assert(rgbFromGetters(&testColorWithGetters{0, 0.7, 0}), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(0.7),
		blue:  colorValue(0),
	})

	c.Assert(rgbFromGetters(&testColorWithGetters{0, 0, 0.5}), DeepEquals, &rgb{
		red:   colorValue(0),
		green: colorValue(0),
		blue:  colorValue(0.5),
	})

	c.Assert(rgbFromGetters(&testColorWithGetters{0.2, 0, 0.25}), DeepEquals, &rgb{
		red:   colorValue(0.2),
		green: colorValue(0),
		blue:  colorValue(0.25),
	})

	c.Assert(rgbFromGetters(&testColorWithGetters{1, 0.777, 1}), DeepEquals, &rgb{
		red:   colorValue(1),
		green: colorValue(0.777),
		blue:  colorValue(1),
	})
}

func (s *ColorsSuite) Test_rgb_lightness_basedOnOneSingleValue(c *C) {
	col := &rgb{
		red: 0.5,
	}
	c.Assert(col.lightness(), Equals, 0.25)
}

func (s *ColorsSuite) Test_rgb_lightness_basedOnMoreThanOneValue(c *C) {
	col := &rgb{
		red:   0.5,
		green: 0.5,
		blue:  0.5,
	}
	c.Assert(col.lightness(), Equals, 0.5)

	col.green = 0.1
	c.Assert(col.lightness(), Equals, 0.3)
}

func (s *ColorsSuite) Test_rgb_isDark_works(c *C) {
	c.Assert((&rgb{1, 1, 1}).isDark(), Equals, false)
	c.Assert((&rgb{0, 0, 0}).isDark(), Equals, true)
	c.Assert((&rgb{0, 0.5, 0}).isDark(), Equals, true)
	c.Assert((&rgb{0, 0.8, 0.9}).isDark(), Equals, true)
	c.Assert((&rgb{1.0, 0.8, 0.8}).isDark(), Equals, false)
}

func (s *ColorsSuite) Test_rgb_toHex_works(c *C) {
	c.Assert((&rgb{0, 0, 0}).toHex(), Equals, "#000000")
	c.Assert((&rgb{1, 1, 1}).toHex(), Equals, "#ffffff")
	c.Assert((&rgb{1, 0, 0}).toHex(), Equals, "#ff0000")
	c.Assert((&rgb{1, 0, 1}).toHex(), Equals, "#ff00ff")
	c.Assert((&rgb{0, 1, 1}).toHex(), Equals, "#00ffff")
	c.Assert((&rgb{0, 1, 0}).toHex(), Equals, "#00ff00")
	c.Assert((&rgb{0, 0, 1}).toHex(), Equals, "#0000ff")
	c.Assert((&rgb{0.5, 0.25, 0.75}).toHex(), Equals, "#7f3fbf")
}

func (s *ColorsSuite) Test_rgb_toCSS_works(c *C) {
	c.Assert((&rgb{0, 0, 0}).toCSS(), Equals, "rgb(0, 0, 0)")
	c.Assert((&rgb{1, 1, 1}).toCSS(), Equals, "rgb(255, 255, 255)")
	c.Assert((&rgb{1, 0, 0}).toCSS(), Equals, "rgb(255, 0, 0)")
	c.Assert((&rgb{1, 0, 1}).toCSS(), Equals, "rgb(255, 0, 255)")
	c.Assert((&rgb{0, 1, 1}).toCSS(), Equals, "rgb(0, 255, 255)")
	c.Assert((&rgb{0, 1, 0}).toCSS(), Equals, "rgb(0, 255, 0)")
	c.Assert((&rgb{0, 0, 1}).toCSS(), Equals, "rgb(0, 0, 255)")
	c.Assert((&rgb{0.5, 0.25, 0.75}).toCSS(), Equals, "rgb(127, 63, 191)")
}

func (s *ColorsSuite) Test_rgba_toCSS_works(c *C) {
	c.Assert(rgbaFrom(0, 0, 0, 1).toCSS(), Equals, "rgba(0, 0, 0, 1)")
	c.Assert(rgbaFrom(255, 255, 255, 0.1).toCSS(), Equals, "rgba(255, 255, 255, 0.1)")
	c.Assert(rgbaFrom(255, 0, 0, 0.2).toCSS(), Equals, "rgba(255, 0, 0, 0.2)")
	c.Assert(rgbaFrom(255, 0, 255, 0.3).toCSS(), Equals, "rgba(255, 0, 255, 0.3)")
	c.Assert(rgbaFrom(0, 255, 255, 0.5).toCSS(), Equals, "rgba(0, 255, 255, 0.5)")
	c.Assert(rgbaFrom(0, 255, 0, 0.7).toCSS(), Equals, "rgba(0, 255, 0, 0.7)")
	c.Assert(rgbaFrom(0, 0, 255, 0.8).toCSS(), Equals, "rgba(0, 0, 255, 0.8)")
	c.Assert(rgbaFrom(127, 63, 191, 0.95).toCSS(), Equals, "rgba(127, 63, 191, 0.95)")
}

func (s *ColorsSuite) Test_rgba_String_works(c *C) {
	c.Assert(rgbaFrom(0, 0, 0, 1).String(), Equals, "rgba(0, 0, 0, 1)")
}

func (s *ColorsSuite) Test_rgb_String_works(c *C) {
	c.Assert((&rgb{0, 0, 0}).String(), Equals, "rgb(0, 0, 0)")
}

func (s *ColorsSuite) Test_cssColorReferenceFrom_works(c *C) {
	c.Assert(cssColorReferenceFrom("hello world"), DeepEquals, &cssColorReference{"hello world"})
	c.Assert(cssColorReferenceFrom("none"), DeepEquals, &cssColorReference{"none"})
	c.Assert(cssColorReferenceFrom("@theme_bg_color"), DeepEquals, &cssColorReference{"@theme_bg_color"})
}

func (s *ColorsSuite) Test_cssColorReference_toCSS_works(c *C) {
	c.Assert(cssColorReferenceFrom("hello world").toCSS(), Equals, "hello world")
	c.Assert(cssColorReferenceFrom("none").toCSS(), Equals, "none")
	c.Assert(cssColorReferenceFrom("@theme_bg_color").toCSS(), Equals, "@theme_bg_color")
}

func (s *ColorsSuite) Test_cssColorReference_String_works(c *C) {
	c.Assert(cssColorReferenceFrom("hello world").String(), Equals, "hello world")
	c.Assert(cssColorReferenceFrom("none").String(), Equals, "none")
	c.Assert(cssColorReferenceFrom("@theme_bg_color").String(), Equals, "@theme_bg_color")
}

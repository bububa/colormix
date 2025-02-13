package colormix

import (
	"fmt"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
	"gonum.org/v1/gonum/mat"
)

type ColorSpace int

const (
	RGB ColorSpace = iota
	// LAB CIE-L*a*b*: A perceptually uniform color space, i.e. distances are meaningful. L* in [0..1] and a*, b* almost in [-1..1].
	LAB
	// HSL Hue in [0..360], Saturation and Luminance in [0..1]. For legacy reasons; please forget that it exists.
	HSL
	// HSV Hue in [0..360], Saturation and Value in [0..1]. You're better off using HCL
	HSV
	// LUV CIE-L*u*v*: Very similar to CIE-L*a*b*, there is no consensus on which one is "better".
	LUV
	// LCH CIE-L*C*h° (HCL): This is generally the most useful one; CIE-L*a*b* space in polar coordinates, i.e. a better HSV. H° is in [0..360], C* almost in [0..1] and L* as in CIE-L*a*b*.
	LCH
	// HSLuv: The better alternative to HSL, see here and here. Hue in [0..360], Saturation and Luminance in [0..1].
	HSLuv
	// HPLuv: A variant of HSLuv. The color space is smoother, but only pastel colors can be included. Because the valid colors are limited, it's easy to get invalid Saturation values way above 1.0, indicating the color can't be represented in HPLuv because it's not pastel.
	HPLuv
	// HCL
	HCL
)

type ColorMeta struct {
	name            string
	alternativeName string
	serialNo        string
	brandName       string
	formatter       func() string
}

// SetName update the color name
func (c *ColorMeta) SetName(name string) {
	c.name = name
}

// Name returns the name of color
func (c ColorMeta) Name() string {
	return c.name
}

// SetAlternativeName update the color alternative name
func (c *ColorMeta) SetAlternativeName(name string) {
	c.alternativeName = name
}

// AlternativeName returns the alternative name of color
func (c ColorMeta) AlternativeName() string {
	return c.alternativeName
}

// SetBrandName update the color brand name
func (c *ColorMeta) SetBrandName(name string) {
	c.name = name
}

// SetSerialNo update the color serial number
func (c *ColorMeta) SetSerialNo(v string) {
	c.serialNo = v
}

// SerialNo returns the serial number of color
func (c ColorMeta) SerialNo() string {
	return c.serialNo
}

// BrandName returns the brand name of color
func (c ColorMeta) BrandName() string {
	return c.brandName
}

// SetFormatter set String formatter for color meta
func (c *ColorMeta) SetFormatter(fn func() string) {
	c.formatter = fn
}

func (c ColorMeta) defaultFormatter() string {
	name := c.name
	if c.alternativeName != "" {
		name = c.alternativeName
	}
	if c.serialNo != "" {
		return fmt.Sprintf("%s#%s", name, c.serialNo)
	}
	return name
}

func (c ColorMeta) String() string {
	if c.formatter != nil {
		return c.formatter()
	}
	return c.defaultFormatter()
}

// Color named color
type Color struct {
	colorful.Color
	ColorMeta
	ratio float64
}

// MakeColor create a color
func MakeColor(v color.Color) Color {
	ret, _ := colorful.MakeColor(v)
	return Color{
		Color: ret,
	}
}

// Hex create a color
func HexColor(v string) Color {
	ret, _ := colorful.Hex(v)
	return Color{
		Color: ret,
	}
}

// Dense returns vec dense of color
func (c Color) Dense(colorSpace ColorSpace) *mat.VecDense {
	s1, s2, s3 := c.SpaceValues(colorSpace)
	return mat.NewVecDense(3, []float64{s1, s2, s3})
}

// Ratio the color ratio in a palette
func (c Color) Ratio() float64 {
	return c.ratio
}

// SetRatio update the color ratio in a palette
func (c *Color) SetRatio(ratio float64) {
	c.ratio = ratio
}

// SpaceValues returns float64 values for a specific color space of color
func (c Color) SpaceValues(space ColorSpace) (float64, float64, float64) {
	switch space {
	case LAB:
		return c.Lab()
	case LUV:
		return c.Luv()
	case HSV:
		return c.Hsv()
	case HSL:
		return c.Hsl()
	case LCH:
		return c.LuvLCh()
	case HSLuv:
		return c.HSLuv()
	case HPLuv:
		return c.HPLuv()
	case HCL:
		return c.Hcl()
	default:
		return c.LinearRgb()
	}
}

// ColorInSpace returns a Color in specific color space
func ColorInSpace(values []float64, space ColorSpace) color.Color {
	var (
		s1 float64
		s2 float64
		s3 float64
	)
	if l := len(values); l > 0 {
		s1 = values[0]
		if l > 1 {
			s2 = values[1]
			if l > 2 {
				s3 = values[2]
			}
		}
	}
	switch space {
	case LAB:
		return colorful.Lab(s1, s2, s3)
	case LUV:
		return colorful.Luv(s1, s2, s3)
	case HSV:
		return colorful.Hsv(s1, s2, s3)
	case HSL:
		return colorful.Hsl(s1, s2, s3)
	case LCH:
		return colorful.LuvLCh(s1, s2, s3)
	case HSLuv:
		return colorful.HSLuv(s1, s2, s3)
	case HPLuv:
		return colorful.HPLuv(s1, s2, s3)
	case HCL:
		return colorful.Hcl(s1, s2, s3)
	default:
		return colorful.LinearRgb(s1, s2, s3)
	}
}

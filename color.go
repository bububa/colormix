package colormix

import (
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

type IColor interface {
	color.Color
	Hex() string
	Dense(ColorSpace) *mat.VecDense
	Ratio() float64
	SetRatio(float64)
	SpaceValues(space ColorSpace) (float64, float64, float64)
}

// Color named color
type Color struct {
	colorful.Color
	name  string
	ratio float64
}

// MakeColor create a color
func MakeColor(v color.Color) Color {
	ret, _ := colorful.MakeColor(v)
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

// SetName update the color name
func (c *Color) SetName(name string) {
	c.name = name
}

// Name returns the name of color
func (c Color) Name() string {
	return c.name
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

// Palette is named color palette
type Palette struct {
	name   string
	colors []IColor
	cache  map[string]struct{}
}

// NewPalette returns a new palette instance
func NewPalette(colors ...color.Color) *Palette {
	ret := new(Palette)
	l := len(colors)
	ret.colors = make([]IColor, 0, l)
	ret.cache = make(map[string]struct{}, l)
	for _, v := range colors {
		val := MakeColor(v)
		key := val.Hex()
		if _, found := ret.cache[key]; found {
			continue
		}
		ret.cache[key] = struct{}{}
		ret.colors = append(ret.colors, &val)
	}
	return ret
}

// AddColors append colors to the palette
func (p *Palette) AddColors(colors ...IColor) *Palette {
	for _, v := range colors {
		key := v.Hex()
		if _, found := p.cache[key]; found {
			continue
		}
		p.cache[key] = struct{}{}
		p.colors = append(p.colors, v)
	}
	return p
}

// Dense returns a mat dense of palette
func (p Palette) Dense(colorSpace ColorSpace) *mat.Dense {
	cols := 3
	rows := len(p.colors)
	values := make([]float64, 0, rows*cols)
	for _, v := range p.colors {
		s1, s2, s3 := v.SpaceValues(colorSpace)
		values = append(values, s1, s2, s3)
	}
	return mat.NewDense(len(p.colors), cols, values)
}

// SetRatio update the color ratio at index(idx) in the palette
func (p Palette) SetRatio(ratio float64, idx int) {
	if idx < 0 || idx > len(p.colors) {
		return
	}
	p.colors[idx].SetRatio(ratio)
}

// Colors returns a list of colors in palette
func (p Palette) Colors() []IColor {
	return p.colors
}

// SetName set palette name
func (p *Palette) SetName(name string) {
	p.name = name
}

// Name returns the name of palette
func (p Palette) Name() string {
	return p.name
}

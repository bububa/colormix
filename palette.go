package colormix

import (
	"image/color"

	"gonum.org/v1/gonum/mat"
)

// Palette is named color palette
type Palette struct {
	name   string
	colors []Color
	cache  map[string]struct{}
}

// NewPalette returns a new palette instance
func NewPalette(colors ...color.Color) *Palette {
	ret := new(Palette)
	l := len(colors)
	ret.colors = make([]Color, 0, l)
	ret.cache = make(map[string]struct{}, l)
	for _, v := range colors {
		val := MakeColor(v)
		key := val.Hex()
		if _, found := ret.cache[key]; found {
			continue
		}
		ret.cache[key] = struct{}{}
		ret.colors = append(ret.colors, val)
	}
	return ret
}

// AddColors append colors to the palette
func (p *Palette) AddColors(colors ...Color) *Palette {
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
func (p Palette) Colors() []Color {
	return p.colors
}

// Len returns the length of colors in palette
func (p Palette) Len() int {
	return len(p.colors)
}

// SetName set palette name
func (p *Palette) SetName(name string) {
	p.name = name
}

// Name returns the name of palette
func (p Palette) Name() string {
	return p.name
}

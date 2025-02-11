package colormix

import "image/color"

// Mix mix the palette colors with proper ratio to match the target color based on a specific color space
// returns a mixed color or error
func Mix(target color.Color, palette *Palette, space ColorSpace) (color.Color, error) {
	targetDense := MakeColor(target).Dense(space)
	colorsDense := palette.Dense(space)

	// 计算最佳混合比例
	ratios, err := computeColorRatios(colorsDense, targetDense)
	if err != nil {
		return nil, err
	}
	for i := 0; i < ratios.Len(); i++ {
		palette.SetRatio(ratios.AtVec(i), i)
	}
	// 计算混合后的颜色
	mixedColor := mixColors(ratios, colorsDense)
	data := mixedColor.RawVector().Data
	return ColorInSpace(data, space), nil
}

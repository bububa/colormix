package colormix

import (
	"fmt"
	"image/color"
	"testing"
)

func TestMix(t *testing.T) {
	targetColor := color.RGBA{200, 100, 50, 255}
	palette := NewPalette(
		// RED
		color.RGBA{255, 0, 0, 255},
		// GREEN
		color.RGBA{0, 255, 0, 255},
		// BLUE
		color.RGBA{0, 0, 255, 255},
		// WHITE
		color.RGBA{255, 255, 255, 255},
		// BLACK
		color.RGBA{0, 0, 0, 255},
	)
	mixedColor, err := Mix(targetColor, palette, RGB)
	if err != nil {
		t.Fatal(err)
		return
	}
	if MakeColor(mixedColor).Hex() != MakeColor(targetColor).Hex() {
		t.Errorf("wrong expected color: %s, expected %s", MakeColor(mixedColor).Hex(), MakeColor(targetColor).Hex())
	}
	fmt.Printf("Target: %s\n", MakeColor(targetColor).Hex())
	fmt.Println("")
	for _, v := range palette.Colors() {
		fmt.Printf("%s: %.0f%%\n", v.Hex(), v.Ratio()*100)
	}
	fmt.Println("")
	fmt.Printf("Mixed: %s\n", MakeColor(mixedColor).Hex())
	// Output:
	// Target: #c86432
	//
	// #ff0000: 56%
	// #00ff00: 11%
	// #0000ff: 1%
	// #ffffff: 2%
	// #000000: 30%
	//
	// Mixed: #c86432
}

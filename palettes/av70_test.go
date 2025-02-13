package palettes

import (
	"testing"

	"github.com/fatih/color"
)

func TestAV70(t *testing.T) {
	const total int = 205
	if l := AV70.Len(); l != total {
		t.Errorf("wrong total number colors, expect %d, got %d", total, l)
		return
	}
	color.NoColor = false
	for _, v := range AV70.Colors() {
		r, g, b := v.RGB255()
		color.BgRGB(int(r), int(g), int(b)).Println(v)
	}
}

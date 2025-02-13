package palettes

import (
	"encoding/csv"
	"errors"
	"io"
	"slices"

	"github.com/bububa/colormix"
)

var AV70 *colormix.Palette

func init() {
	fn, err := assets.Open("assets/AV70.csv")
	if err != nil {
		panic(err)
	}
	defer fn.Close()
	r := csv.NewReader(fn)
	const brandName = "Acrylicos Vallejo"
	AV70 = colormix.NewPalette()
	for {
		records, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		if len(records) < 4 {
			continue
		}
		slices.Reverse(records)
		color := colormix.HexColor(records[0])
		color.SetName(records[1])
		color.SetAlternativeName(records[2])
		color.SetSerialNo(records[3])
		color.SetBrandName(brandName)
		AV70.AddColors(color)
	}
}

# ColorMix - A Go Library for Color Mixing Optimization

[![Go Reference](https://pkg.go.dev/badge/github.com/bububa/colormix.svg)](https://pkg.go.dev/github.com/bububa/colormix)
[![Go](https://github.com/bububa/colormix/actions/workflows/go.yml/badge.svg)](https://github.com/bububa/colormix/actions/workflows/go.yml)
[![goreleaser](https://github.com/bububa/colormix/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/bububa/colormix/actions/workflows/goreleaser.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/bububa/colormix.svg)](https://github.com/bububa/colormix)
[![GoReportCard](https://goreportcard.com/badge/github.com/bububa/colormix)](https://goreportcard.com/report/github.com/bububa/colormix)
[![GitHub license](https://img.shields.io/github/license/bububa/colormix.svg)](https://github.com/bububa/colormix/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/bububa/colormix.svg)](https://GitHub.com/bububa/colormix/releases/)

_ColorMix_ is a Go library that optimizes color mixing by adjusting the proportions of a given set of colors to match a target color. The library uses optimization algorithms to find the best weights (proportions) for a set of colors in order to minimize the difference between the mixed color and the target color. The optimization is based on color differences, and it supports advanced optimization techniques.

## Features

- Mix colors to match a target color.
- Uses CIEDE2000 color difference metric for better perceptual accuracy.
- Supports different optimization algorithms like Nelder-Mead for general use cases.
- Prevents zero color contributions through penalty functions.

## Installation

To install the ColorMix library, use the go get command:

```bash
go get -u github.com/bububa/colormix
```

## Example Usage

### Import the package

```golong
import "github.com/yourusername/colormix"
```

### Basic Example: Mix Colors to Match a Target

Here is a simple example of how to use the library to mix a set of colors and match a target color.

```golang
package main

import (
	"fmt"
	"image/color"

	"github.com/bububa/colormix"
)

func main() {
	targetColor := color.RGBA{200, 100, 50, 255}
	palette := colormix.NewPalette(
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
	mixedColor, err := colormix.Mix(targetColor, palette, RGB)
	if err != nil {
		return
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
```

### Explanation

1. _palette_: This is a matrix containing the available colors in RGB format.
2. _targetColor_: The target color that we want to match through mixing.
3. _Mix_: This function finds the optimal proportions (weights) of the available colors that best approximate the target color.
4. _mixedColor_: Once the optimal weights are calculated, this function computes the resulting mixed color based on the color table and the calculated weights.

## Color Metrics

The library uses the CIEDE2000 color difference metric to calculate the color differences, which provides a more perceptually accurate measure of color difference than Euclidean distance in RGB space.

## Available Functions

MixColors(targetColor color.Color, palette \*colormix.Palette, colorspace colormix.ColorMix) (mixedColor color.Color, err error): Computes the optimal weights for the given color table to match the target color and Calculates the final mixed color based on the optimal weights and the color table.

## Optimizers Supported

- _Nelder-Mead_: A simplex-based optimization method that doesn't require derivatives. Suitable for non-differentiable functions.
- _LBFGS (Limited-memory Broyden–Fletcher–Goldfarb–Shanno)_: A gradient-based method that is more efficient for smooth, differentiable functions.

## Contributing

We welcome contributions! If you find a bug or want to suggest a feature, feel free to open an issue or submit a pull request.

### Steps to Contribute:

1. Fork the repository.
2. Create a new branch (git checkout -b feature-name).
3. Make your changes and commit them (git commit -am 'Add new feature').
4. Push to the branch (git push origin feature-name).
5. Open a pull request.

## License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.

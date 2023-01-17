package processing

import (
	"github.com/ProninIgorr/fingerprint/internal/imgs/matrix"
	"image"
	"math"
)

func doNormalize(in, out *matrix.M, bounds image.Rectangle, min, max float64) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := in.At(x, y)
			normalizedPixel := math.MaxUint8 * (pixel - min) / (max - min)
			out.Set(x, y, normalizedPixel)
		}
	}
}

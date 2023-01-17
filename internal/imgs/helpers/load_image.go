package helpers

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/ProninIgorr/fingerprint/internal/imgs/matrix"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
)

func maybeResizeImage(img image.Image, maxDimension int) image.Image {
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	if dx <= maxDimension && dy <= maxDimension {
		return img
	}

	xp, yp := 0, 0
	if dx > dy {
		xp = maxDimension
		yp = int(float64(dy) / (float64(dx) / float64(maxDimension)))
	} else if dy > dx {
		yp = maxDimension
		xp = int(float64(dx) / (float64(dy) / float64(maxDimension)))
	} else {
		xp, yp = maxDimension, maxDimension
	}
	return resize.Resize(uint(xp), uint(yp), img, resize.Bilinear)
}

// LoadImage opens a file and attempts to decode the image
// If the dimensions of the image are bigger than expected, then
// the image is resized to fit the expected resolution.
func LoadImage(fname string, maxDimension int) (*matrix.M, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var img image.Image

	ext := strings.ToLower(path.Ext(fname))
	if ext == ".jpg" {
		img, err = jpeg.Decode(f)
	} else if ext == ".png" {
		img, err = png.Decode(f)
	} else if ext == ".bmp" {
		img, err = bmp.Decode(f)
	} else {
		return nil, fmt.Errorf("unsupported image extention: %s", ext)
	}
	if err != nil {
		return nil, err
	}
	resizedImg := maybeResizeImage(img, maxDimension)

	bounds := resizedImg.Bounds()
	gray := image.NewGray(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			c := resizedImg.At(x, y)
			gray.Set(x, y, color.GrayModel.Convert(c))
		}
	}

	return matrix.NewFromGray(gray), nil
}

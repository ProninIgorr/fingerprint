package helpers

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"

	"github.com/ProninIgorr/fingerprint/internal/matrix"
	"github.com/nfnt/resize"
	"github.com/sergeymakinen/go-bmp"
)

const maxDimension = 300

func maybeResizeImage(img image.Image) image.Image {
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
	log.Printf("resizing image from (%d,%d) to (%d,%d)", dx, dy, xp, yp)
	return resize.Resize(uint(xp), uint(yp), img, resize.Bilinear)
}

// LoadImage открывает файл и пытается декодировать изображение
// Если размер изображения больше чем ожидалось, то
// размер изображения изменяется в соответсвии с ожидаемым решением.
func LoadImage(fname string) *matrix.M {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("cannot open image %s, err: %s", fname, err)
	}
	defer f.Close()

	var img image.Image

	ext := path.Ext(fname)
	if ext == ".jpg" {
		img, err = jpeg.Decode(f)
	} else if ext == ".png" {
		img, err = png.Decode(f)
	} else if ext == ".bmp" {
		img, err = bmp.Decode(f)

	} else {
		log.Fatalf("%q extension not supported", ext)
	}
	if err != nil {
		log.Fatalf("cannot decode image %s, err: %s", fname, err)
	}

	resizedImg := maybeResizeImage(img)
	bounds := resizedImg.Bounds()
	gray := image.NewGray(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			c := resizedImg.At(x, y)
			gray.Set(x, y, color.GrayModel.Convert(c))
		}
	}

	return matrix.NewFromGray(gray)
}

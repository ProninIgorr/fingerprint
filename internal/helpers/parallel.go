package helpers

import (
	"image"
	"sync"

	"github.com/ProninIgorr/fingerprint/internal/matrix"
)

// GenerateSubBounds на основе входной матрицы.Полезно для запуска алгоритмов
// паралелльно, каждый состредотачиваясь на подгранице
func GenerateSubBounds(in *matrix.M, offset int) []image.Rectangle {
	subBounds := []image.Rectangle{}
	bounds := in.Bounds()
	blockSize := bounds.Max.X / 2
	for x := bounds.Min.X; x < bounds.Max.X; x += blockSize {
		xi := x - offset
		if xi < bounds.Min.X {
			xi = bounds.Min.X
		}
		xp := x + blockSize + offset
		if xp > bounds.Max.X {
			xp = bounds.Max.X
		}
		for y := bounds.Min.Y; y < bounds.Max.Y; y += blockSize {
			yi := y - offset
			if yi < bounds.Min.Y {
				yi = bounds.Min.Y
			}
			yp := y + blockSize + offset
			if yp > bounds.Max.Y {
				yp = bounds.Max.Y
			}
			subBound := image.Rect(xi, yi, xp, yp)
			subBounds = append(subBounds, subBound)
		}
	}
	return subBounds
}

func RunInParallel(in *matrix.M, offset int, fn func(wg *sync.WaitGroup, bounds image.Rectangle)) {
	subBounds := GenerateSubBounds(in, offset)
	wg := new(sync.WaitGroup)
	wg.Add(len(subBounds))
	for _, bounds := range subBounds {
		go fn(wg, bounds)
	}
	wg.Wait()
}

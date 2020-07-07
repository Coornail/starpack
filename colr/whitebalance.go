package colr

import (
	"image"
	"image/color"

	colorful "github.com/lucasb-eyer/go-colorful"
)

// ModifiedGrayWorld algorithm for white balance.
// Based on https://ieeexplore.ieee.org/document/6269338/
func ModifiedGrayWorld(img image.Image) *image.NRGBA64 {
	var rAvg, gAvg, bAvg float64
	var aAvg float64

	var c color.Color
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c = img.At(x, y)

			r, g, b, _ := c.RGBA()

			rAvg += float64(r)
			gAvg += float64(g)
			bAvg += float64(b)
		}
	}

	n := float64(bounds.Max.X * bounds.Max.Y)
	rAvg = rAvg / n
	gAvg = gAvg / n
	bAvg = bAvg / n

	aAvg = (rAvg + gAvg + bAvg) / 3.0

	rScale := aAvg - rAvg
	gScale := aAvg - gAvg
	bScale := aAvg - bAvg

	res := image.NewNRGBA64(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c = img.At(x, y)

			r, g, b, _ := c.RGBA()

			newColor := colorful.Color{
				R: cap((float64(r)+rScale)/65535.0, 0.0, 1.0),
				G: cap((float64(g)+gScale)/65535.0, 0.0, 1.0),
				B: cap((float64(b)+bScale)/65535.0, 0.0, 1.0),
			}

			res.Set(x, y, newColor)
		}
	}

	return res
}

func cap(v, min, max float64) float64 {
	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}

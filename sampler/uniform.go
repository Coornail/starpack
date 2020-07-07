package sampler

import (
	"image"
	"math"
)

type UniformSampler struct {
	Reference image.Image
	XSamples  int
	YSamples  int

	bounds image.Rectangle

	x int
	y int

	startX int
	startY int

	stepX int
	stepY int
}

func (us UniformSampler) HasMore() bool {
	return us.y < us.bounds.Max.Y
}

func (us *UniformSampler) Next() (x, y int) {
	x, y = us.x, us.y

	us.x += us.stepX

	if us.x > us.bounds.Max.X {
		us.y += us.stepY
		us.x = us.startX
	}

	return
}

func (us *UniformSampler) Reset() {
	us = NewUniformSampler(us.Reference, us.XSamples*us.XSamples)
}

func NewUniformSampler(img image.Image, samples int) *UniformSampler {
	bounds := img.Bounds()

	// @todo use ratio instead.
	xSamples := int(math.Sqrt(float64(samples)))
	ySamples := int(math.Sqrt(float64(samples)))

	stepX := bounds.Max.X / xSamples
	stepY := bounds.Max.Y / ySamples

	startX := (bounds.Max.X / xSamples) / 2
	startY := (bounds.Max.Y / ySamples) / 2

	us := UniformSampler{
		Reference: img,
		bounds:    bounds,
		x:         startX,
		y:         startY,

		XSamples: xSamples,
		YSamples: ySamples,

		stepX: stepX,
		stepY: stepY,

		startX: startX,
		startY: startY,
	}

	return &us
}

package sampler

import (
	"image"
	"math/rand"
	"time"
)

type GaussSampler struct {
	Reference        image.Image
	MaxSamples       int
	RemainingSamples int

	bounds image.Rectangle
}

func (rs GaussSampler) HasMore() bool {
	return rs.RemainingSamples > 0
}

func (rs *GaussSampler) Next() (x, y int) {
	rs.RemainingSamples--

	xMax := float64(rs.bounds.Max.X) / 2.0
	yMax := float64(rs.bounds.Max.Y) / 2.0

	return int(rand.NormFloat64()*(xMax/5) + xMax), int(rand.NormFloat64()*(yMax/5) + yMax)
}

func (rs *GaussSampler) Reset() {
	rs.RemainingSamples = rs.MaxSamples
}

func NewGaussSampler(img image.Image, samples int) *GaussSampler {
	rand.Seed(time.Now().UTC().UnixNano())
	return &GaussSampler{
		Reference:        img,
		RemainingSamples: samples,
		MaxSamples:       samples,
		bounds:           img.Bounds(),
	}
}

//@TODO combine samplers using a new general type

//func combineSamplers(...Sampler) Sampler

//@TODO seed properly
//@TODO Check different image sizes

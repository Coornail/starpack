package sampler

import (
	"image"
)

type CombinedSampler struct {
	Samplers []ImageSampler

	MaxSamples       int
	RemainingSamples int
	i                int
}

func (cs CombinedSampler) HasMore() bool {
	return !(cs.i == len(cs.Samplers)-1 && cs.RemainingSamples == 0)
}

func (cs *CombinedSampler) Next() (x, y int) {
	if cs.RemainingSamples == 0 {
		cs.i++
		cs.RemainingSamples = cs.MaxSamples
	}
	cs.RemainingSamples--

	return cs.Samplers[cs.i].Next()
}

func (cs *CombinedSampler) Reset() {
	cs.RemainingSamples = cs.MaxSamples
	for i := range cs.Samplers {
		cs.Samplers[i].Reset()
	}
	cs.i = 0
}

func NewCombinedSampler(img image.Image, samples int, samplers ...ImageSampler) *CombinedSampler {
	return &CombinedSampler{
		Samplers:         samplers,
		MaxSamples:       samples,
		RemainingSamples: samples,
	}
}

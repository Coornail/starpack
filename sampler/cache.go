package sampler

import (
	"image"
	"sort"
)

type ImageSamplerCache struct {
	Sampler *ImageSampler // @todo we will not need this later

	serveCache bool
	cache      []image.Point
	i          int
}

func (isc *ImageSamplerCache) HasMore() bool {
	var hasMore bool
	if !isc.serveCache {
		hasMore = (*isc.Sampler).HasMore()
		if !hasMore {
			isc.serveCache = true
			isc.Compress()
		}
	} else {
		hasMore = isc.i < len(isc.cache)
	}

	return hasMore
}

func (isc *ImageSamplerCache) Next() (x, y int) {
	if isc.serveCache {
		x, y := isc.cache[isc.i].X, isc.cache[isc.i].Y
		isc.i++

		return x, y
	} else {
		x, y := (*isc.Sampler).Next()
		isc.cache = append(isc.cache, image.Point{X: x, Y: y})
		return x, y
	}
}

func (isc *ImageSamplerCache) Reset() {
	isc.i = 0
}

func (isc *ImageSamplerCache) Compress() {
	// Sort points.
	// Hopefully it will provide better cache-locality.
	sort.Slice(isc.cache, func(i, j int) bool {
		return isc.cache[i].Y < isc.cache[j].Y &&
			isc.cache[i].Y == isc.cache[j].Y && isc.cache[i].X < isc.cache[j].X
	})

	// Remove duplicate points.
	length := len(isc.cache)
	for i := 0; i < length-1; i++ {
		if isc.cache[i] == isc.cache[i+1] {
			length--
			isc.cache = isc.cache[:i+copy(isc.cache[i:], isc.cache[i+1:])]
		}
	}

	// Free up the sampler.
	// Generally samplers hold the image in ram, we can save a few gigabytes here.
	isc.Sampler = nil
}

func NewSamplerCache(sampler ImageSampler) *ImageSamplerCache {
	cache := make([]image.Point, 0)
	return &ImageSamplerCache{
		Sampler: &sampler,

		cache: cache,
	}
}

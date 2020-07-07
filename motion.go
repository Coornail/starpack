package main

import (
	"bytes"
	"encoding/json"
	"image"
	"io/ioutil"
	"math"
	"os"

	"github.com/Coornail/superres/sampler"
	colorful "github.com/lucasb-eyer/go-colorful"
)

type Motion struct {
	X    int
	Y    int
	Diff float64
}

const (
	// We make  (imageX*maxMotion)*(imageY*maxMotion)*(imageX*ImageSampleX)*(imageY*ImageSampleY)
	// comparison on the potentially supersampled image.
	maxMotion = 0.05

	ImageSamples = 1024 * 8

	// Direction change in the spiral since we saw improvement in the picture differences.
	// We change direction 4 times to go "full circle".
	MaxDirectionChangeSinceImprovement = 32
)

// getOutliers returns the indexes for every motion that is over one standard deviation from the mean.
func getOutliers(motions []Motion) []int {
	outliers := make([]int, 0)

	var mean float64
	var sum float64
	for _, m := range motions {
		sum += m.Diff
	}

	N := float64(len(motions))

	mean = sum / N

	sum = 0
	for _, m := range motions {
		d := mean - m.Diff
		sum += d * d
	}

	deviation := math.Sqrt(sum / N)

	for i := range motions {
		d := motions[i].Diff - mean
		if d > deviation {
			outliers = append(outliers, i)
		}
	}

	return outliers
}

// estimateMotion tries to move the candidate image to best match the reference image.
// Comparing the reference image works by taking a sample (@see GetSampler) from both images and calculate the sum of square color differences.
func estimateMotion(reference, candidate image.Image) Motion {
	bounds := reference.Bounds()
	ref := NewImageCache(reference)

	var bestXMotion, bestYMotion int
	var bestDist = math.MaxFloat64

	maxXMotion := math.Round(float64(bounds.Max.X) * maxMotion)
	maxYMotion := math.Round(float64(bounds.Max.Y) * maxMotion)

	max := int(math.Max(maxXMotion, maxYMotion))
	m2 := max * max

	var xMotion, yMotion int
	var dx, dy = 0, -1 // Direction.

	directionChangeSinceImprovement := 0

	smp := GetSampler(reference, ImageSamples)

	var currentDist float64

	// Based on: https://stackoverflow.com/questions/398299/looping-in-a-spiral
	for i := 0; i < m2; i++ {
		if (-max/2 < xMotion && xMotion <= max/2) && (-max/2 < yMotion && yMotion <= max/2) {
			currentDist = 0
			numberOfPixelsCompared := 0

			smp.Reset()
			for smp.HasMore() {
				x, y := smp.Next()
				if x+xMotion < bounds.Min.X || x+xMotion > bounds.Max.X ||
					y+yMotion < bounds.Min.Y || y+yMotion > bounds.Max.Y {
					//fmt.Printf("Out of bounds: %d %d\n", x, y)
					// @todo why does it go out of bounds?
					continue
				}

				referencePoint := ref.At(x, y)
				candidatePoint := candidate.At(x+xMotion, y+yMotion)

				d := distance(referencePoint, rgbaToColorful(candidatePoint))
				currentDist += d * d
				numberOfPixelsCompared++
			}

			currentDist = currentDist / float64(numberOfPixelsCompared)

			if numberOfPixelsCompared > 0 && currentDist < bestDist {
				bestXMotion = xMotion
				bestYMotion = yMotion
				bestDist = currentDist
				directionChangeSinceImprovement = 0
			}
		}

		// Change direction.
		if xMotion == yMotion || (xMotion < 0 && xMotion == -yMotion) || (xMotion > 0 && xMotion == 1-yMotion) {
			dx, dy = -dy, dx
			directionChangeSinceImprovement++
		}
		xMotion, yMotion = xMotion+dx, yMotion+dy

		// If we haven't found an improvement for a long time, we give up.
		if directionChangeSinceImprovement > MaxDirectionChangeSinceImprovement {
			return Motion{X: bestXMotion, Y: bestYMotion, Diff: bestDist}
		}
	}

	return Motion{X: bestXMotion, Y: bestYMotion, Diff: bestDist}
}

// GetSampler returns a sampling implementation for the image.
// Comparing the whole picture would be too computational intensive, so we are forced to choose a subset of pixels to compare.
func GetSampler(img image.Image, samples int) sampler.ImageSampler {
	switch samplerName {
	case "uniform":
		return sampler.NewUniformSampler(img, samples)
	case "edge":
		// Although it's slow to calculate the edges, it gives us the best indication when the intensity will change, hence delivering the best result.
		// Unlike the others it is determinstic.
		return sampler.NewSamplerCache(sampler.NewEdgeDetector(img, samples))
	case "gauss":
		return sampler.NewSamplerCache(sampler.NewGaussSampler(img, samples))
	default:
		s1 := sampler.NewGaussSampler(img, samples/2)
		s2 := sampler.NewEdgeDetector(img, samples/2)
		s := sampler.NewCombinedSampler(img, samples, s1, s2)
		return sampler.NewSamplerCache(s)
	}
}

type ImageCache struct {
	Img       image.Image
	cache     map[int]map[int]colorful.Color
	CacheHit  int
	CacheMiss int
}

func (ic *ImageCache) At(x, y int) (res colorful.Color) {
	if cacheX, foundX := ic.cache[x]; foundX {
		if res, foundY := cacheX[y]; foundY {
			ic.CacheHit++
			return res
		}
	} else {
		ic.cache[x] = make(map[int]colorful.Color)
	}

	ic.CacheMiss++
	color := rgbaToColorful(ic.Img.At(x, y))
	ic.cache[x][y] = color

	return color
}

func NewImageCache(img image.Image) *ImageCache {
	var ic ImageCache

	ic.Img = img
	ic.cache = make(map[int]map[int]colorful.Color, 0)

	return &ic
}

type MotionCache map[string]Motion

func (ms MotionCache) WriteToFile(filename string) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(ms); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf.Bytes(), os.FileMode(0666|os.O_CREATE|os.O_TRUNC))
}

func (ms *MotionCache) ReadFromFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, &ms)
}

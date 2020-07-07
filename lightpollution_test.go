package main

import (
	"os"
	"testing"

	"golang.org/x/image/tiff"
)

const maskOutputFile = "./test_images/light_pollution_out.tif"

// Tests light pollution mask.
func TestLightPollution(*testing.T) {
	loadedImages := loadImages([]string{"./test_images/light_pollution_in.jpg"})

	mask := EstimateLightPollutionMask(loadedImages[0])
	// EstimateLightPollutionMask(loadedImages[0])

	f, _ := os.Create(maskOutputFile)
	defer f.Close()
	tiff.Encode(f, mask, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
}

package main

import (
	"image"

	"github.com/disintegration/imaging"
)

const downSamplePoints = 16

// EstimateLightPollutionMask generates a mask to remove it from the image.
// Based on the idea from https://benedikt-bitterli.me/astro/ .
func EstimateLightPollutionMask(img image.Image) image.Image {
	downsampled := imaging.Resize(img, downSamplePoints, downSamplePoints, imaging.Gaussian)
	// @todo improve on upscaling.
	upsampled := imaging.Resize(downsampled, img.Bounds().Max.X, img.Bounds().Max.Y, imaging.Lanczos)
	upsampled = imaging.Blur(upsampled, 1.5)

	return upsampled
}

package main

import (
	"os"

	starpack "github.com/coornail/starpack/lib"
	"golang.org/x/image/tiff"
)

func main() {
	mask := starpack.EstimateLightPollutionMask(starpack.LoadImage(os.Args[1]))

	f, _ := os.Create(os.Args[2])
	defer f.Close()
	tiff.Encode(f, mask, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
}

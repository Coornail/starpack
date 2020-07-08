package main

import (
	"os"

	starpack "github.com/coornail/starpack/lib"
)

func main() {
	mask := starpack.EstimateLightPollutionMask(starpack.LoadImage(os.Args[1]))
	starpack.SaveImage(os.Args[2], mask)
}

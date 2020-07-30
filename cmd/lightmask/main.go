package main

import (
	"os"

	starpack "github.com/Coornail/starpack/lib"
)

func main() {
	mask := starpack.EstimateLightPollutionMask(starpack.LoadImage(os.Args[1]))
	starpack.SaveImage(os.Args[2], mask)
}

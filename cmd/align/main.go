package main

import (
	"fmt"
	"image/png"
	"os"

	starpack "github.com/Coornail/starpack/lib"
	"github.com/Coornail/starpack/starmap"
)

func main() {
	ref := starpack.LoadImage(os.Args[1])
	target := starpack.LoadImage(os.Args[2])

	sm1 := starpack.GetStarmap(ref)
	sm2 := starpack.GetStarmap(target)

	// sm2.Offset(1, 40)

	diff := starmap.Starmaps{sm1, sm2}.VisualizeDifference()
	diffOut := starmap.Starmaps{sm1, sm2}.CorrectPixels()
	fmt.Printf("%d\n", diffOut)

	f, _ := os.Create("./difference.png")
	defer f.Close()

	png.Encode(f, diff)
}

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	starpack "github.com/Coornail/starpack/lib"
	"github.com/Coornail/starpack/starmap"
)

func main() {
	ref := starpack.LoadImage(os.Args[1])
	target := starpack.LoadImage(os.Args[2])

	sm1, _ := starpack.GetStarmap(ref, 0)
	sm2, _ := starpack.GetStarmap(target, 0)

	m1 := starmap.Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 20, Y: 20}},
		Stars:  sm1.Stars,
	}

	m2 := starmap.Starmap{
		Bounds: image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 20, Y: 20}},
		Stars:  sm2.Stars,
	}

	diff := starmap.Starmaps{sm1, sm2}.VisualizeDifference()
	f, _ := os.Create("./difference.png")
	defer f.Close()
	png.Encode(f, diff)

	offset, _ := m1.FindOffset(m2)
	fmt.Printf("%#v\n", offset)
	m2 = m2.Offset(float64(offset.X), float64(offset.Y))
	m2 = m2.Rotate(offset.Rotation)

	sm2.Stars = m2.Stars

	diff = starmap.Starmaps{sm1, sm2}.VisualizeDifference()
	f2, _ := os.Create("./difference_after.png")
	defer f2.Close()
	png.Encode(f2, diff)
}

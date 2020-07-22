package main

import (
	"os"

	starpack "github.com/Coornail/starpack/lib"
)

func main() {
	img := starpack.LoadImage(os.Args[1])
	sm := starpack.GetStarmap(img)

	sm.WriteFile("./out.png")
}

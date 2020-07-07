package main

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	colorful "github.com/lucasb-eyer/go-colorful"
)

type ColorMerge func([]colorful.Color) colorful.Color

func averageColor(colors []colorful.Color) colorful.Color {
	var h, c, l float64
	for i := range colors {
		currH, currC, currL := colors[i].Hcl()
		h += currH
		c += currC
		l += currL
	}

	count := float64(len(colors))

	return colorful.Hcl(h/count, c/count, l/count).Clamped()
}

func brightestColor(colors []colorful.Color) colorful.Color {
	brightestValue := -1.0
	brightestColor := 0

	for i := range colors {
		_, _, l := colors[i].Hcl()
		if brightestValue < l {
			brightestColor = i
		}
	}

	return colors[brightestColor]
}

func medianColor(colors []colorful.Color) colorful.Color {
	c := len(colors)
	if c == 1 {
		return colors[0]
	}

	l := make([]float64, c)
	a := make([]float64, c)
	b := make([]float64, c)

	for i := range colors {
		l[i], a[i], b[i] = colors[i].Lab()
	}

	// How do I order colors? By luminoscence?
	sort.Slice(colors, func(i, j int) bool {
		return l[i] < l[j]
	})

	if c%2 == 1 {
		return colorful.Lab(l[c/2], a[c/2], b[c/2])
	}

	i := int(math.Floor(float64(c) / 2.0))
	return colorful.Lab((l[i]+l[i-1])/2.0, (a[i]+a[i-1])/2.0, (b[i]+b[i-1])/2.0).Clamped()
}

func rgbaToColorful(c color.Color) colorful.Color {
	res, _ := colorful.MakeColor(c)
	return res.Clamped()
}

func distance(c1, c2 colorful.Color) float64 {
	d := c1.DistanceCIEDE2000(c2)
	if math.IsNaN(d) {
		fmt.Printf("%s + %s = %#v\n", c1.Hex(), c2.Hex(), d)
		panic("Color distance is NaN")
	}

	if d < -1.0 {
		return -1.0
	}

	if d > 1.0 {
		return 1.0
	}

	return d
}

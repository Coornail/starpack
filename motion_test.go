package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"testing"
)

func TestMotionSameImage(t *testing.T) {
	file1, _ := os.Open("motion_1.jpg")
	defer file1.Close()
	img1, _, err := image.Decode(file1)
	if err != nil {
		panic(err)
	}

	m := estimateMotion(img1, img1)
	if m.X != 0 || m.Y != 0 {
		fmt.Printf("%#v\n", m)
		t.Errorf("Same image should not detect motion")
	}
}

func TestMotionSameImageUpscaled(t *testing.T) {
	file1, _ := os.Open("motion_1.jpg")
	defer file1.Close()
	img1, _, err := image.Decode(file1)
	if err != nil {
		panic(err)
	}

	upscaled := upscale([]image.Image{img1})[0]

	m := estimateMotion(upscaled, upscaled)
	if m.X != 0 || m.Y != 0 {
		fmt.Printf("%#v\n", m)
		t.Errorf("Same image should not detect motion")
	}
}

func TestMotion(t *testing.T) {
	file1, _ := os.Open("motion_1.jpg")
	defer file1.Close()
	img1, _, err := image.Decode(file1)
	if err != nil {
		panic(err)
	}

	file2, _ := os.Open("motion_2.jpg")
	defer file1.Close()
	img2, _, err := image.Decode(file2)
	if err != nil {
		panic(err)
	}

	m := estimateMotion(img1, img2)
	if m.X != 16 || m.Y != 22 {
		t.Errorf("Did not find correct motion for the example images")
	}
}

func TestMotionUpscaled(t *testing.T) {
	file1, _ := os.Open("motion_1.jpg")
	defer file1.Close()
	img1, _, err := image.Decode(file1)
	if err != nil {
		panic(err)
	}

	file2, _ := os.Open("motion_2.jpg")
	defer file1.Close()
	img2, _, err := image.Decode(file2)
	if err != nil {
		panic(err)
	}

	m := estimateMotion(upscale([]image.Image{img1})[0], upscale([]image.Image{img2})[0])
	if m.X != 32 || m.Y != 44 {
		t.Errorf("Did not find correct motion for the example images")
	}
}

func TestOutliers(t *testing.T) {
	motions := []Motion{
		{X: 0, Y: 0, Diff: 0.017066}, // <- Outlier
		{X: -5, Y: -1, Diff: 0.001339},
		{X: 6, Y: 1, Diff: 0.001792},
		{X: -10, Y: 13, Diff: 0.002762},
		{X: -13, Y: 1, Diff: 0.002262},
		{X: -32, Y: 22, Diff: 0.001811},
		{X: -47, Y: 46, Diff: 0.002215},
		{X: 0, Y: -8, Diff: 0.002053},
	}

	res := getOutliers(motions)
	if len(res) != 1 {
		t.Errorf("Could not find outlier")
	}

	if res[0] != 0 {
		t.Errorf("Could not mark first item as outlier")
	}
}

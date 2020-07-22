package starpack

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync"

	"github.com/Coornail/starpack/starmap"
	"github.com/disintegration/imaging"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/errors"
	"golang.org/x/image/tiff"
)

const (
	delta = 0.1
)

func Starpack(images []image.Image, colorMergeMethod ColorMerge) *image.NRGBA64 {
	bounds := images[0].Bounds()
	output := image.NewNRGBA64(bounds)

	var currentColor []colorful.Color
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			currentColor = make([]colorful.Color, len(images))

			wg.Add(1)
			go func(x, y int, currentColor []colorful.Color) {
				defer wg.Done()
				for i := range images {
					currX := x
					currY := y
					if currX < bounds.Min.X || currX >= bounds.Max.X ||
						currY < bounds.Min.Y || currY >= bounds.Max.Y {
						continue
					}

					currentColor[i] = rgbaToColorful(images[i].At(currX, currY))
				}
				output.Set(x, y, colorMergeMethod(currentColor))
			}(x, y, currentColor)
		}
		wg.Wait()
		fmt.Printf("Merging: %.2f%%\r", float64(y)/float64(bounds.Max.Y)*100.0)
	}
	fmt.Printf("\n")

	return output
}

func RemoveLightPollutionImage(img, mask image.Image) image.Image {
	bounds := img.Bounds()
	output := image.NewNRGBA64(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			currColor := rgbaToColorful(img.At(x, y))
			maskColor := rgbaToColorful(mask.At(x, y))

			currH, currS, currV := currColor.Hsv()
			maskH, maskS, maskV := maskColor.Hsv()
			output.Set(x, y, colorful.Hsv(currH-maskH, currS-maskS, currV-maskV).Clamped())
		}
	}

	return output
}

func LoadImages(images []string) []image.Image {
	var files []string
	for _, file := range images {
		files = append(files, collectFiles(file)...)
	}
	images = files

	loadedImages := make([]image.Image, len(images))

	var wg sync.WaitGroup
	for i := range images {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			loadedImages[i] = LoadImage(images[i])
		}(i)
	}

	wg.Wait()

	return loadedImages
}

func fileVisit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		ext := filepath.Ext(info.Name())
		if !info.IsDir() && (ext == ".jpg" || ext == ".jpeg" || ext == ".tif" || ext == ".tiff" || ext == ".png") {
			*files = append(*files, path)
		}

		return nil
	}
}

func collectFiles(file string) []string {
	var files []string
	filepath.Walk(file, fileVisit(&files))

	return files
}

func LoadImage(filename string) image.Image {
	var currImg *os.File
	var err error

	handleError := func(err error) {
		err = errors.Wrapf(err, "for image: %s", filename)
		panic(err)
	}

	if currImg, err = os.Open(filename); err != nil {
		handleError(err)
	}
	defer currImg.Close()

	decoded, _, err := image.Decode(currImg)
	if err != nil {
		handleError(err)
	}

	return decoded
}

// Single pixel denoising.
func DenoiseImage(img image.Image) image.Image {
	bounds := img.Bounds()
	output := image.NewNRGBA64(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			averageColor := getNeighborAverageColor(img, x, y)
			currColor := img.At(x, y)
			if distance(rgbaToColorful(currColor), rgbaToColorful(averageColor)) > delta {
				output.Set(x, y, averageColor)
			} else {
				output.Set(x, y, currColor)
			}
		}
	}

	return output
}

func getNeighborAverageColor(img image.Image, x, y int) color.Color {
	colors := make([]colorful.Color, 0)
	bounds := img.Bounds()

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if outOfBounds(x+i, y+j, bounds) || (i == 0 && j == 0) {
				continue
			}
			colors = append(colors, rgbaToColorful(img.At(x+i, y+j)))
		}
	}

	return AverageColor(colors)
}

func outOfBounds(x, y int, bounds image.Rectangle) bool {
	return x < bounds.Min.X || x > bounds.Max.X || y < bounds.Min.Y || y > bounds.Max.Y
}

func Upscale(images []image.Image) []image.Image {
	bounds := images[0].Bounds()
	width := bounds.Max.X * 2
	height := bounds.Max.Y * 2

	for i := range images {
		images[i] = imaging.Resize(images[i], width, height, imaging.Gaussian)
	}

	return images
}

func StarTrack(images []image.Image) []image.Image {
	reference := images[0]
	referenceMap := GetStarmap(reference)

	for i := 1; i < len(images); i++ {
		sMap := GetStarmap(images[i])
		config := sMap.FindOffset(referenceMap)
		fmt.Printf("%#v\n", config)
		images[i] = Transform(images[i], config)
	}

	return images
}

func Transform(img image.Image, config starmap.OffsetConfig) image.Image {
	img = Translate(img, config.X, config.Y)
	if config.Rotation != 0 {
		img = imaging.Rotate(img, config.Rotation, color.RGBA{R: 0, G: 0, B: 0, A: 0})
	}

	return img
}

func Translate(img image.Image, dx, dy int) *image.NRGBA64 {
	bounds := img.Bounds()
	output := image.NewNRGBA64(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x+dx < bounds.Max.X && y+dy < bounds.Max.Y {
				output.Set(x+dx, y+dy, img.At(x, y))
			}
		}
	}
	return output
}

func SaveImage(fileName string, image image.Image) error {
	f, _ := os.Create(fileName)
	defer f.Close()

	return tiff.Encode(f, image, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
}

func GetStarmap(img image.Image) starmap.Starmap {
	// Filter for brightness.
	bounds := img.Bounds()
	var brightPoints []image.Point

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, v := rgbaToColorful(img.At(x, y)).Hsv()
			// @todo dynamic range
			if v > 0.9 {
				brightPoints = append(brightPoints, image.Point{X: x, Y: y})
			}
		}
	}

	sm := starmap.Starmap{Bounds: bounds}
	for i := range brightPoints {
		sm.Stars = append(sm.Stars, starmap.Star{X: float64(brightPoints[i].X), Y: float64(brightPoints[i].Y), Size: 1})
	}

	return sm.Compress()
}

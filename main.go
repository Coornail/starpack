package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	"github.com/Coornail/superres/colr"
	"github.com/disintegration/imaging"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/errors"
	"golang.org/x/image/tiff"
)

const (
	motionCachePath = "/tmp/motion.json"

	delta = 0.1
)

var (
	supersample  bool
	verbose      bool
	fast         bool
	whiteBalance bool
	denoise      bool
	parallelism  int
	mergeMethod  string
	samplerName  string
	outputFile   string
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	flag.BoolVar(&supersample, "supersample", true, "Supersample image")
	flag.BoolVar(&verbose, "verbose", true, "Verbose output")
	flag.BoolVar(&fast, "fast", true, "Process images faster, trading quality")
	flag.BoolVar(&whiteBalance, "whiteBalance", false, "White balancing")
	flag.BoolVar(&denoise, "denoise", false, "Denoise input images")
	flag.IntVar(&parallelism, "parallelism", runtime.NumCPU()*2, "Number of threads to process images")
	flag.StringVar(&mergeMethod, "mergeMethod", "average", "Method to merge pixels from the input images (median, average, brightest)")
	flag.StringVar(&samplerName, "sampler", "combined", "Sample images for motion detection (gauss, uniform, edge)")
	flag.StringVar(&outputFile, "output", "output.tif", "Output file name")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	images := flag.Args()
	if fast {
		samplerName = "gauss"
		supersample = false
	}

	flag.VisitAll(func(f *flag.Flag) {
		verboseOutput("%s:\t%v\n", f.Name, f.Value)
	})

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	verboseOutput("Loading images\n")
	loadedImages := loadImages(images)

	if supersample {
		verboseOutput("Upscaling\n")
		loadedImages = upscale(loadedImages)
	}

	var colorMergeMethod ColorMerge = medianColor
	if mergeMethod == "average" {
		colorMergeMethod = averageColor
	} else if mergeMethod == "brightest" {
		colorMergeMethod = brightestColor
	}

	output := superres(loadedImages, colorMergeMethod)

	if whiteBalance {
		verboseOutput("White balancing\n")
		output = colr.ModifiedGrayWorld(output)
	}

	verboseOutput("Writing output\n")
	f, _ := os.Create(outputFile)
	defer f.Close()
	tiff.Encode(f, output, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
}

func superres(images []image.Image, colorMergeMethod ColorMerge) *image.NRGBA64 {
	bounds := images[0].Bounds()
	output := image.NewNRGBA64(bounds)

	var currentColor []colorful.Color
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			currentColor = make([]colorful.Color, len(images))

			wg.Add(1)
			go func(x, y int, currentColor []colorful.Color) {
				for i := range images {
					currX := x
					currY := y
					if currX < bounds.Min.X || currX >= bounds.Max.X ||
						currY < bounds.Min.Y || currY >= bounds.Max.Y {
						continue
					}

					currentColor[i] = rgbaToColorful(images[i].At(currX, currY))
				}
				defer wg.Done()
				output.Set(x, y, colorMergeMethod(currentColor))
			}(x, y, currentColor)
		}
		verboseOutput("Merging: %.2f%%\r", float64(y)/float64(bounds.Max.Y)*100.0)
	}
	wg.Wait()
	verboseOutput("\n")

	return output
}

func verboseOutput(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format, args...)

	}
}

func loadImages(images []string) []image.Image {
	loadedImages := make([]image.Image, len(images))

	var wg sync.WaitGroup
	for i := range images {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var currImg *os.File
			var err error

			if currImg, err = os.Open(images[i]); err != nil {
				err = errors.Wrapf(err, "for image: %s", currImg.Name())
				panic(err)
			}
			defer currImg.Close()

			decoded, _, err := image.Decode(currImg)
			if err != nil {
				err = errors.Wrapf(err, "for image: %s", currImg.Name())
				panic(err)
			}

			// @todo lift me up one level
			if denoise {
				loadedImages[i] = denoiseImage(decoded)
			} else {
				loadedImages[i] = decoded
			}
		}(i)
	}

	wg.Wait()

	return loadedImages
}

// Single pixel denoising.
func denoiseImage(img image.Image) image.Image {
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

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if outOfBounds(img, x+i, y+j) || (i == 0 && j == 0) {
				continue
			}
			colors = append(colors, rgbaToColorful(img.At(x+i, y+j)))
		}
	}

	return averageColor(colors)
}

func outOfBounds(img image.Image, x, y int) bool {
	return x < img.Bounds().Min.X || x > img.Bounds().Max.X || y < img.Bounds().Min.Y || y > img.Bounds().Max.Y
}

func upscale(images []image.Image) []image.Image {
	bounds := images[0].Bounds()
	width := bounds.Max.X * 2
	height := bounds.Max.Y * 2

	for i := range images {
		images[i] = imaging.Resize(images[i], width, height, imaging.Gaussian)
	}

	return images
}

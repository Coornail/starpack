package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/Coornail/starpack/colr"
	starpack "github.com/Coornail/starpack/lib"
)

var (
	supersample          bool
	verbose              bool
	whiteBalance         bool
	denoise              bool
	removeLightPollution bool
	mergeMethod          string
	outputFile           string
)

func verboseOutput(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format, args...)

	}
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	flag.BoolVar(&supersample, "supersample", false, "Supersample image")
	flag.BoolVar(&verbose, "verbose", true, "Verbose output")
	flag.BoolVar(&whiteBalance, "whiteBalance", false, "White balancing")
	flag.BoolVar(&denoise, "denoise", false, "Denoise input images")
	flag.BoolVar(&removeLightPollution, "removeLightPollution", true, "Remove light pollution")
	flag.StringVar(&mergeMethod, "mergeMethod", "average", "Method to merge pixels from the input images (median, average, brightest)")
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

	flag.VisitAll(func(f *flag.Flag) {
		verboseOutput("%s:\t%v\n", f.Name, f.Value)
	})

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	verboseOutput("Loading images\n")
	loadedImages := starpack.LoadImages(images)

	var wg sync.WaitGroup

	if denoise {
		verboseOutput("Denoising\n")
		for i := range loadedImages {
			wg.Add(1)
			go func(i int) {
				loadedImages[i] = starpack.DenoiseImage(loadedImages[i])
				wg.Done()
			}(i)
		}
		wg.Wait()
	}

	if removeLightPollution {
		verboseOutput("Removing light pollution\n")
		mask := starpack.EstimateLightPollutionMask(loadedImages[0])
		for i := range loadedImages {
			wg.Add(1)
			go func(i int) {
				loadedImages[i] = starpack.RemoveLightPollutionImage(loadedImages[i], mask)
				wg.Done()
			}(i)
		}

		wg.Wait()
	}

	if supersample {
		verboseOutput("Upscaling\n")
		loadedImages = starpack.Upscale(loadedImages)
	}

	var colorMergeMethod starpack.ColorMerge = starpack.MedianColor
	if mergeMethod == "average" {
		colorMergeMethod = starpack.AverageColor
	} else if mergeMethod == "brightest" {
		colorMergeMethod = starpack.BrightestColor
	}

	output := starpack.Starpack(loadedImages, colorMergeMethod)

	if whiteBalance {
		verboseOutput("White balancing\n")
		output = colr.ModifiedGrayWorld(output)
	}

	verboseOutput("Writing output\n")
	starpack.SaveImage(outputFile, output)
}

package sampler

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"testing"
)

func TestUniformSampler(t *testing.T) {
	bounds := image.Rect(0, 0, 1024, 768)
	output := image.NewNRGBA(bounds)
	us := NewUniformSampler(output, 1024)

	renderSampler(us, output, "uniform.png")
}

/*
func TestRandomSampler(t *testing.T) {
	bounds := image.Rect(0, 0, 1024, 768)
	output := image.NewNRGBA(bounds)
	us := NewRandomSampler(output, 32)

	renderSampler(us, output, "random.png")
}
*/

func TestGaussSampler(t *testing.T) {
	bounds := image.Rect(0, 0, 1024, 768)
	output := image.NewNRGBA(bounds)
	us := NewGaussSampler(output, 1024)

	renderSampler(us, output, "gauss.png")
}

func TestEdgeDetector(t *testing.T) {
	file, _ := os.Open("edge_target.jpg")
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	output := image.NewNRGBA(img.Bounds())
	ed := NewEdgeDetector(img, 10240)

	renderSampler(ed, output, "edge-detector.png")
	fmt.Printf("%d edges\n", ed.Edges)
}

func renderSampler(sampler ImageSampler, output *image.NRGBA, fileName string) {
	bounds := output.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			output.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}

	for sampler.HasMore() {
		x, y := sampler.Next()
		/*
			output.Set(x-1, y-1, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			output.Set(x-1, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			output.Set(x-1, y+1, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		*/
		output.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	}

	f, _ := os.Create(fileName)
	defer f.Close()
	png.Encode(f, output)
}

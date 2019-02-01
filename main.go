package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/tracer"
	"github.com/DheerendraRathor/GoTracer/utils"
	"gopkg.in/cheggaaa/pb.v1"
)

var renderSpecFile string
var doCpuProfile bool
var showProgress bool

func init() {
	flag.StringVar(&renderSpecFile, "spec", "./examples/fiveSpheresWithLights.json", "Name of JSON file containing rendering spec")
	flag.BoolVar(&doCpuProfile, "cpu", false, "Enable CPU Profile")
	flag.BoolVar(&showProgress, "progress", false, "Show progress by rendering pixel by pixel")
}

func main() {
	flag.Parse()

	// Profiler Code
	if doCpuProfile {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	file, e := ioutil.ReadFile(renderSpecFile)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}

	var env models.Specification
	json.Unmarshal(file, &env)

	pngImage := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{env.Image.Width, env.Image.Height},
	})

	if showProgress {
		progress := make(chan *models.Pixel, 100)
		defer close(progress)

		var pbWg sync.WaitGroup
		var progressBar *pb.ProgressBar

		pbWg.Add(1)
		go func() {
			defer pbWg.Done()

			total := env.Image.Width * env.Image.Height
			progressBar = pb.StartNew(total)
			progressBar.ShowFinalTime = true
			progressBar.ShowTimeLeft = false
			progressBar.ShowBar = false

			for pixel := range progress {
				if pixel == nil {
					break
				}

				updateImage(pngImage, pixel)
				progressBar.Increment()
			}
		}()

		tracer.GoTrace(&env, true, progress, false, nil)

		pbWg.Wait()
		progressBar.Finish()
	} else {
		tracerOutput := tracer.GoTrace(&env, false, nil, false, nil)

		for _, row := range tracerOutput.Pixels {
			for _, pixel := range row {
				updateImage(pngImage, pixel)
			}
		}

	}

	pngFile := utils.CreateNestedFile(env.Image.OutputFile)
	defer pngFile.Close()

	png.Encode(pngFile, pngImage)

}

func updateImage(image *image.RGBA, pixel *models.Pixel) {
	rgbaColor := color.RGBA{pixel.Color[0], pixel.Color[1], pixel.Color[2], 255}
	image.Set(pixel.I, pixel.J, rgbaColor)
}

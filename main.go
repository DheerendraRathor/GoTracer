package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"sync"

	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/tracer"
	"github.com/DheerendraRathor/GoTracer/utils"
	"gopkg.in/cheggaaa/pb.v1"
)

var renderSpecFile string

func init() {
	flag.StringVar(&renderSpecFile, "spec", "sample_world.json", "Name of JSON file containing rendering spec")
}

func main() {
	flag.Parse()

	file, e := ioutil.ReadFile(renderSpecFile)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}

	var env models.World
	json.Unmarshal(file, &env)

	progress := make(chan *models.Pixel, 100)
	defer close(progress)

	pngImage := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{env.Image.Width, env.Image.Height},
	})

	var pbWg sync.WaitGroup
	var progressBar *pb.ProgressBar

	pbWg.Add(1)
	go func() {
		defer pbWg.Done()
		if env.Settings.ShowProgress {
			total := env.Image.Width * env.Image.Height
			progressBar = pb.StartNew(total)
			progressBar.ShowFinalTime = true
			progressBar.ShowTimeLeft = false
		}
		for pixel := range progress {
			if pixel == nil {
				break
			}

			rgbaColor := color.RGBA{pixel.Color[0], pixel.Color[1], pixel.Color[2], 255}
			pngImage.Set(pixel.I, pixel.J, rgbaColor)

			if env.Settings.ShowProgress {
				progressBar.Increment()
			}
		}
	}()

	closeChan := make(chan bool)
	goTracer.GoTrace(&env, progress, closeChan)

	pngFile := utils.CreateNestedFile(env.Image.OutputFile)
	defer pngFile.Close()

	pbWg.Wait()

	png.Encode(pngFile, pngImage)

}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"io/ioutil"
	"math"
	"sync"

	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/tracer"
	"github.com/DheerendraRathor/GoTracer/utils"
	"gopkg.in/cheggaaa/pb.v1"
)

var renderSpecFile string
var showProgress bool

func init() {
	flag.StringVar(&renderSpecFile, "spec", "./examples/dolly.json", "Name of JSON file containing rendering spec")
	flag.BoolVar(&showProgress, "progress", false, "Show progress by rendering pixel by pixel")
}

func main() {
	flag.Parse()

	file, e := ioutil.ReadFile(renderSpecFile)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}

	var env models.Specification
	json.Unmarshal(file, &env)

	outputFileFormat := env.Image.OutputFile

	cameraInput := env.Scene.Camera

	constHalfHeight := math.Tan(cameraInput.FieldOfView*math.Pi/360) * cameraInput.Focus
	initialZDistance := cameraInput.LookFrom[2]
	initialFieldOfView := cameraInput.FieldOfView
	initialFocus := cameraInput.Focus

	for zoomMode := 0; zoomMode < 2; zoomMode++ {
		directoryName := "out"
		distanceChange := -0.05
		if zoomMode%2 == 1 {
			directoryName = "in"
			distanceChange = 0.05
		}
		cameraInput.LookFrom[2] = initialZDistance
		cameraInput.FieldOfView = initialFieldOfView
		cameraInput.Focus = initialFocus

		outGif := &gif.GIF{}
		env.Image.OutputFile = fmt.Sprintf(outputFileFormat, directoryName)

		for i := 0; i < 40; i++ {
			fmt.Printf("Rendering frame: %d. ZoomMode: Zoom %s\n", i, directoryName)
			env.Scene.Camera = cameraInput

			progress := make(chan *models.Pixel, 1000)

			imageRect := image.Rectangle{
				Min: image.Point{},
				Max: image.Point{X: env.Image.Width, Y: env.Image.Height},
			}

			palleted := image.NewPaletted(imageRect, palette.Plan9)

			var pbWg sync.WaitGroup
			var progressBar *pb.ProgressBar

			pbWg.Add(1)
			go func() {
				defer pbWg.Done()
				if showProgress {
					total := env.Image.Width * env.Image.Height
					progressBar = pb.StartNew(total)
					progressBar.ShowFinalTime = true
					progressBar.ShowTimeLeft = false
					progressBar.ShowBar = false
				}
				for pixel := range progress {
					if pixel == nil {
						break
					}

					rgbaColor := color.RGBA{R: pixel.Color[0], G: pixel.Color[1], B: pixel.Color[2], A: 255}
					palleted.Set(pixel.I, pixel.J, rgbaColor)

					if showProgress {
						progressBar.Increment()
					}
				}
			}()

			tracer.GoTrace(&env, true, progress, false, nil)

			pbWg.Wait()
			close(progress)
			progressBar.Finish()

			outGif.Image = append(outGif.Image, palleted)
			outGif.Delay = append(outGif.Delay, 0)

			// Changing camera distance and maintaining FoV
			cameraInput.LookFrom[2] += distanceChange
			newFocus := models.NewVectorFromArray(cameraInput.LookFrom).
				SubtractVector(models.NewVectorFromArray(cameraInput.LookAt)).
				Length()
			cameraInput.FieldOfView = math.Atan(constHalfHeight/newFocus) * 360 / math.Pi
			cameraInput.Focus = newFocus
		}

		gifFile := utils.CreateNestedFile(env.Image.OutputFile)
		gif.EncodeAll(gifFile, outGif)
		gifFile.Close()
	}
}

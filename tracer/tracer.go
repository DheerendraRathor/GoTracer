package tracer

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/DheerendraRathor/GoTracer/models"
)

var MaxRenderDepth int = 10

type TracerOutput struct {
	Pixels [][]*models.Pixel
}

func GoTrace(
	env *models.Specification,
	sharePixelProgress bool, progress chan<- *models.Pixel,
	isClosable bool, closeChan <-chan bool,
) *TracerOutput {
	if env.Settings.RenderDepth > 0 {
		MaxRenderDepth = env.Settings.RenderDepth
	}

	scene := env.GetScene()

	width, height := env.Image.Width, env.Image.Height

	askedRenderRoutines := env.Settings.RenderRoutines
	if askedRenderRoutines <= 0 {
		askedRenderRoutines = runtime.NumCPU()
	}

	renderRoutines := askedRenderRoutines - 2
	if renderRoutines < 1 {
		renderRoutines = 1
	}

	processingGroupData := make([][][2]int, renderRoutines)
	for i := range processingGroupData {
		processingGroupData[i] = make([][2]int, 2)
	}

	imin, imax, jmin, jmax := env.Image.GetPatch()

	var output *TracerOutput
	if !sharePixelProgress {
		output = &TracerOutput{
			Pixels: make([][]*models.Pixel, imax),
		}
		for i := range output.Pixels {
			output.Pixels[i] = make([]*models.Pixel, jmax)
		}
	}

	division := 0
	processingGroup := 0
	for i := imax - 1; i >= imin; i-- {
		for j := jmin; j < jmax; j++ {

			processingGroup = division % renderRoutines
			processingGroupData[processingGroup] = append(processingGroupData[processingGroup], [2]int{i, j})
			division += 1

			/*
				// Check if go tracer is asked for close
				if isClosable {
					select {
					case <-closeChan:
						break IMAGE_PROCESSING
					default:
					}
				}
			*/

		}
	}

	var renderWg sync.WaitGroup
	renderWg.Add(renderRoutines)
	for _, _data := range processingGroupData {
		go func(samples int, scene *models.Scene, data [][2]int) {

			defer func() {
				renderWg.Done()
			}()

			rng := rand.New(rand.NewSource(time.Now().Unix() + rand.Int63()))

			for _, point := range data {
				i, j := point[0], point[1]
				pixel := processPixel(i, j, width, height, samples, scene, rng)
				if sharePixelProgress {
					progress <- pixel
				} else {
					output.Pixels[i][j] = pixel
				}
			}
		}(env.Image.Samples, scene, _data)
	}

	renderWg.Wait()

	if sharePixelProgress {
		progress <- nil
	}

	return output
}

func processPixel(i, j, imageWidth, imageHeight, sample int, scene *models.Scene, rng *rand.Rand) *models.Pixel {
	pixel := models.NewEmptyVector()
	for s := 0; s < sample; s++ {
		randFloatu, randFloatv := rng.Float64(), rng.Float64()
		u, v := (float64(j)+randFloatu)/float64(imageWidth), (float64(i)+randFloatv)/float64(imageHeight)
		ray := scene.Camera.RayAt(u, v, rng)
		pixel.AddVector(getColor(ray, scene, 0, rng))
	}

	pixel.Scale(1 / float64(sample)).Gamma2()
	uint8Pixel := pixel.ToPixel(j, imageHeight-i-1)
	return uint8Pixel
}

func getColor(r *models.Ray, scene *models.Scene, renderDepth int, rng *rand.Rand) *models.Vector {

	// tmin is 0.0001 to avoid self intersection
	didHit, hitRecord := scene.HitableList.Hit(r, 0.0001, math.MaxFloat64)
	if didHit {
		shouldScatter, attenuation, ray := hitRecord.Material.Scatter(r, hitRecord, rng)

		if hitRecord.Material.IsLight() {
			return attenuation
		}

		if renderDepth < MaxRenderDepth && shouldScatter {
			return attenuation.MultiplyVector(getColor(ray, scene, renderDepth+1, rng))
		} else {
			return models.NewEmptyVector()
		}
	}

	return scene.AmbientLight
}

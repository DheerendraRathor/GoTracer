package goTracer

import (
	"math"
	"math/rand"
	"runtime"
	"sync"

	"github.com/DheerendraRathor/GoTracer/models"
)

var MaxRenderDepth int = 10

func GoTrace(env *models.World, progress chan<- *models.Pixel, closeChan <-chan bool) {
	if env.Settings.RenderDepth > 0 {
		MaxRenderDepth = env.Settings.RenderDepth
	}

	camera := env.GetCamera()
	world := env.GetHitableList()

	width, height := env.Image.Width, env.Image.Height

	var renderWg sync.WaitGroup

	renderRoutines := env.Settings.RenderRoutines
	if renderRoutines <= 0 {
		renderRoutines = runtime.NumCPU()
	}
	renderer := make(chan bool, renderRoutines)
	defer close(renderer)

	imin, imax, jmin, jmax := env.Image.GetPatch()

IMAGE_PROCESSING:
	for i := imax - 1; i >= imin; i-- {
		for j := jmin; j < jmax; j++ {

			// Check if go tracer is asked for close
			select {
			case <-closeChan:
				break IMAGE_PROCESSING
			default:
			}

			renderer <- true
			renderWg.Add(1)

			go func(i, j, samples int, camera *models.Camera, world *models.HitableList) {
				defer func() {
					<-renderer
					renderWg.Done()
				}()
				pixel := processPixel(i, j, width, height, samples, camera, world)
				progress <- pixel
			}(i, j, env.Image.Samples, camera, world)
		}
	}
	renderWg.Wait()

	progress <- nil
}

func processPixel(i, j, imageWidth, imageHeight, sample int, camera *models.Camera, world *models.HitableList) *models.Pixel {
	pixel := models.Vector{0, 0, 0}
	for s := 0; s < sample; s++ {
		randFloatu, randFloatv := rand.Float64(), rand.Float64()
		u, v := (float64(j)+randFloatu)/float64(imageWidth), (float64(i)+randFloatv)/float64(imageHeight)
		ray := camera.RayAt(u, v)
		pixel.Add(getColor(ray, world, 0))
	}

	pixel.DivideByScalar(float64(sample)).Gamma2()
	uint8Pixel := pixel.ToUint8(j, imageHeight-i-1)
	return uint8Pixel
}

func getColor(r *models.Ray, world *models.HitableList, renderDepth int) models.Vector {

	// tmin is 0.0001 to avoid self intersection
	willHit, hitRecord := world.Hit(r, 0.0001, math.MaxFloat64)
	if willHit {
		shouldScatter, attenuation, ray := hitRecord.Material.Scatter(r, hitRecord)
		if renderDepth < MaxRenderDepth && shouldScatter {
			return models.MultiplyVectors(attenuation, getColor(ray, world, renderDepth+1))
		} else {
			return []float64{0, 0, 0}
		}
	}

	unitDir := models.UnitVector(r.Direction)
	t := 0.5 * (unitDir.Y() + 1.0)
	var startBlend, endBlend models.Vector
	startBlend = models.Vector{1.0, 1.0, 1.0}.MultiplyScalar(1 - t)
	endBlend = models.Vector{0.5, 0.7, 1.0}.MultiplyScalar(t)

	return startBlend.Add(endBlend)
}

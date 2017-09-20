package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/utils"
	"gopkg.in/cheggaaa/pb.v1"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"runtime"
	"sync"
)

var MaxRenderDepth int = 10

var renderSpecFile string

func init() {
	flag.StringVar(&renderSpecFile, "spec", "sample_world.json", "Name of JSON file containing rendering spec")
}

func main() {

	file, e := ioutil.ReadFile(renderSpecFile)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}

	var env models.World
	json.Unmarshal(file, &env)

	if env.Settings.RenderDepth > 0 {
		MaxRenderDepth = env.Settings.RenderDepth
	}

	camera := env.GetCamera()

	rows, columns := env.Image.Rows, env.Image.Columns
	sample := env.Image.Samples

	world := env.GetHitableList()

	progress := make(chan bool, 100)
	defer close(progress)

	var pbWg, renderWg sync.WaitGroup

	if env.Settings.ShowProgress {
		// Progress Bar
		pbWg.Add(1)
		go func() {
			defer pbWg.Done()
			total := rows * columns
			bar := pb.StartNew(total)
			bar.ShowFinalTime = true
			bar.ShowTimeLeft = false
			for value := range progress {
				if value {
					break
				}
				bar.Increment()
			}
		}()
	}

	pngImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{columns, rows}})
	pngFile := utils.CreateNestedFile(env.Image.OutputFile)
	defer pngFile.Close()

	renderRoutines := env.Settings.RenderRoutines
	if renderRoutines <= 0 {
		renderRoutines = runtime.NumCPU()
	}
	renderer := make(chan bool, renderRoutines)
	defer close(renderer)

	for i := env.Image.IMax - 1; i >= env.Image.IMin; i-- {
		for j := env.Image.JMin; j < env.Image.JMax; j++ {
			renderer <- true
			renderWg.Add(1)
			go func(i, j int) {
				defer func() {
					<-renderer
					renderWg.Done()
				}()
				ProcessPixel(i, j, rows, columns, sample, &camera, &world, pngImage)
				if env.Settings.ShowProgress {
					progress <- false
				}
			}(i, j)
		}
	}
	renderWg.Wait()

	if env.Settings.ShowProgress {
		progress <- true
		pbWg.Wait()
	}

	png.Encode(pngFile, pngImage)
}

func ProcessPixel(i, j, rows, columns, sample int, camera *models.Camera, world *models.HitableList, pngImage *image.RGBA) {
	colorVector := models.NewVector3D(0, 0, 0)
	for s := 0; s < sample; s++ {
		randFloatu, randFloatv := rand.Float64(), rand.Float64()
		u, v := (float64(j)+randFloatu)/float64(columns), (float64(i)+randFloatv)/float64(rows)
		ray := camera.RayAt(u, v)
		colorVector = models.AddVectors(colorVector, Color(ray, *world, 0))
	}

	pixel := models.NewPixelFromVector(
		models.DivideScalar(colorVector, float64(sample)),
	)
	pixel.Gamma2()
	uint8Pixel := pixel.UInt8Pixel()
	rgba := color.RGBA{uint8Pixel.R, uint8Pixel.G, uint8Pixel.B, 255}
	pngImage.Set(j, rows-i-1, rgba)
}

func Color(r models.Ray, world models.HitableList, renderDepth int) models.Pixel {

	willHit, hitRecord := world.Hit(r, 0.0, math.MaxFloat64)
	if willHit {
		shouldScatter, attenuation, ray := hitRecord.Material.Scatter(r, hitRecord)
		if renderDepth < MaxRenderDepth && shouldScatter {
			colorVector := models.MultiplyVectors(attenuation, Color(ray, world, renderDepth+1))
			return models.NewPixelFromVector(colorVector)
		}
	}

	var unitDir models.Vector3D = models.UnitVector(r.Direction)
	t := 0.5 * (unitDir.Y() + 1.0)
	var startValue, endValue models.Vector3D
	startValue = models.NewVector3D(1.0, 1.0, 1.0)
	endValue = models.NewVector3D(0.5, 0.7, 1.0)
	var startBlend models.Vector3D = models.MultiplyScalar(startValue, 1-t)
	var endBlend models.Vector3D = models.MultiplyScalar(endValue, t)
	return models.Pixel{
		Vector3D: models.AddVectors(startBlend, endBlend),
	}

}

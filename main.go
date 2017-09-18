package main

import (
	"github.com/DheerendraRathor/GoTracer/models"
	"gopkg.in/cheggaaa/pb.v1"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
)

const (
	MaxDepth = 10
)

func main() {
	lookFrom := models.NewPoint(-3, 1, 0.5)
	lookAt := models.NewPoint(0, 0, -0.5)
	vUp := models.NewVector3D(0, 1, 0)
	camera := models.NewCamera(lookFrom, lookAt, vUp, 45, 2)

	rows, columns := 200, 400
	sample := 10

	world := models.HitableList{}

	world.AddHitable(models.NewSphere(0, 0, -1, 0.5, models.NewLambertian(0.8, 0.3, 0.3)))
	world.AddHitable(models.NewSphere(0, -1000.5, -1, 1000, models.NewLambertian(0.8, 0.8, 0)))
	world.AddHitable(models.NewSphere(1, 0, -1, 0.5, models.NewMetal(0.8, 0.6, 0.2, 0.2)))
	world.AddHitable(models.NewSphere(-1, 0, -1, 0.5, models.NewDielectric(1.3)))
	world.AddHitable(models.NewSphere(-1, 0, -1.75, 0.25, models.NewLambertian(0.2, 0.2, 0.7)))

	progress := make(chan bool, 100)
	var pbWg, renderWg sync.WaitGroup

	// Progress Bar
	pbWg.Add(1)
	go func() {
		defer pbWg.Done()
		total := rows * columns
		bar := pb.StartNew(total)
		for value := range progress {
			if value {
				break
			}
			bar.Increment()
		}
	}()

	pngImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{columns, rows}})
	pngFile, _ := os.Create("myTestImage.png")
	defer pngFile.Close()

	renderer := make(chan bool, 4)
	for i := rows - 1; i >= 0; i-- {
		for j := 0; j < columns; j++ {
			renderer <- true
			renderWg.Add(1)
			go func(i, j int) {
				defer func() {
					<-renderer
					renderWg.Done()
				}()
				ProcessPixel(i, j, rows, columns, sample, &camera, &world, pngImage)
				progress <- false
			}(i, j)
		}
	}
	renderWg.Wait()

	progress <- true
	pbWg.Wait()

	png.Encode(pngFile, pngImage)
}

func ProcessPixel(i, j, rows, columns, sample int, camera *models.Camera, world *models.HitableList, pngImage *image.RGBA) {
	pixel := models.NewPixel(0, 0, 0)
	for s := 0; s < sample; s++ {
		randFloatu, randFloatv := rand.Float64(), rand.Float64()
		u, v := (float64(j)+randFloatu)/float64(columns), (float64(i)+randFloatv)/float64(rows)
		ray := camera.RayAt(u, v)
		pixel = models.NewPixelFromVector(
			models.AddVectors(pixel, Color(ray, *world, 0)),
		)
	}
	pixel = models.NewPixelFromVector(
		models.DivideScalar(pixel, float64(sample)),
	)
	pixel.Gamma2()
	uint8Pixel := pixel.UInt8Pixel()
	rgba := color.RGBA{uint8Pixel.R, uint8Pixel.G, uint8Pixel.B, 255}
	pngImage.Set(j, rows-i-1, rgba)
}

func Color(r models.Ray, world models.HitableList, depth int) models.Pixel {

	willHit, hitRecord := world.Hit(r, 0.0, math.MaxFloat64)
	if willHit {
		shouldScatter, attenuation, ray := hitRecord.Material.Scatter(r, hitRecord)
		if depth < MaxDepth && shouldScatter {
			colorVector := models.MultiplyVectors(attenuation, Color(ray, world, depth+1))
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

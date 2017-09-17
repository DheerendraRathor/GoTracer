package main

import (
	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/utils"
	"math"
	"math/rand"
)

const (
	MaxDepth = 10
)

func main() {
	camera := models.NewCamera()
	camera.LowerLeftCorner = models.NewPoint(-4.0, -2.0, -2.0)
	camera.Horizontal = models.NewVector3D(8, 0, 0)
	camera.Vertical = models.NewVector3D(0, 4, 0)
	camera.Origin = models.NewPoint(0, 0.5, 2)

	rows, columns := 400, 800
	sample := 100

	image := make([][]models.Pixel, rows)
	for i := 0; i < rows; i++ {
		image[i] = make([]models.Pixel, columns)
	}

	world := models.HitableList{}

	world.AddHitable(models.NewSphere(0, 0, -1, 0.5, models.NewLambertian(0.8, 0.3, 0.3)))
	world.AddHitable(models.NewSphere(0, -100.5, -1, 100, models.NewLambertian(0.8, 0.8, 0)))
	world.AddHitable(models.NewSphere(1, 0, -1, 0.5, models.NewMetal(0.8, 0.6, 0.2, 0.6)))
	world.AddHitable(models.NewSphere(-1, 0, -1, 0.5, models.NewMetal(0.8, 0.8, 0.8, 0.1)))

	for i := rows - 1; i >= 0; i-- {
		for j := 0; j < columns; j++ {
			color := models.NewPixel(0, 0, 0)
			for s := 0; s < sample; s++ {
				randFloatu, randFloatv := rand.Float64(), rand.Float64()
				u, v := (float64(j)+randFloatu)/float64(columns), (float64(i)+randFloatv)/float64(rows)
				ray := camera.RayAt(u, v)
				color = models.NewPixelFromVector(
					models.AddVectors(color, Color(ray, world, 0)),
				)
			}
			color = models.NewPixelFromVector(
				models.DivideScalar(color, float64(sample)),
			)
			color.Gamma2()
			image[rows-i-1][j] = color
		}
	}

	utils.RenderPPM(image, "myTestImage2.ppm")
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

package main

import (
	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/utils"
	"math"
	"math/rand"
)

var sphere1 models.Sphere = models.NewSphere(0, 0, -1, 0.5)

var sphere2 models.Sphere = models.NewSphere(0, -100.5, -1, 100)

func main() {
	camera := models.NewCamera()
	rows, columns := 100, 200
	sample := 10

	image := make([][]models.Pixel, rows)
	for i := 0; i < rows; i++ {
		image[i] = make([]models.Pixel, columns)
	}

	world := models.HitableList{}

	world.List = append(world.List, sphere1)
	world.List = append(world.List, sphere2)

	for i := rows - 1; i >= 0; i-- {
		for j := 0; j < columns; j++ {
			color := models.NewPixel(0, 0, 0)
			for s := 0; s < sample; s++ {
				randFloatu, randFloatv := rand.Float64(), rand.Float64()
				u, v := (float64(j)+randFloatu)/float64(columns), (float64(i)+randFloatv)/float64(rows)
				ray := camera.RayAt(u, v)
				color = models.NewPixelFromVector(
					models.AddVectors(color, Color(ray, world)),
				)
			}
			color = models.NewPixelFromVector(
				models.DivideScalar(color, float64(sample)),
			)
			color.Gamma2()
			image[rows-i-1][j] = color
		}
	}

	utils.RenderPPM(image, "myTestImage.ppm")
}

func Color(r models.Ray, world models.HitableList) models.Pixel {

	willHit, hitRecord := world.Hit(r, 0.0, math.MaxFloat64)
	if willHit {
		pN := models.AddVectors(hitRecord.P, hitRecord.N)
		targetPoint := models.AddVectors(pN, utils.RandomPointInUnitSphere())
		rayToTargetPoint := models.Ray{
			Origin:    hitRecord.P,
			Direction: models.SubtractVectors(targetPoint, hitRecord.P),
		}
		outputPixelVector := models.MultiplyScalar(Color(rayToTargetPoint, world), 0.5)
		return models.NewPixelFromVector(outputPixelVector)
		//pixel := models.MultiplyScalar(models.NewPixel(hitRecord.N.X()+1, hitRecord.N.Y()+1, hitRecord.N.Z()+1), 0.5)
		//return models.NewPixelFromVector(pixel)
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

package main

import (
	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/utils"
	"math"
)

var sphere1 models.Sphere = models.Sphere{
	models.NewPoint(0, 0, -1),
	0.5,
}

var sphere2 models.Sphere = models.Sphere{
	models.NewPoint(0, -100.5, -1),
	100,
}

func main() {
	var lowerLeftCorner, horizontal, vertical models.Vector3D
	lowerLeftCorner = models.NewVector3D(-2.0, -1.0, -1.0)
	horizontal = models.NewVector3D(4.0, 0.0, 0.0)
	vertical = models.NewVector3D(0.0, 2.0, 0.0)

	var origin models.Point = models.NewPoint(0.0, 0.0, 0.0)

	rows, columns := 100, 200

	image := make([][]models.Pixel, rows)
	for i := 0; i < rows; i++ {
		image[i] = make([]models.Pixel, columns)
	}

	world := models.HitableList{}

	world.List = append(world.List, sphere1)
	world.List = append(world.List, sphere2)

	for i := rows - 1; i >= 0; i-- {
		for j := 0; j < columns; j++ {
			u, v := float64(j)/float64(columns), float64(i)/float64(rows)
			var r models.Ray
			horizontalDir := models.AddVectors(lowerLeftCorner, models.MultiplyScalar(horizontal, u))
			compositeDir := models.AddVectors(horizontalDir, models.MultiplyScalar(vertical, v))
			r = models.Ray{
				origin,
				compositeDir,
			}
			image[rows-i-1][j] = Color(r, world)
		}
	}

	utils.RenderPPM(image, "myTestImage.ppm")
}

func Color(r models.Ray, world models.HitableList) models.Pixel {

	willHit, hitRecord := world.Hit(r, 0.0, math.MaxFloat64)
	if willHit {
		pixel := models.MultiplyScalar(models.NewPixel(hitRecord.N.X()+1, hitRecord.N.Y()+1, hitRecord.N.Z()+1), 0.5)
		return models.NewPixelFromVector(pixel)
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

func RenderHelloWorld() {
	nx, ny := 200, 100

	image := make([][]models.Pixel, ny)
	for i := 0; i < ny; i++ {
		image[i] = make([]models.Pixel, nx)
	}

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			r, g, b := float64(i)/float64(nx), float64(j)/float64(ny), float64(0.2)
			image[j][i] = models.NewPixel(r, g, b)
		}
	}

	utils.RenderPPM(image, "myTestImage.ppm")
}

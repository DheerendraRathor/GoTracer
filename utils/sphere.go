package utils

import (
	"github.com/DheerendraRathor/GoTracer/models"
	"math/rand"
)

func RandomPointInUnitSphere() models.Point {
	var p models.Point
	for {
		x, y, z := 2*rand.Float64()-1, 2*rand.Float64()-1, 2*rand.Float64()-1
		p = models.NewPoint(x, y, z)
		if models.VectorDotProduct(p, p) < 1.0 {
			break
		}
	}

	return p
}

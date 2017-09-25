package models

import (
	"math"
	"math/rand"
)

type Camera struct {
	LowerLeftCorner, Origin Vector
	Horizontal, Vertical    Vector
	LensRadius              float64
	U, V, W                 Vector
}

func (c *Camera) RayAt(u, v float64) *Ray {

	rd := MultiplyScalar(RandomPointInUnitDisk(), c.LensRadius)
	offset := AddVectors(MultiplyScalar(c.U, rd[0]), MultiplyScalar(c.V, rd[1]))
	origin := AddVectors(c.Origin, offset)

	compositeDir := AddVectors(c.LowerLeftCorner, MultiplyScalar(c.Horizontal, u)).
		Add(MultiplyScalar(c.Vertical, v)).
		Subtract(origin)
	return &Ray{
		origin,
		compositeDir,
	}
}

func NewCamera(lookFrom, lookAt, vup Vector, vfov, aspect, aperture, focus float64) *Camera {
	theta := vfov * math.Pi / 180
	half_height := math.Tan(theta / 2)

	half_width := aspect * half_height

	w := UnitVector(SubtractVectors(lookFrom, lookAt))
	u := UnitVector(VectorCrossProduct(vup, w))
	v := VectorCrossProduct(w, u)

	llc := SubtractVectors(lookFrom, MultiplyScalar(u, half_width*focus)).
		Subtract(MultiplyScalar(v, half_height*focus)).
		Subtract(MultiplyScalar(w, focus))

	camera := &Camera{
		LowerLeftCorner: llc,
		Horizontal:      MultiplyScalar(u, 2*half_width*focus),
		Vertical:        MultiplyScalar(v, 2*half_height*focus),
		Origin:          lookFrom,
		LensRadius:      aperture / 2,
		U:               u,
		V:               v,
		W:               w,
	}

	return camera
}

func RandomPointInUnitDisk() Vector {
	var p Vector
	origin := []float64{1, 1, 0}
	for {
		p = []float64{rand.Float64(), rand.Float64(), 0}
		p.MultiplyScalar(2).Subtract(origin)

		if VectorDotProduct(p, p) < 1.0 {
			break
		}
	}

	return p
}

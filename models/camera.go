package models

import (
	"math"
	"math/rand"
)

type Camera struct {
	LowerLeftCorner, Origin *Vector
	Horizontal, Vertical    *Vector
	LensRadius              float64
	U, V, W                 *Vector
}

func (c *Camera) RayAt(u, v float64, rng *rand.Rand) *Ray {

	rd := RandomPointInUnitDisk(rng).Scale(c.LensRadius)
	origin := c.Origin.Copy().
		AddScaledVector(c.U, rd.X()).
		AddScaledVector(c.V, rd.Y())

	compositeDir := c.LowerLeftCorner.Copy().
		AddScaledVector(c.Horizontal, u).
		AddScaledVector(c.Vertical, v).
		SubtractVector(origin)

	return &Ray{
		origin,
		compositeDir,
	}
}

func NewCamera(lookFrom, lookAt, vup *Vector, vfov, aspect, aperture, focus float64) *Camera {
	theta := vfov * math.Pi / 180
	half_height := math.Tan(theta / 2)

	half_width := aspect * half_height

	w := lookFrom.Copy().SubtractVector(lookAt).MakeUnitVector()
	u := NewEmptyVector().VectorCrossProduct(vup, w).MakeUnitVector()
	v := NewEmptyVector().VectorCrossProduct(w, u)

	llc := lookFrom.Copy().SubtractVector(u.Copy().Scale(half_width * focus)).
		SubtractVector(v.Copy().Scale(half_height * focus)).
		SubtractVector(w.Copy().Scale(focus))

	camera := &Camera{
		LowerLeftCorner: llc,
		Horizontal:      u.Copy().Scale(2 * half_width * focus),
		Vertical:        v.Copy().Scale(2 * half_height * focus),
		Origin:          lookFrom,
		LensRadius:      aperture / 2,
		U:               u,
		V:               v,
		W:               w,
	}

	return camera
}

func RandomPointInUnitDisk(rng *rand.Rand) *Vector {
	var p = NewEmptyVector()
	var x, y float64
	for {
		x, y = rng.Float64(), rng.Float64()
		p.Update(2*x-1, 2*y-1, 0)

		if p.Dot(p) < 1.0 {
			break
		}
	}

	return p
}

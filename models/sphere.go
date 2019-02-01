package models

import (
	"math"
	"math/rand"
)

type Sphere struct {
	Center   *Vector
	Radius   float64
	Material Material
}

func NewSphere(x, y, z, r float64, material Material) *Sphere {
	return &Sphere{
		NewVector(x, y, z),
		r,
		material,
	}
}

func (s *Sphere) Hit(r *Ray, tmin, tmax float64) (bool, *HitRecord) {
	oc := r.Origin.Copy().SubtractVector(s.Center)
	var a, b, c, d float64
	a = r.Direction.Dot(r.Direction)
	b = 2.0 * oc.Dot(r.Direction)
	c = oc.SquaredLength() - s.Radius*s.Radius
	d = b*b - 4*a*c

	if d > 0 {
		sqrtD := math.Sqrt(d)
		a2 := 2 * a
		root := (-b - sqrtD) / a2
		if root > tmin && root < tmax {

			tempP := r.PointAtParameter(root)

			record := &HitRecord{
				T:        root,
				P:        tempP,
				N:        tempP.Copy().SubtractVector(s.Center).MakeUnitVector(),
				Material: s.Material,
			}

			return true, record
		}
		root = (-b + sqrtD) / a2
		if root > tmin && root < tmax {
			tempP := r.PointAtParameter(root)

			record := &HitRecord{
				T:        root,
				P:        tempP,
				N:        tempP.Copy().SubtractVector(s.Center).MakeUnitVector(),
				Material: s.Material,
			}

			return true, record
		}
	}
	return false, nil
}

func RandomPointInUnitSphere(rng *rand.Rand) *Vector {
	p := NewEmptyVector()
	var x, y, z float64
	for {
		x, y, z = rng.Float64(), rng.Float64(), rng.Float64()
		p = p.Update(2*x-1, 2*y-1, 2*z-1)
		if p.Dot(p) < 1.0 {
			break
		}
	}

	return p
}

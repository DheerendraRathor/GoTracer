package models

import (
	"math"
	"math/rand"
)

type Sphere struct {
	Center   Point
	Radius   float64
	Material Material
}

func NewSphere(x, y, z, r float64, albedo Material) Sphere {
	return Sphere{
		NewPoint(x, y, z),
		r,
		albedo,
	}
}

func (s Sphere) Hit(r Ray, tmin, tmax float64) (bool, HitRecord) {
	oc := SubtractVectors(r.Origin, s.Center)
	var a, b, c, d float64
	a = VectorDotProduct(r.Direction, r.Direction)
	b = 2.0 * VectorDotProduct(oc, r.Direction)
	c = VectorDotProduct(oc, oc) - s.Radius*s.Radius
	d = b*b - 4*a*c

	record := HitRecord{}
	if d > 0 {
		sqrtD := math.Sqrt(d)
		a2 := 2 * a
		temp := (-b - sqrtD) / a2
		if temp > tmin && temp < tmax {
			record.T = temp
			record.P = r.PointAtParameter(temp)
			record.N = UnitVector(SubtractVectors(record.P, s.Center))
			record.Material = s.Material
			return true, record
		}
		temp = (-b + sqrtD) / a2
		if temp > tmin && temp < tmax {
			record.T = temp
			record.P = r.PointAtParameter(temp)
			record.N = UnitVector(SubtractVectors(record.P, s.Center))
			record.Material = s.Material
			return true, record
		}
	}
	return false, record
}

func (s Sphere) IsHitByRay(r Ray) float64 {
	oc := SubtractVectors(r.Origin, s.Center)
	var a, b, c, d float64
	a = VectorDotProduct(r.Direction, r.Direction)
	b = 2.0 * VectorDotProduct(oc, r.Direction)
	c = VectorDotProduct(oc, oc) - s.Radius*s.Radius
	d = b*b - 4*a*c
	if d < 0 {
		return -1
	} else {
		return (-b - math.Sqrt(d)) / (2 * a)
	}
}

func RandomPointInUnitSphere() Point {
	var p Point
	for {
		x, y, z := 2*rand.Float64()-1, 2*rand.Float64()-1, 2*rand.Float64()-1
		p = NewPoint(x, y, z)
		if VectorDotProduct(p, p) < 1.0 {
			break
		}
	}

	return p
}

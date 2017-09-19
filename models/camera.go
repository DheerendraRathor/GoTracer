package models

import (
	"math"
	"math/rand"
)

type Camera struct {
	LowerLeftCorner, Origin Point
	Horizontal, Vertical    Vector3D
	LensRadius              float64
	U, V, W                 Vector3D
}

func (c Camera) RayAt(u, v float64) Ray {

	rd := MultiplyScalar(RandomPointInUnitDisk(), c.LensRadius)
	offset := AddVectors(MultiplyScalar(c.U, rd.X()), MultiplyScalar(c.V, rd.Y()))
	origin := AddVectors(c.Origin, offset)

	horizontalDir := AddVectors(c.LowerLeftCorner, MultiplyScalar(c.Horizontal, u))
	compositeDir := AddVectors(horizontalDir, MultiplyScalar(c.Vertical, v))
	compositeDir = SubtractVectors(compositeDir, origin)
	return Ray{
		NewPointByVector(origin),
		compositeDir,
	}
}

func NewCamera(lookFrom, lookAt Point, vup Vector3D, vfov, aspect, aperture, focus float64) Camera {
	theta := vfov * math.Pi / 180
	half_height := math.Tan(theta / 2)

	//half_height *= wVector.Length()
	half_width := aspect * half_height

	w := UnitVector(SubtractVectors(lookFrom, lookAt))
	u := UnitVector(VectorCrossProduct(vup, w))
	v := VectorCrossProduct(w, u)

	//llc := NewPoint(-half_width, -half_height, -1.0)
	llc := SubtractVectors(lookFrom, MultiplyScalar(u, half_width*focus))
	llc = SubtractVectors(llc, MultiplyScalar(v, half_height*focus))
	llc = SubtractVectors(llc, MultiplyScalar(w, focus))

	return Camera{
		LowerLeftCorner: NewPointByVector(llc),
		Horizontal:      MultiplyScalar(u, 2*half_width*focus),
		Vertical:        MultiplyScalar(v, 2*half_height*focus),
		Origin:          lookFrom,
		LensRadius:      aperture / 2,
		U:               u,
		V:               v,
		W:               w,
	}
}

func RandomPointInUnitDisk() Point {
	var p Point
	for {
		x, y, z := 2*rand.Float64()-1, 2*rand.Float64()-1, 0.0
		p = NewPoint(x, y, z)
		if VectorDotProduct(p, p) < 1.0 {
			break
		}
	}

	return p
}

package models

import "math"

type Camera struct {
	LowerLeftCorner Point
	Horizontal      Vector3D
	Vertical        Vector3D
	Origin          Point
}

func (c Camera) RayAt(u, v float64) Ray {
	horizontalDir := AddVectors(c.LowerLeftCorner, MultiplyScalar(c.Horizontal, u))
	compositeDir := AddVectors(horizontalDir, MultiplyScalar(c.Vertical, v))
	compositeDir = SubtractVectors(compositeDir, c.Origin)
	return Ray{
		c.Origin,
		compositeDir,
	}
}

func NewCamera(lookFrom, lookAt Point, vup Vector3D, vfov, aspect float64) Camera {
	theta := vfov * math.Pi / 180
	half_height := math.Tan(theta / 2)

	wVector := SubtractVectors(lookFrom, lookAt)

	half_height *= wVector.Length()
	half_width := aspect * half_height

	w := UnitVector(wVector)
	u := UnitVector(VectorCrossProduct(vup, w))
	v := VectorCrossProduct(w, u)

	//llc := NewPoint(-half_width, -half_height, -1.0)
	llc := SubtractVectors(lookFrom, MultiplyScalar(u, half_width))
	llc = SubtractVectors(llc, MultiplyScalar(v, half_height))
	llc = SubtractVectors(llc, wVector)

	return Camera{
		LowerLeftCorner: NewPointByVector(llc),
		Horizontal:      MultiplyScalar(u, 2*half_width),
		Vertical:        MultiplyScalar(v, 2*half_height),
		Origin:          lookFrom,
	}
}

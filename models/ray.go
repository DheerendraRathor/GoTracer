package models

type Ray struct {
	Origin    Point
	Direction Vector3D
}

func (r *Ray) PointAtParameter(t float64) Point {
	dirVec := MultiplyScalar(r.Direction, t)
	return NewPointByVector(AddVectors(r.Origin, dirVec))
}

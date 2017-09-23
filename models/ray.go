package models

type Ray struct {
	Origin    Vector
	Direction Vector
}

func (r *Ray) PointAtParameter(t float64) Vector {
	dirVec := MultiplyScalar(r.Direction, t)
	return AddVectors(r.Origin, dirVec)
}

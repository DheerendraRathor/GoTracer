package models

type Ray struct {
	Origin    *Vector
	Direction *Vector
}

func (r *Ray) PointAtParameter(t float64) *Vector {
	return r.Origin.Copy().AddScaledVector(r.Direction, t)
}

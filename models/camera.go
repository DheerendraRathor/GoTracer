package models

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

func NewCamera() Camera {
	return Camera{
		LowerLeftCorner: NewPoint(-2.0, -1.0, -1.0),
		Horizontal:      NewVector3D(4.0, 0.0, 0.0),
		Vertical:        NewVector3D(0.0, 2.0, 0.0),
		Origin:          NewPoint(0.0, 0.0, 0.0),
	}
}

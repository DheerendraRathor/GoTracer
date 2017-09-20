package models

type Point struct {
	Vector3D
}

func NewPoint(x, y, z float64) Point {
	return Point{
		Vector3D{
			x,
			y,
			z,
		},
	}
}

func NewPointByArray(input [3]float64) Point {
	return Point{
		NewVector3DFromArray(input),
	}
}

func NewPointByVector(v Vector) Point {
	return Point{
		Vector3D{
			v.X(),
			v.Y(),
			v.Z(),
		},
	}
}

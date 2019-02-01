package models

import (
	"math"
)

type Vector struct {
	data [3]float64
}

func NewEmptyVector() *Vector {
	return &Vector{
		[3]float64{0, 0, 0},
	}
}

func NewVectorFromArray(data [3]float64) *Vector {
	return &Vector{data}
}

func NewVector(x, y, z float64) *Vector {
	return &Vector{
		[3]float64{x, y, z},
	}
}

func (v *Vector) Update(x, y, z float64) *Vector {
	v.data[0] = x
	v.data[1] = y
	v.data[2] = z
	return v
}

func (v *Vector) Copy() *Vector {
	return &Vector{
		[3]float64{v.data[0], v.data[1], v.data[2]},
	}
}

func (v *Vector) X() float64 {
	return v.data[0]
}

func (v *Vector) Y() float64 {
	return v.data[1]
}

func (v *Vector) Z() float64 {
	return v.data[2]
}

func Reflect(v, n *Vector) *Vector {
	b2 := n.Copy().Scale(2 * v.Dot(n))
	return v.Copy().SubtractVector(b2)
}

func Refract(v, n *Vector, ni, nt float64) (bool, *Vector) {
	uv := v.Copy()
	uv.MakeUnitVector()
	cosθ := uv.Dot(n)
	snellRatio := ni / nt
	discriminator := 1 - snellRatio*snellRatio*(1-cosθ*cosθ)
	if discriminator > 0 {
		//(uv - n*cosθ)*snellRatio - n*sqrt(disc)
		refracted := uv.SubtractVector(n.Copy().Scale(cosθ)).Scale(snellRatio).
			SubtractVector(n.Copy().Scale(math.Sqrt(discriminator)))
		return true, refracted
	}
	return false, nil
}

func (v *Vector) Negate() *Vector {
	v.data[0] = -v.data[0]
	v.data[1] = -v.data[1]
	v.data[2] = -v.data[2]

	return v
}

func (v *Vector) SquaredLength() float64 {
	return v.Dot(v)
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.SquaredLength())
}

func (v *Vector) MakeUnitVector() *Vector {
	length := v.Length()
	v.data[0] /= length
	v.data[1] /= length
	v.data[2] /= length

	return v
}

func (v *Vector) Scale(t float64) *Vector {
	v.data[0] *= t
	v.data[1] *= t
	v.data[2] *= t

	return v
}

func (v *Vector) AddVector(v1 *Vector) *Vector {
	v.data[0] += v1.data[0]
	v.data[1] += v1.data[1]
	v.data[2] += v1.data[2]

	return v
}

func (v *Vector) AddScaledVector(v1 *Vector, t float64) *Vector {
	v.data[0] += v1.data[0] * t
	v.data[1] += v1.data[1] * t
	v.data[2] += v1.data[2] * t

	return v
}

func (v *Vector) SubtractVector(v1 *Vector) *Vector {
	v.data[0] -= v1.data[0]
	v.data[1] -= v1.data[1]
	v.data[2] -= v1.data[2]

	return v
}

func (v *Vector) MultiplyVector(v1 *Vector) *Vector {
	v.data[0] *= v1.data[0]
	v.data[1] *= v1.data[1]
	v.data[2] *= v1.data[2]

	return v
}

func (v *Vector) Dot(v1 *Vector) float64 {
	return v1.data[0]*v.data[0] + v1.data[1]*v.data[1] + v1.data[2]*v.data[2]
}

func (v *Vector) VectorCrossProduct(v1, v2 *Vector) *Vector {
	v.data[0] = v1.data[1]*v2.data[2] - v1.data[2]*v2.data[1]
	v.data[1] = v1.data[2]*v2.data[0] - v1.data[0]*v2.data[2]
	v.data[2] = v1.data[0]*v2.data[1] - v1.data[1]*v2.data[0]

	return v
}

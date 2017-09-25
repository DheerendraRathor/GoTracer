package models

import (
	"math"
)

type Vector []float64

func (v Vector) X() float64 {
	return v[0]
}

func (v Vector) Y() float64 {
	return v[1]
}

func (v Vector) Z() float64 {
	return v[2]
}

func Reflect(v, n Vector) Vector {
	b2 := MultiplyScalar(n, 2*VectorDotProduct(v, n))
	return SubtractVectors(v, b2)
}

func Refract(v, n Vector, ni, nt float64) (bool, Vector) {
	uv := UnitVector(v)
	cosθ := VectorDotProduct(uv, n)
	snellRatio := ni / nt
	discriminator := 1 - snellRatio*snellRatio*(1-cosθ*cosθ)
	if discriminator > 0 {
		//(uv - n*cosθ)*snellRatio - n*sqrt(disc)
		refracted := uv.Subtract(MultiplyScalar(n, cosθ)).MultiplyScalar(snellRatio).
			Subtract(MultiplyScalar(n, math.Sqrt(discriminator)))
		return true, refracted
	}
	return false, nil
}

func Negate(v Vector) Vector {
	return Vector{
		-v[0],
		-v[1],
		-v[2],
	}
}

func (v Vector) SquaredLength() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.SquaredLength())
}

func (v Vector) Add(v1 Vector) Vector {
	v[0] += v1[0]
	v[1] += v1[1]
	v[2] += v1[2]
	return v
}

func (v Vector) Subtract(v1 Vector) Vector {
	v[0] -= v1[0]
	v[1] -= v1[1]
	v[2] -= v1[2]
	return v
}

func (v Vector) MultiplyScalar(val float64) Vector {
	v[0] *= val
	v[1] *= val
	v[2] *= val
	return v
}

func (v Vector) DivideByScalar(val float64) Vector {
	v[0] /= val
	v[1] /= val
	v[2] /= val
	return v
}

func (v Vector) MakeUnitVector() Vector {
	length := v.Length()
	v[0] /= length
	v[1] /= length
	v[2] /= length
	return v
}

func UnitVector(v Vector) Vector {
	length := v.Length()
	return Vector{
		v[0] / length,
		v[1] / length,
		v[2] / length,
	}
}

func MultiplyScalar(v Vector, t float64) Vector {
	return Vector{
		v[0] * t,
		v[1] * t,
		v[2] * t,
	}
}

func AddVectors(v1, v2 Vector) Vector {
	return Vector{
		v1[0] + v2[0],
		v1[1] + v2[1],
		v1[2] + v2[2],
	}
}

func SubtractVectors(v1, v2 Vector) Vector {
	return Vector{
		v1[0] - v2[0],
		v1[1] - v2[1],
		v1[2] - v2[2],
	}
}

func MultiplyVectors(v1, v2 Vector) Vector {
	return Vector{
		v1[0] * v2[0],
		v1[1] * v2[1],
		v1[2] * v2[2],
	}
}

func VectorDotProduct(v1, v2 Vector) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

func VectorCrossProduct(v1, v2 Vector) Vector {
	return Vector{
		v1[1]*v2[2] - v1[2]*v2[1],
		v1[2]*v2[0] - v1[0]*v2[2],
		v1[0]*v2[1] - v1[1]*v2[0],
	}
}

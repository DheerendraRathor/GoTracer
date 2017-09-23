package models

import "math"

type Vector interface {
	X() float64
	Y() float64
	Z() float64
	Length() float64
}

type Vector3D struct {
	x float64
	y float64
	z float64
}

func Reflect(v, n Vector) *Vector3D {
	b2 := MultiplyScalar(n, 2*VectorDotProduct(v, n))
	return SubtractVectors(v, b2)
}

func Refract(v, n Vector, ni, nt float64) (bool, *Vector3D) {
	uv := UnitVector(v)
	cosθ := VectorDotProduct(uv, n)
	snellRatio := ni / nt
	discriminator := 1 - snellRatio*(1-cosθ*cosθ)
	if discriminator > 0 {
		v1 := MultiplyScalar(SubtractVectors(v, MultiplyScalar(n, cosθ)), snellRatio)
		v2 := MultiplyScalar(n, math.Sqrt(discriminator))
		return true, SubtractVectors(v1, v2)
	}
	return false, nil
}

func NewVector3D(x, y, z float64) *Vector3D {
	return &Vector3D{x, y, z}
}

func NewVector3DFromArray(input [3]float64) *Vector3D {
	return &Vector3D{input[0], input[1], input[2]}
}

func (v *Vector3D) X() float64 {
	return v.x
}

func (v *Vector3D) Y() float64 {
	return v.y
}

func (v *Vector3D) Z() float64 {
	return v.z
}

func Negate(v *Vector3D) *Vector3D {
	return &Vector3D{
		-v.x,
		-v.y,
		-v.z,
	}
}

func (v Vector3D) SquaredLength() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v Vector3D) Length() float64 {
	return math.Sqrt(v.SquaredLength())
}

func (v *Vector3D) Add(v2 Vector) {
	v.x += v2.X()
	v.y += v2.Y()
	v.z += v2.Z()
}

func (v *Vector3D) Subtract(v2 Vector) {
	v.x -= v2.X()
	v.y -= v2.Y()
	v.z -= v2.Z()
}

func (v *Vector3D) MultiplyScalar(val float64) {
	v.x *= val
	v.y *= val
	v.z *= val
}

func (v *Vector3D) DivideByScalar(val float64) {
	v.x /= val
	v.y /= val
	v.z /= val
}

func (v *Vector3D) MakeUnitVector() {
	length := v.Length()
	v.x /= length
	v.y /= length
	v.z /= length
}

func UnitVector(v Vector) *Vector3D {
	length := v.Length()
	return &Vector3D{
		v.X() / length,
		v.Y() / length,
		v.Z() / length,
	}
}

func MultiplyScalar(v Vector, t float64) *Vector3D {
	return &Vector3D{
		v.X() * t,
		v.Y() * t,
		v.Z() * t,
	}
}

func DivideScalar(v Vector, t float64) *Vector3D {
	return &Vector3D{
		v.X() / t,
		v.Y() / t,
		v.Z() / t,
	}
}

func AddVectors(v1, v2 Vector) *Vector3D {
	return &Vector3D{
		v1.X() + v2.X(),
		v1.Y() + v2.Y(),
		v1.Z() + v2.Z(),
	}
}

func SubtractVectors(v1, v2 Vector) *Vector3D {
	return &Vector3D{
		v1.X() - v2.X(),
		v1.Y() - v2.Y(),
		v1.Z() - v2.Z(),
	}
}

func MultiplyVectors(v1, v2 Vector) *Vector3D {
	return &Vector3D{
		v1.X() * v2.X(),
		v1.Y() * v2.Y(),
		v1.Z() * v2.Z(),
	}
}

func VectorDotProduct(v1, v2 Vector) float64 {
	return v1.X()*v2.X() + v1.Y()*v2.Y() + v1.Z()*v2.Z()
}

func VectorCrossProduct(v1, v2 Vector) *Vector3D {
	return &Vector3D{
		v1.Y()*v2.Z() - v1.Z()*v2.Y(),
		-v1.X()*v2.Z() - v1.X()*v2.X(),
		v1.X()*v2.Y() - v1.Y()*v2.X(),
	}
}

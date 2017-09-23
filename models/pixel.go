package models

import "math"

func (v Vector) R() float64 {
	return v[0]
}

func (v Vector) G() float64 {
	return v[1]
}

func (v Vector) B() float64 {
	return v[2]
}

func (v Vector) Gamma2() {
	v[0] = math.Sqrt(v[0])
	v[1] = math.Sqrt(v[1])
	v[2] = math.Sqrt(v[2])
}

func (v Vector) ToUint8() []uint8 {
	result := MultiplyScalar(v, 255.99)
	return []uint8{
		uint8(result[0]),
		uint8(result[1]),
		uint8(result[2]),
	}
}

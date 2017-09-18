package utils

import "math"

func Schlick(cosine, ni, nt float64) float64 {
	r0 := (ni - nt) / (ni + nt)
	r0 *= r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

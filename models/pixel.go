package models

import "math"

func (v *Vector) Gamma2() {
	v.data[0] = math.Sqrt(v.data[0])
	v.data[1] = math.Sqrt(v.data[1])
	v.data[2] = math.Sqrt(v.data[2])
}

type Pixel struct {
	Color [3]uint8
	I, J  int
}

func (v *Vector) ToPixel(i, j int) *Pixel {
	r := v.X() * 255.99
	if r > 255 {
		r = 255
	}

	g := v.Y() * 255.99
	if g > 255 {
		g = 255
	}

	b := v.Z() * 255.99
	if b > 255 {
		b = 255
	}

	return &Pixel{
		[3]uint8{
			uint8(r),
			uint8(g),
			uint8(b),
		},
		i,
		j,
	}
}

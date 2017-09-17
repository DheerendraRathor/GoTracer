package models

type Pixel struct {
	Vector3D
}

func (p *Pixel) R() float64 {
	return p.x
}

func (p *Pixel) G() float64 {
	return p.y
}

func (p *Pixel) B() float64 {
	return p.z
}

func (p *Pixel) UInt8Pixel() Uint8Pixel {
	return Uint8Pixel{
		uint8(255.99 * p.x),
		uint8(255.99 * p.y),
		uint8(255.99 * p.z),
	}
}

func NewPixel(r, g, b float64) Pixel {
	return Pixel{
		Vector3D{
			r,
			g,
			b,
		},
	}
}

func NewPixelFromVector(v Vector) Pixel {
	return Pixel{
		Vector3D{
			v.X(),
			v.Y(),
			v.Z(),
		},
	}
}

type Uint8Pixel struct {
	R uint8
	G uint8
	B uint8
}

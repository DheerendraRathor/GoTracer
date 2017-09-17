package models

type Material interface {
	Scatter(Ray, HitRecord) (bool, Vector3D, Ray)
}

type Lambertian struct {
	albedo Vector3D
}

func NewLambertian(x, y, z float64) Lambertian {
	return Lambertian{
		albedo: NewVector3D(x, y, z),
	}
}

func (l Lambertian) Scatter(ray Ray, hitRecord HitRecord) (bool, Vector3D, Ray) {
	pN := AddVectors(hitRecord.P, hitRecord.N)
	targetPoint := AddVectors(pN, RandomPointInUnitSphere())
	scattered := Ray{
		Origin:    hitRecord.P,
		Direction: SubtractVectors(targetPoint, hitRecord.P),
	}
	return true, l.albedo, scattered
}

type Metal struct {
	albedo Vector3D
	fuzz   float64
}

func NewMetal(x, y, z, fuzz float64) Metal {
	return Metal{
		albedo: NewVector3D(x, y, z),
		fuzz:   fuzz,
	}
}

func (m Metal) Scatter(ray Ray, hitRecord HitRecord) (bool, Vector3D, Ray) {
	reflected := UnitVector(ray.Direction.Reflect(hitRecord.N))
	scattered := Ray{
		hitRecord.P,
		AddVectors(reflected, MultiplyScalar(RandomPointInUnitSphere(), m.fuzz)),
	}
	shouldScatter := VectorDotProduct(reflected, hitRecord.N) > 0
	return shouldScatter, m.albedo, scattered
}

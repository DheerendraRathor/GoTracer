package models

import (
	"github.com/DheerendraRathor/GoTracer/utils"
	"math/rand"
)

type Material interface {
	Scatter(*Ray, *HitRecord) (bool, Vector, *Ray)
}

type BaseMaterial struct {
	Albedo Vector
}

type Lambertian struct {
	*BaseMaterial
}

func NewLambertian(albedo Vector) *Lambertian {
	return &Lambertian{
		&BaseMaterial{
			Albedo: albedo,
		},
	}
}

func (l *Lambertian) Scatter(ray *Ray, hitRecord *HitRecord) (bool, Vector, *Ray) {
	pN := AddVectors(hitRecord.P, hitRecord.N)
	targetPoint := AddVectors(pN, RandomPointInUnitSphere())
	scattered := Ray{
		Origin:    hitRecord.P,
		Direction: SubtractVectors(targetPoint, hitRecord.P),
	}
	return true, l.Albedo, &scattered
}

type Metal struct {
	*BaseMaterial
	fuzz float64
}

func NewMetal(albedo Vector, fuzz float64) *Metal {
	return &Metal{
		BaseMaterial: &BaseMaterial{
			albedo,
		},
		fuzz: fuzz,
	}
}

func (m Metal) Scatter(ray *Ray, hitRecord *HitRecord) (bool, Vector, *Ray) {
	reflected := UnitVector(Reflect(ray.Direction, hitRecord.N))
	scattered := Ray{
		hitRecord.P,
		AddVectors(reflected, MultiplyScalar(RandomPointInUnitSphere(), m.fuzz)),
	}
	shouldScatter := VectorDotProduct(reflected, hitRecord.N) > 0
	return shouldScatter, m.Albedo, &scattered
}

type Dielectric struct {
	*BaseMaterial
	RefIndex float64
}

func (d *Dielectric) Scatter(ray *Ray, hitRecord *HitRecord) (bool, Vector, *Ray) {
	reflected := Reflect(ray.Direction, hitRecord.N)
	var outwardNormal Vector
	var ni, nt float64 = 1, 1
	var cosine, reflectionProb float64
	if VectorDotProduct(ray.Direction, hitRecord.N) > 0 {
		outwardNormal = Negate(hitRecord.N)
		ni = d.RefIndex
		nt = 1
		cosine = d.RefIndex * VectorDotProduct(ray.Direction, hitRecord.N) * ray.Direction.Length()
	} else {
		outwardNormal = hitRecord.N
		ni = 1
		nt = d.RefIndex
		cosine = -VectorDotProduct(ray.Direction, hitRecord.N) * ray.Direction.Length()
	}

	var scattered *Ray
	willRefract, refractedVec := Refract(ray.Direction, outwardNormal, ni, nt)
	if willRefract {
		reflectionProb = utils.Schlick(cosine, 1, d.RefIndex)
		scattered = &Ray{hitRecord.P, refractedVec}
	} else {
		reflectionProb = 1.0
	}

	if rand.Float64() < reflectionProb {
		scattered = &Ray{hitRecord.P, reflected}
	}

	return true, d.Albedo, scattered
}

func NewDielectric(albedo Vector, r float64) *Dielectric {
	return &Dielectric{
		BaseMaterial: &BaseMaterial{
			Albedo: albedo,
		},
		RefIndex: r,
	}
}

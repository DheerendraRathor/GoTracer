package models

import (
	"math/rand"

	"github.com/DheerendraRathor/GoTracer/utils"
)

type Material interface {
	Scatter(*Ray, *HitRecord, *rand.Rand) (bool, *Vector, *Ray)
	IsLight() bool
}

type BaseMaterial struct {
	Albedo  *Vector
	isLight bool
}

func NewBaseMaterial(albedo *Vector, isLight bool) *BaseMaterial {
	return &BaseMaterial{
		Albedo:  albedo,
		isLight: isLight,
	}
}

func (b *BaseMaterial) IsLight() bool {
	return b.isLight
}

type Lambertian struct {
	*BaseMaterial
}

func NewLambertian(albedo *Vector) *Lambertian {
	return &Lambertian{
		BaseMaterial: NewBaseMaterial(albedo, false),
	}
}

func (l *Lambertian) Scatter(ray *Ray, hitRecord *HitRecord, rng *rand.Rand) (bool, *Vector, *Ray) {

	pN := RandomPointInUnitSphere(rng).
		AddVector(hitRecord.N)

	scattered := Ray{
		Origin:    hitRecord.P,
		Direction: pN,
	}

	return true, l.Albedo.Copy(), &scattered
}

type Metal struct {
	*BaseMaterial
	fuzz float64
}

func NewMetal(albedo *Vector, fuzz float64) *Metal {
	return &Metal{
		BaseMaterial: NewBaseMaterial(albedo, false),
		fuzz:         fuzz,
	}
}

func (m *Metal) Scatter(ray *Ray, hitRecord *HitRecord, rng *rand.Rand) (bool, *Vector, *Ray) {
	reflected := Reflect(ray.Direction, hitRecord.N).MakeUnitVector()
	scattered := Ray{
		hitRecord.P,
		reflected.AddScaledVector(RandomPointInUnitSphere(rng), m.fuzz),
	}
	shouldScatter := scattered.Direction.Dot(hitRecord.N) > 0
	return shouldScatter, m.Albedo.Copy(), &scattered
}

type Dielectric struct {
	*BaseMaterial
	RefIndex float64
}

func (d *Dielectric) Scatter(ray *Ray, hitRecord *HitRecord, rng *rand.Rand) (bool, *Vector, *Ray) {
	reflected := Reflect(ray.Direction, hitRecord.N)
	var outwardNormal *Vector
	var ni, nt, cosine, reflectionProb float64
	if ray.Direction.Dot(hitRecord.N) > 0 {
		outwardNormal = hitRecord.N.Copy().Negate()
		ni = d.RefIndex
		nt = 1
		cosine = d.RefIndex * ray.Direction.Dot(hitRecord.N) / ray.Direction.Length()
	} else {
		outwardNormal = hitRecord.N
		ni = 1
		nt = d.RefIndex
		cosine = -ray.Direction.Dot(hitRecord.N) / ray.Direction.Length()
	}

	var scattered *Ray
	willRefract, refractedVec := Refract(ray.Direction, outwardNormal, ni, nt)
	if willRefract {
		reflectionProb = utils.Schlick(cosine, ni, nt)
	} else {
		reflectionProb = 1.0
	}

	if rng.Float64() < reflectionProb {
		scattered = &Ray{hitRecord.P, reflected}
	} else {
		scattered = &Ray{hitRecord.P, refractedVec}
	}

	return true, d.Albedo.Copy(), scattered
}

func NewDielectric(albedo *Vector, r float64) *Dielectric {
	return &Dielectric{
		BaseMaterial: NewBaseMaterial(albedo, false),
		RefIndex:     r,
	}
}

type Light struct {
	*BaseMaterial
}

func NewLight(albedo *Vector) *Light {
	return &Light{
		BaseMaterial: NewBaseMaterial(albedo, true),
	}
}

func (l *Light) Scatter(ray *Ray, hitRecord *HitRecord, rng *rand.Rand) (bool, *Vector, *Ray) {
	return false, l.Albedo.Copy(), nil
}

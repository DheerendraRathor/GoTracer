package models

import (
	"github.com/DheerendraRathor/GoTracer/utils"
	"math/rand"
)

type Material interface {
	Scatter(Ray, HitRecord) (bool, Vector3D, Ray)
}

type Lambertian struct {
	albedo Vector3D
}

func NewLambertian(albedo Vector3D) Lambertian {
	return Lambertian{
		albedo: albedo,
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

func NewMetal(albedo Vector3D, fuzz float64) Metal {
	return Metal{
		albedo: albedo,
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

type Dielectric struct {
	RefIndex float64
	albedo   Vector3D
}

func (d Dielectric) Scatter(ray Ray, hitRecord HitRecord) (bool, Vector3D, Ray) {
	reflected := ray.Direction.Reflect(hitRecord.N)
	var outwardNormal Vector3D
	var ni, nt float64 = 1, 1
	var cosine, reflectionProb float64
	if VectorDotProduct(ray.Direction, hitRecord.N) > 0 {
		outwardNormal = hitRecord.N.Negate()
		ni = d.RefIndex
		nt = 1
		cosine = d.RefIndex * VectorDotProduct(ray.Direction, hitRecord.N) * ray.Direction.Length()
	} else {
		outwardNormal = hitRecord.N
		ni = 1
		nt = d.RefIndex
		cosine = -VectorDotProduct(ray.Direction, hitRecord.N) * ray.Direction.Length()
	}

	var scattered Ray
	willRefract, refractedVec := ray.Direction.Refract(outwardNormal, ni, nt)
	if willRefract {
		reflectionProb = utils.Schlick(cosine, 1, d.RefIndex)
		scattered = Ray{hitRecord.P, refractedVec}
	} else {
		reflectionProb = 1.0
	}

	if rand.Float64() < reflectionProb {
		scattered = Ray{hitRecord.P, reflected}
	}

	return true, d.albedo, scattered
}

func NewDielectric(albedo Vector3D, r float64) Dielectric {
	return Dielectric{
		albedo:   albedo,
		RefIndex: r,
	}
}

package models

import "fmt"

type ObjectType int

const (
	LambertianMaterial = "Lambertian"
	MetalMaterial      = "Metal"
	DielectricMaterial = "Dielectric"
	KaboomMaterial = "Kaboom"
)

type ImageInput struct {
	OutputFile string
	Height     int
	Width      int
	Samples    int
	Patch      [4]int
}

func (i *ImageInput) GetPatch() (int, int, int, int) {
	return i.Patch[0], i.Patch[1], i.Patch[2], i.Patch[3]
}

type CameraInput struct {
	LookFrom    Vector
	LookAt      Vector
	UpVector    Vector
	FieldOfView float64
	AspectRatio float64
	Focus       float64
	Aperture    float64
}

type SurfaceInput struct {
	Type     string
	Albedo   []float64
	Fuzz     float64
	RefIndex float64
}

type SphereInput struct {
	Center  Vector
	Radius  float64
	Surface SurfaceInput
}

type ObjectsInput struct {
	Spheres []SphereInput
}

type Setting struct {
	ShowProgress   bool
	RenderRoutines int
	RenderDepth    int
}

type World struct {
	Settings Setting
	Image    ImageInput
	Camera   CameraInput
	Objects  ObjectsInput
}

func (w World) GetCamera() *Camera {
	return NewCamera(w.Camera.LookFrom, w.Camera.LookAt, w.Camera.UpVector, w.Camera.FieldOfView,
		w.Camera.AspectRatio, w.Camera.Aperture, w.Camera.Focus)
}

func (w World) GetHitableList() *HitableList {
	world := HitableList{}

	for _, sphere := range w.Objects.Spheres {
		world.AddHitable(sphere.getSphere())
	}

	return &world
}

func (s SphereInput) getSphere() *Sphere {
	return &Sphere{
		Radius:   s.Radius,
		Center:   s.Center,
		Material: s.Surface.getMaterial(),
	}
}

func (s *SurfaceInput) getMaterial() Material {

	var material Material

	switch s.Type {
	case LambertianMaterial:
		material = NewLambertian(s.Albedo)
	case MetalMaterial:
		material = NewMetal(s.Albedo, s.Fuzz)
	case DielectricMaterial:
		material = NewDielectric(s.Albedo, s.RefIndex)
	case KaboomMaterial:
		material = NewKaboom(s.Albedo)
	default:
		panic(fmt.Sprintf("Got invalid surface type: %s", s.Type))
	}

	return material
}

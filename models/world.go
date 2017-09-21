package models

import "fmt"

type VectorInput [3]float64
type ObjectType int

const (
	LambertianMaterial = "Lambertian"
	MetalMaterial      = "Metal"
	DielectricMaterial = "Dielectric"
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
	LookFrom    VectorInput
	LookAt      VectorInput
	UpVector    VectorInput
	FieldOfView float64
	AspectRatio float64
	Focus       float64
	Aperture    float64
}

type SurfaceInput struct {
	Type     string
	Albedo   VectorInput
	Fuzz     float64
	RefIndex float64
}

type SphereInput struct {
	Center  VectorInput
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
	lookFrom := NewPointByArray(w.Camera.LookFrom)
	lookAt := NewPointByArray(w.Camera.LookAt)
	vup := NewVector3DFromArray(w.Camera.UpVector)
	return NewCamera(lookFrom, lookAt, vup, w.Camera.FieldOfView, w.Camera.AspectRatio, w.Camera.Aperture, w.Camera.Focus)
}

func (w World) GetHitableList() *HitableList {
	world := HitableList{}

	for _, sphere := range w.Objects.Spheres {
		world.AddHitable(sphere.getSphere())
	}

	return &world
}

func (s SphereInput) getSphere() Sphere {
	return Sphere{
		Radius:   s.Radius,
		Center:   NewPointByArray(s.Center),
		Material: s.Surface.getMaterial(),
	}
}

func (s SurfaceInput) getMaterial() Material {

	var material Material

	switch s.Type {
	case LambertianMaterial:
		material = NewLambertian(NewVector3DFromArray(s.Albedo))
	case MetalMaterial:
		material = NewMetal(NewVector3DFromArray(s.Albedo), s.Fuzz)
	case DielectricMaterial:
		material = NewDielectric(NewVector3DFromArray(s.Albedo), s.RefIndex)
	default:
		panic(fmt.Sprintf("Got invalid surface type: %s", s.Type))
	}

	return material
}

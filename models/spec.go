package models

import "fmt"

type ObjectType int

const (
	LambertianMaterial = "Lambertian"
	MetalMaterial      = "Metal"
	DielectricMaterial = "Dielectric"
	LightMaterial      = "Light"
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
	LookFrom    [3]float64
	LookAt      [3]float64
	UpVector    [3]float64
	FieldOfView float64
	AspectRatio float64
	Focus       float64
	Aperture    float64
}

type SurfaceInput struct {
	Type     string
	Albedo   [3]float64
	Fuzz     float64
	RefIndex float64
}

func (s *SurfaceInput) getMaterial() Material {

	var material Material

	albedo := NewVectorFromArray(s.Albedo)
	switch s.Type {
	case LambertianMaterial:
		material = NewLambertian(albedo)
	case MetalMaterial:
		material = NewMetal(albedo, s.Fuzz)
	case DielectricMaterial:
		material = NewDielectric(albedo, s.RefIndex)
	case LightMaterial:
		material = NewLight(albedo)
	default:
		panic(fmt.Sprintf("Got invalid surface type: %s", s.Type))
	}

	return material
}

type SphereInput struct {
	Center  [3]float64
	Radius  float64
	Surface SurfaceInput
}

func (s SphereInput) getSphere() *Sphere {
	return NewSphere(s.Center[0], s.Center[1], s.Center[2], s.Radius, s.Surface.getMaterial())
}

type ObjectsInput struct {
	Spheres []SphereInput
}

type Setting struct {
	RenderRoutines int
	RenderDepth    int
}

type SceneInput struct {
	Camera       CameraInput
	Objects      ObjectsInput
	AmbientLight [3]float64
}

type Specification struct {
	Settings Setting
	Image    ImageInput
	Scene    SceneInput
}

type Scene struct {
	Camera       *Camera
	HitableList  *HitableList
	AmbientLight *Vector
}

func (w Specification) GetCamera() *Camera {
	camera := w.Scene.Camera
	return NewCamera(
		NewVectorFromArray(camera.LookFrom),
		NewVectorFromArray(camera.LookAt),
		NewVectorFromArray(camera.UpVector),
		camera.FieldOfView,
		camera.AspectRatio, camera.Aperture, camera.Focus)
}

func (w Specification) GetHitableList() *HitableList {
	world := HitableList{}

	for _, sphere := range w.Scene.Objects.Spheres {
		world.AddHitable(sphere.getSphere())
	}

	return &world
}

func (w Specification) GetScene() *Scene {
	return &Scene{
		Camera:       w.GetCamera(),
		HitableList:  w.GetHitableList(),
		AmbientLight: NewVectorFromArray(w.Scene.AmbientLight),
	}
}

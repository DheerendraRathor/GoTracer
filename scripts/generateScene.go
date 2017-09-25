package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"

	"github.com/DheerendraRathor/GoTracer/models"
)

func PositiveRandom() float64 {
	return rand.Float64() * rand.Float64()
}

func AnotherPositiveRandom() float64 {
	return 0.5 * (1 + rand.Float64())
}

func main() {
	flag.Parse()

	spheres := []models.SphereInput{}
	spheres = append(spheres, models.SphereInput{
		Center: []float64{0, -1000, 0},
		Radius: 1000,
		Surface: models.SurfaceInput{
			Type:   models.LambertianMaterial,
			Albedo: []float64{0.5, 0.5, 0.5},
		},
	})

	var matProb float64
	var temp models.Vector
	var cleaner models.Vector = []float64{4, 0.2, 0}

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			matProb = rand.Float64()
			center := []float64{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}

			temp = models.SubtractVectors(center, cleaner)
			if temp.Length() > 0.9 {
				sphere := models.SphereInput{
					Center:  center,
					Radius:  0.2,
					Surface: models.SurfaceInput{},
				}
				if matProb < 0.5 { //diffuse
					sphere.Surface.Type = models.LambertianMaterial
					sphere.Surface.Albedo = []float64{PositiveRandom(), PositiveRandom(), PositiveRandom()}

				} else if matProb < 0.9 { //Metal
					sphere.Surface.Type = models.MetalMaterial
					sphere.Surface.Albedo = []float64{AnotherPositiveRandom(), AnotherPositiveRandom(), AnotherPositiveRandom()}
					sphere.Surface.Fuzz = 0.5 * rand.Float64()
				} else { //glass
					sphere.Surface.Type = models.DielectricMaterial
					sphere.Surface.Albedo = []float64{1.0, 1.0, 1.0}
					sphere.Surface.RefIndex = (rand.Float64() * 0.5) + 1.5
				}
				spheres = append(spheres, sphere)
			}
		}
	}

	spheres = append(
		spheres,
		models.SphereInput{
			Center: []float64{0, 1, 0},
			Radius: 1,
			Surface: models.SurfaceInput{
				Type:     models.DielectricMaterial,
				Albedo:   []float64{1, 1, 1},
				RefIndex: 1.5,
			},
		},
		models.SphereInput{
			Center: []float64{-4, 1, 0},
			Radius: 1,
			Surface: models.SurfaceInput{
				Type:   models.LambertianMaterial,
				Albedo: []float64{0.4, 0.2, 0.1},
			},
		},
		models.SphereInput{
			Center: []float64{4, 1, 0},
			Radius: 1,
			Surface: models.SurfaceInput{
				Type:   models.MetalMaterial,
				Albedo: []float64{0.7, 0.6, 0.5},
				Fuzz:   0,
			},
		},
	)

	jsonData, _ := json.MarshalIndent(spheres, "", "    ")
	fmt.Println(string(jsonData))
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/tracer"
	"github.com/DheerendraRathor/GoTracer/utils"
	"gopkg.in/cheggaaa/pb.v1"
	"image/png"
	"io/ioutil"
	"sync"
)

var renderSpecFile string

func init() {
	flag.StringVar(&renderSpecFile, "spec", "sample_world.json", "Name of JSON file containing rendering spec")
}

func main() {
	flag.Parse()

	file, e := ioutil.ReadFile(renderSpecFile)
	if e != nil {
		panic(fmt.Sprintf("File error: %v\n", e))
	}

	var env models.World
	json.Unmarshal(file, &env)

	progress := make(chan bool, 100)
	defer close(progress)

	var pbWg sync.WaitGroup

	if env.Settings.ShowProgress {
		// Progress Bar
		pbWg.Add(1)
		go func() {
			defer pbWg.Done()
			total := env.Image.Width * env.Image.Height
			bar := pb.StartNew(total)
			bar.ShowFinalTime = true
			bar.ShowTimeLeft = false
			for value := range progress {
				if value {
					break
				}
				bar.Increment()
			}
		}()
	}

	pngImage := goTracer.GoTrace(&env, progress)

	pngFile := utils.CreateNestedFile(env.Image.OutputFile)
	defer pngFile.Close()

	png.Encode(pngFile, pngImage)

	if env.Settings.ShowProgress {
		pbWg.Wait()
	}

}

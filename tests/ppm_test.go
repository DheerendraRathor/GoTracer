package tests

import (
	"testing"

	"os"

	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/utils"
)

func TestRenderPPM(t *testing.T) {
	image := [][]models.Pixel{
		{models.NewPixel(1.0, 0, 0), models.NewPixel(0, 1.0, 0)},
		{models.NewPixel(0, 0, 1.0), models.NewPixel(1.0, 1.0, 1.0)},
	}

	filename := "myTestImage.ppm"
	utils.RenderPPM(image, filename)

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		t.Error("Test failed as file %s doesn't exist", filename)
	}
}

package utils

import (
	"os"

	"fmt"

	"github.com/DheerendraRathor/GoTracer/models"
)

func RenderPPM(pixels [][]models.Pixel, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("P3\n")
	rows, columns := 0, 0
	rows = len(pixels)
	if rows > 0 {
		columns = len(pixels[0])
	}

	file.WriteString(fmt.Sprintf("%d %d\n", columns, rows))
	file.WriteString("255\n")

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			currentPixel := pixels[i][j].UInt8Pixel()
			file.WriteString(fmt.Sprintf("%d %d %d ", currentPixel.R, currentPixel.G, currentPixel.B))
		}
		file.WriteString("\n")
	}
}

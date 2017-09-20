package utils

import (
	"math"
	"os"
	"path"
)

func Schlick(cosine, ni, nt float64) float64 {
	r0 := (ni - nt) / (ni + nt)
	r0 *= r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

func CreateNestedFile(filePath string) *os.File {
	basepath := path.Dir(filePath)

	_, err := os.Stat(basepath)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(basepath, 0777)
	}

	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	return file
}

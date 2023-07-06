package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/fogleman/fauxgl"
	log "github.com/sirupsen/logrus"
)

const (
	width  = 1000
	height = 1000
	fovy   = 40
	near   = 1
	far    = 50
)

var (
	eye    = V(-2.5, 2.5, 2.5) // camera position
	center = V(0, 0, 0)        // view center position
	up     = V(0, 1, 0)

	basePath = "/home/dave/Documents/Tech/3D Prints"
	//file     = "FilamentGrommet"
	outPath = fmt.Sprintf("%s/images", basePath)
	files   = []string{}
)

func main() {
	GetFiles()
	fmt.Println(outPath)
	for _, v := range files {
		err := os.MkdirAll(fmt.Sprintf("%s/%s", outPath, filepath.Dir(v)), os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		DrawImage(v)
	}
}

func DrawImage(fileName string) {
	fmt.Printf("Opening: %s", fmt.Sprintf("%s/%s\n", basePath, fileName))
	mesh, err := LoadSTL(fmt.Sprintf("%s/%s", basePath, fileName))
	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}
	mesh.BiUnitCube()
	mesh.SmoothNormalsThreshold(Radians(30))

	context := NewContext(int(width*mesh.BoundingBox().Size().X), int(height*mesh.BoundingBox().Size().Y))
	context.ClearColor = Black
	context.ClearColorBuffer()

	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	light := V(-2, 0, 1).Normalize()
	color := Color{0, 0.5, 0.65, 1}

	shader := NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader
	context.DrawMesh(mesh)

	outFile := fmt.Sprintf("%s/%s.png", outPath, fileName)
	err = SavePNG(outFile, context.Image())

	if err != nil {
		fmt.Errorf("Write Error: %v\n", err.Error())
		os.Exit(12)

	} else {
		fmt.Printf("Drew: %s\n", outFile)
	}

}

func GetFiles() {

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if strings.HasSuffix(strings.ToLower(info.Name()), ".stl") {
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				log.Error(err)
			}
			files = append(files, relPath)
		}

		return nil
	})

}

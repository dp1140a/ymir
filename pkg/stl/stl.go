package stl

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	. "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
)

const (
	width  = 640
	height = 480
	fovy   = 40
	near   = 1
	far    = 50
)

var (
	eye    = V(.75, 2, 3) // camera position
	center = V(0, 0, 0)   // view center position
	up     = V(0, 1, 0)

	//file     = "FilamentGrommet"
	outPath = fmt.Sprintf("%s/images", basePath)
	files   = []string{}
)

func SaveImage(path string, image image.Image) {
	err := SavePNG(path, image)
	if err != nil {
		_ = fmt.Errorf("Write Error: %v\n", err.Error())
		os.Exit(12)

	} else {
		fmt.Printf("Drew: %s\n", path)
	}
}

func Image(fileName string) image.Image {
	log.Infof("creating image from %s", fileName)
	mesh, err := LoadSTL(fileName)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	mesh.BiUnitCube()
	mesh.SmoothNormalsThreshold(Radians(2))

	context := NewContext(int(width*mesh.BoundingBox().Size().X), int(height*mesh.BoundingBox().Size().Y))
	context.ClearColor = Color{R: 0.5, G: 0.5, B: 0.5}
	context.ClearColorBuffer()

	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	light := V(2, 2, 2).Normalize()
	color := Color{R: .65, G: 0.24, B: 0.14, A: 1}

	shader := NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader
	shader.AmbientColor = Color{R: 0.6, G: 0.6, B: 0.6, A: 1}
	context.DrawMesh(mesh)

	return context.Image()
}

func Thumbnail(srcImage image.Image, maxHeight uint, maxWidth uint) image.Image {
	return resize.Thumbnail(maxWidth, maxHeight, srcImage, resize.Bilinear)
}

func ThumbnailBase64(srcImage image.Image, maxHeight uint, maxWidth uint) string {
	img := Thumbnail(srcImage, maxHeight, maxWidth)
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		log.Error(err)
	}

	return fmt.Sprintf("data:image/png;base64, %s", base64.StdEncoding.EncodeToString(buf.Bytes()))
}

func Png(image image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, image); err != nil {
		log.Error(err)
		return nil, errors.New("unable to encode png")
	}

	return buf.Bytes(), nil
}

func JPG(image image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, image, nil); err != nil {
		log.Error(err)
		return nil, errors.New("unable to encode jpg")
	}

	return buf.Bytes(), nil
}

func GetFiles(basePath string) {
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

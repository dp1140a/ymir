package stl

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var basePath = "/home/dave/Documents/Tech/3D Prints"

func main() {
	GetFiles(basePath)
	fmt.Println(outPath)
	for _, v := range files {
		err := os.MkdirAll(fmt.Sprintf("%s/%s", outPath, filepath.Dir(v)), os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		Image(v)
	}
}

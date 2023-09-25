package utils

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

func MakeDirIfNotExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Errorf("error creating dir %v: %v", path, err)
		}
	}
}

package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
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

func GenId() string {
	buf := make([]byte, 8)
	// then we can call rand.Read.
	_, _ = rand.Read(buf)

	return fmt.Sprintf("%x", buf)
}

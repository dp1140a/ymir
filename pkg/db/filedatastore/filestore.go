package filedatastore

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"ymir/pkg/db"
)

const _INDEX_FILE = "ymir.json"

type FileDataStore struct {
	db.Datastore
	Config *FileDataStoreConfig
	index  *os.File
}

func NewFileDataStore() (datastore *FileDataStore) {
	f := &FileDataStore{
		Config: NewFileDataStoreConfig(),
	}
	err := f.open()
	if err != nil {
		log.Error("cant create file store")
		return nil
	}
	return f
}

/*
*
Satisfies Store Interface
*/
func (fs *FileDataStore) Close() error {
	return fs.index.Close()
}

func (fs *FileDataStore) GetType() string {
	return "fle"
}

func (fs FileDataStore) open() error {
	if f, err := os.Open(filepath.Join(fs.Config.BasePath, _INDEX_FILE)); err == nil {
		fmt.Printf("File exists\n")
		fs.index = f
		return nil
	} else {
		log.Errorf("cant open index file: %v", err)
		return err
	}
}

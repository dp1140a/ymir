package filedatastore

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_FILE_DATASTORE = "datastore.file"
)

type FileDataStoreConfig struct {
	BasePath string `toml:"basePath"`
}

func NewFileDataStoreConfig() *FileDataStoreConfig {
	c := &FileDataStoreConfig{
		BasePath: "ymir",
	}

	h := viper.Sub(_FILE_DATASTORE)
	if h != nil {
		err := h.Unmarshal(c)
		if err != nil {
			log.Error(_FILE_DATASTORE, " config error: ", err.Error())
		}
	}
	return c
}

func (c *FileDataStoreConfig) StringJSON() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *FileDataStoreConfig) StringToml() (config string) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(c)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

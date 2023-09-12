package model

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_MODELS = "models"
)

/*
[models]
uploadsTempDir="uploads/tmp"
uploadsFilesDir="uploads/modelFiles"
*/
type ModelsConfig struct {
	UploadsTempDir string `toml:"uploadsTempDir"`
	ModelsDir      string `toml:"ModelsDir"`
}

func NewModelsConfig() *ModelsConfig {
	c := &ModelsConfig{
		ModelsDir:      "uploads/modelFiles",
		UploadsTempDir: "uploads/tmp",
	}

	h := viper.Sub(_MODELS)
	if h != nil {
		err := h.Unmarshal(c)
		if err != nil {
			log.Error(_MODELS, " config error: ", err.Error())
		}
	}
	return c
}

func (m *ModelsConfig) StringJSON() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func (m *ModelsConfig) StringToml() (config string) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(m)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

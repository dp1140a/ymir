package boltdatastore

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_BOLTDB_DATASTORE = "datastore"
)

type BoltDBDataStoreConfig struct {
	DBFile string `toml:"dbFile"`
}

func NewBoltDBDataStoreConfig() *BoltDBDataStoreConfig {
	c := &BoltDBDataStoreConfig{
		DBFile: "ymir.db",
	}

	h := viper.Sub(_BOLTDB_DATASTORE)
	if h != nil {
		err := h.Unmarshal(c)
		if err != nil {
			log.Error(_BOLTDB_DATASTORE, " config error: ", err.Error())
		}
	}
	return c
}

func (c *BoltDBDataStoreConfig) StringJSON() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *BoltDBDataStoreConfig) StringToml() (config string) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(c)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

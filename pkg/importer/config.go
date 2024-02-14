package importer

import (
	"github.com/spf13/viper"
	"ymir/pkg/db/boltdatastore"
)

type ImporterConfig struct {
	storeConfig *boltdatastore.BoltDBDataStoreConfig
	modelsBase  string
}

func NewImporterConfig() *ImporterConfig {
	return &ImporterConfig{
		storeConfig: boltdatastore.NewBoltDBDataStoreConfig(),
		modelsBase:  viper.GetString("models.modelsDir"),
	}
}

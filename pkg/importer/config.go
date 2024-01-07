package importer

import (
	"ymir/pkg/config"
	"ymir/pkg/db/boltdatastore"
)

type ImporterConfig struct {
	storeConfig *boltdatastore.BoltDBDataStoreConfig
}

func NewImporterConfig() *ImporterConfig {
	config.InitConfig()
	return &ImporterConfig{storeConfig: boltdatastore.NewBoltDBDataStoreConfig()}
}

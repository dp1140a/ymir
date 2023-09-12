package importer

import (
	"ymir/pkg/db"
)

type ImporterConfig struct {
	DB *db.DBConfig
}

func NewImporterConfig() *ImporterConfig {
	return &ImporterConfig{DB: db.NewDBConfig()}
}

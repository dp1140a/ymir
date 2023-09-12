package db

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_DB = "db"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string `toml:"dbname"`
}

func NewDBConfig() *DBConfig {
	c := &DBConfig{
		Host:     "localhost",
		Port:     "5984",
		Username: "admin",
		Password: "password",
		DBName:   "ymir",
	}

	h := viper.Sub(_DB)
	if h != nil {
		err := h.Unmarshal(c)
		if err != nil {
			log.Error(_DB, " config error: ", err.Error())
		}
	}

	return c
}

func (db *DBConfig) StringJSON() string {
	b, _ := json.Marshal(db)
	return string(b)
}

func (db *DBConfig) StringToml() (config string) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(db)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

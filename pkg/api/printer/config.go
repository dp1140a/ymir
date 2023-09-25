package printer

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_PRINTERS = "printers"
)

type PrintersConfig struct {
	PrintersDir string `toml:"printersDir"`
}

func NewPrintersConfig() *PrintersConfig {
	c := &PrintersConfig{
		PrintersDir: "uploads/printers",
	}

	h := viper.Sub(_PRINTERS)
	if h != nil {
		err := h.Unmarshal(c)
		if err != nil {
			log.Error(_PRINTERS, " config error: ", err.Error())
		}
	}
	return c
}

func (p *PrintersConfig) StringJSON() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (p *PrintersConfig) StringToml() (config string) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

package httplogger

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_LOGGING_HTTP = "http.logging"
)

type HttpLoggerConfig struct {
	Enabled bool   `json:"enabled" toml:"enabled"`
	StdOut  bool   `json:"stdOut" toml:"stdOut"`
	FileOut bool   `json:"fileOut" toml:"fileOut"`
	LogFile string `json:"logFile" toml:"logFile"`
}

//New
/**
Returns a default populated *HttpLoggerConfig
*/
func NewHttpLoggerConfig() *HttpLoggerConfig {
	httpLoggerConf := &HttpLoggerConfig{}
	h := viper.Sub(_LOGGING_HTTP)
	//fmt.Printf("H: %v", h)
	if h != nil {
		err := h.Unmarshal(httpLoggerConf)
		if err != nil {
			log.Error(_LOGGING_HTTP, " config error: ", err.Error())
			return httpLoggerConf
		}
	}
	return httpLoggerConf
}

func (lc *HttpLoggerConfig) String() string {
	b, _ := json.Marshal(lc)
	return string(b)
}

package logger

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	_LOGGING = "logging"
)

type LoggerConfig struct {
	LogFile  string `json:"logFile" toml:"logFile"`
	LogLevel string `json:"logLevel" toml:"logLevel"`
	StdOut   bool   `json:"stdOut" toml:"stdOut"`
	FileOut  bool   `json:"fileOut" toml:"fileOut"`
}

//New
/**
Returns a default populated *LoggerConfig
*/
func NewLoggerConfig() *LoggerConfig {
	loggerConfig := &LoggerConfig{
		LogFile:  "log/ymir.log",
		LogLevel: "INFO",
		StdOut:   true,
		FileOut:  false,
	}
	h := viper.Sub(_LOGGING)
	if h != nil {
		err := h.Unmarshal(loggerConfig)
		if err != nil {
			log.Error(_LOGGING, " config error: ", err.Error())
			return nil
		}
	}

	return loggerConfig
}

func (lc *LoggerConfig) String() string {
	b, _ := json.Marshal(lc)
	return string(b)
}

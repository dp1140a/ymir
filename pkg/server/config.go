package server

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"ymir/pkg/logger/httplogger"
)

const (
	_HTTP = "http"
)

type ServerConfig struct {
	Hostname             string                       `json:"hostname" toml:"hostname"`
	Port                 string                       `json:"port" toml:"port"`
	UseHttps             bool                         `json:"usehttps" toml:"usehttps"`
	TLSMinVersion        string                       `json:"TLSMinVersion" toml:"TLSMinVersion"`
	HttpTLSStrictCiphers bool                         `json:"HttpTLSStrictCiphers" toml:"HttpTLSStrictCiphers"`
	TLSCert              string                       `json:"TLSCert" toml:"TLSCert"`
	TLSKey               string                       `json:"TLSKey" toml:"TLSKey"`
	EnableCORS           bool                         `json:"enableCORS" toml:"enableCORS"`
	JWTSecret            string                       `json:"JWTSecret" toml:"JWTSecret"`
	HttpLogConfig        *httplogger.HttpLoggerConfig `json:"logging" toml:"logging"`
}

/*
*
Returns a populated *ServerConfig
*/
func NewServerConfig() *ServerConfig {
	serverConfig := &ServerConfig{}
	h := viper.Sub(_HTTP)
	if h != nil {
		err := h.Unmarshal(serverConfig)
		if err != nil {
			log.Error(_HTTP, " config error: ", err.Error())
			return nil
		}
	}

	serverConfig.HttpLogConfig = httplogger.NewHttpLoggerConfig()

	return serverConfig
}

func (lc *ServerConfig) String() string {
	b, _ := json.Marshal(lc)
	return string(b)
}

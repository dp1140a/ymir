package httplogger

import (
	"bytes"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type RequestLoggerTestSuite struct {
	suite.Suite
}

func (suite *RequestLoggerTestSuite) SetupSuite() {
	fmt.Println("SetupSuite()")
	viper.SetConfigType("toml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var tomlExample = []byte(`
[http.logging]
enabled = true
stdOut = true
fileOut = true
logFile = "log/ymir_http.log"
`)

	err := viper.ReadConfig(bytes.NewBuffer(tomlExample))
	if err != nil {
		suite.T().Errorf("Error: %v", err)
	}
}

func (suite *RequestLoggerTestSuite) SetupTest() {

}

func (suite *RequestLoggerTestSuite) TearDownTest() {
}

func (suite *RequestLoggerTestSuite) TearDownSuite() {
}

func (suite *RequestLoggerTestSuite) TestNewRequestLogger() {
	config := NewHttpLoggerConfig()
	fmt.Println(config.String())
	NewStructuredLogger(log.New(), config)
}

func TestRequestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(RequestLoggerTestSuite))
}

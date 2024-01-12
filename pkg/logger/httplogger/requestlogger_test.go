package httplogger

import (
	"bytes"
	"errors"
	"fmt"
	"os"
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
	fmt.Println("TearDownSuite()")
	fmt.Println(viper.GetString("http.logging.logfile"))
	if _, err := os.Stat(viper.GetString("http.logging.logfile")); errors.Is(err, os.ErrNotExist) {
		fmt.Println("File Does NOT exist")
	}
	fmt.Println("removing file")
	err := os.RemoveAll("log")

	if err != nil {
		fmt.Printf("error removing: %v\n", err.Error())
	}
}

func (suite *RequestLoggerTestSuite) TestNewRequestLogger() {
	config := NewHttpLoggerConfig()
	fmt.Println(config.String())
	NewStructuredLogger(log.New(), config)
}

func TestRequestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(RequestLoggerTestSuite))
}

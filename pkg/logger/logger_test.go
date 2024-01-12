package logger

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func (suite *LoggerTestSuite) SetupSuite() {
	fmt.Println("SetupSuite()")
	viper.SetConfigType("toml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var tomlExample = []byte(`
[logging]
logFile = "log/ymir.log"
logLevel="DEBUG"
stdOut = true
fileOut = false
`)

	err := viper.ReadConfig(bytes.NewBuffer(tomlExample))
	if err != nil {
		suite.T().Errorf("Error: %v", err)
	}
}

func (suite *LoggerTestSuite) SetupTest() {

}

func (suite *LoggerTestSuite) TearDownTest() {
}

func (suite *LoggerTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite()")
	if _, err := os.Stat(viper.GetString("logging.logfile")); errors.Is(err, os.ErrNotExist) {
		return
	} else {
		err := os.Remove("log")
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}
}

func (suite *LoggerTestSuite) TestInitLogger() {
	err := InitLogger()
	assert.NoError(suite.T(), err, "should be nil")
}

func (suite *LoggerTestSuite) TestNewLogger() {
	logger := NewLogger()
	assert.IsType(suite.T(), &Logger{}, logger, "should be *Logger")
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

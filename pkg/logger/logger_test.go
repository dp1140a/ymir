package logger

import (
	"bytes"
	"fmt"
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
[datastore]
dbFile = "test.db"
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

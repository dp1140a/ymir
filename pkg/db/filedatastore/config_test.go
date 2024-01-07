package filedatastore

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileDataStoreConfigTestSuite struct {
	suite.Suite
	config *FileDataStoreConfig
}

func (suite *FileDataStoreConfigTestSuite) SetupSuite() {
	fmt.Println("SetupSuite()")
	suite.config = &FileDataStoreConfig{}
}

func (suite *FileDataStoreConfigTestSuite) SetupTest() {

}

func (suite *FileDataStoreConfigTestSuite) TearDownTest() {

}

func (suite *FileDataStoreConfigTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite()")
}

func (suite *FileDataStoreConfigTestSuite) TestNewFileDataStoreConfig_Default() {
	config := NewFileDataStoreConfig()

	assert.NotNil(suite.T(), config, "Should not be nil")
	assert.IsType(suite.T(), &FileDataStoreConfig{}, config, "Should be of type *FileDataStoreConfig")
	assert.Equal(suite.T(), "ymir", config.BasePath, "should be ymir by default")
}

func (suite *FileDataStoreConfigTestSuite) TestNewFileDataStoreConfig_With_Config() {
	viper.SetConfigType("toml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var tomlExample = []byte(`[datastore]
[datastore.file]
basePath="ymir/models"`)

	err := viper.ReadConfig(bytes.NewBuffer(tomlExample))
	if err != nil {
		suite.T().Errorf("Error: %v", err)
	}
	config := NewFileDataStoreConfig()

	assert.NotNil(suite.T(), config, "Should not be nil")
	assert.IsType(suite.T(), &FileDataStoreConfig{}, config, "Should be of type *FileDataStoreConfig")
	assert.Equal(suite.T(), "ymir/models", config.BasePath, "should be ymir/models")
}

func TestFileDataStoreConfigTestSuite(t *testing.T) {
	suite.Run(t, new(FileDataStoreConfigTestSuite))
}

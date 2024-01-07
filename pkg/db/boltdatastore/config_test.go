package boltdatastore

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoltDataStoreConfigTestSuite struct {
	suite.Suite
	config *BoltDBDataStoreConfig
}

func (suite *BoltDataStoreConfigTestSuite) SetupSuite() {
	fmt.Println("SetupSuite()")
	suite.config = &BoltDBDataStoreConfig{}
}

func (suite *BoltDataStoreConfigTestSuite) SetupTest() {

}

func (suite *BoltDataStoreConfigTestSuite) TearDownTest() {

}

func (suite *BoltDataStoreConfigTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite()")
}

func (suite *BoltDataStoreConfigTestSuite) TestNewBoltDBDataStoreConfig_Default() {
	config := NewBoltDBDataStoreConfig()

	assert.NotNil(suite.T(), config, "Should not be nil")
	assert.IsType(suite.T(), &BoltDBDataStoreConfig{}, config, "Should be of type *BoltDBDataStoreConfig")
	assert.Equal(suite.T(), "ymir.db", config.DBFile, "should be ymir.db by default")
}

func (suite *BoltDataStoreConfigTestSuite) TestNewBoltDBDataStoreConfig_With_Config() {
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
	config := NewBoltDBDataStoreConfig()

	assert.NotNil(suite.T(), config, "Should not be nil")
	assert.IsType(suite.T(), &BoltDBDataStoreConfig{}, config, "Should be of type *BoltDBDataStoreConfig")
	assert.Equal(suite.T(), "test.db", config.DBFile, "should be test.db")
}

func TestBoltDataStoreConfigTestSuite(t *testing.T) {
	suite.Run(t, new(BoltDataStoreConfigTestSuite))
}

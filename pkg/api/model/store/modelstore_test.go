package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"ymir/pkg/api/model/types"
)

const (
	TEST_DB = "test.db"
)

type ModelStoreTestSuite struct {
	suite.Suite
	store      ModelStore
	testModels []types.Model
}

func (suite *ModelStoreTestSuite) SetupSuite() {
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
	suite.store = NewModelDataStore().(ModelStore)
	suite.testModels = getTestModels()
}

func (suite *ModelStoreTestSuite) SetupTest() {
	err := suite.store.Create(suite.testModels[0])
	assert.NoError(suite.T(), err, "should be no error on test setup")
}

func (suite *ModelStoreTestSuite) TearDownTest() {
	err := suite.store.Delete(suite.testModels[1].Id)
	assert.NoError(suite.T(), err, "should be no error on test teardown")
}

func (suite *ModelStoreTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite()")
	if _, err := os.Stat(TEST_DB); errors.Is(err, os.ErrNotExist) {
		fmt.Println("DB Does NOT exist")
	}
	err := os.Remove(TEST_DB)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func (suite *ModelStoreTestSuite) TestNewModelStore() {
	assert.NotNil(suite.T(), suite.store, "Should not be nil")
	assert.Equal(suite.T(), TEST_DB, suite.store.ds.GetDB().Path(), "should be test.db")
	assert.Equal(suite.T(), "boltDB", suite.store.ds.GetType(), "Should be boltDB")
}

func (suite *ModelStoreTestSuite) TestModelCreate() {
	err := suite.store.Create(suite.testModels[1])
	assert.NoError(suite.T(), err)
	numModels, _ := suite.store.ds.GetNumKeys(MODELS_BUCKET)
	assert.Equal(suite.T(), numModels, 2, "should be 2")
}

func (suite *ModelStoreTestSuite) TestModelUpdate() {
	newMod := suite.testModels[0]
	newMod.Description = "TestModelUpdate"
	err := suite.store.Update(newMod)
	assert.NoError(suite.T(), err, "should be no error")
	mod, _ := suite.store.Inspect(suite.testModels[0].Id)
	assert.Equal(suite.T(), newMod, mod, "should be equal")
}

func (suite *ModelStoreTestSuite) TestModelDelete() {
	err := suite.store.Delete(suite.testModels[0].Id)
	assert.NoError(suite.T(), err, "should be no error")
	numModels, _ := suite.store.ds.GetNumKeys(MODELS_BUCKET)
	assert.Equal(suite.T(), numModels, 0, "should be 0")
}

func (suite *ModelStoreTestSuite) TestModelList() {
	models, err := suite.store.List()
	assert.NoError(suite.T(), err, "should be no error")
	assert.Len(suite.T(), models, 1, "should be 1")
	assert.IsType(suite.T(), map[string]types.Model{}, models, "should be type []model.Model{}")
	assert.IsType(suite.T(), types.Model{}, models[suite.testModels[0].Id], "should be of type model")
	assert.Equal(suite.T(), suite.testModels[0], models[suite.testModels[0].Id], "should be equal")
}

func (suite *ModelStoreTestSuite) TestModelInspect() {
	mod, err := suite.store.Inspect(suite.testModels[0].Id)
	assert.NoError(suite.T(), err, "should be no error")
	assert.IsType(suite.T(), types.Model{}, mod, "should be of type model")
	assert.Equal(suite.T(), suite.testModels[0], mod, "should be equal")
}

func TestModelStoreTestSuite(t *testing.T) {
	suite.Run(t, new(ModelStoreTestSuite))
}

/*
Utility Functions
*/
func getTestModels() []types.Model {
	models := []types.Model{}
	cwd, _ := os.Getwd()

	d := filepath.Join(cwd, "../testdata")

	dirs, err := os.ReadDir(d)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return nil
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			modelFilePath := filepath.Join(d, dir.Name(), "model.json")
			if fileExists(modelFilePath) {
				mod, err := readModelFile(modelFilePath)
				if err != nil {
					fmt.Printf("Error reading model.json in %s: %v\n", dir.Name(), err)
					continue
				}

				models = append(models, mod)
			}
		}
	}

	return models
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func readModelFile(filePath string) (types.Model, error) {
	fmt.Println(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return types.Model{}, err
	}

	var mod types.Model
	err = json.Unmarshal(data, &mod)
	if err != nil {
		return types.Model{}, err
	}
	return mod, nil
}

package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	suite.Suite
	models []Model
}

func (suite *ModelTestSuite) SetupSuite() {
	suite.models = getTestModels()
	fmt.Printf("%v models\n", len(suite.models))
}

func (suite *ModelTestSuite) SetupTest() {

}

func (suite *ModelTestSuite) TeardownTest() {

}

func (suite *ModelTestSuite) TeardownSuite() {

}

func (suite *ModelTestSuite) TestModel_mock(t *testing.T) {
	assert.True(suite.T(), true, "all good")
}

func getTestModels() []Model {
	models := []Model{}
	cwd, _ := os.Getwd()

	err := filepath.Walk(filepath.Join(cwd, "testdata"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != "testData" {
			modelFilePath := filepath.Join(path, "model.json")
			model, err := readModelFile(modelFilePath)
			if err != nil {
				fmt.Printf("Error reading model.json in %s: %v\n", path, err)
				return nil
			}

			models = append(models, model)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
		return nil
	}

	return models

}

func readModelFile(filePath string) (Model, error) {
	fmt.Println(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Model{}, err
	}

	var model Model
	err = json.Unmarshal(data, &model)
	if err != nil {
		return Model{}, err
	}

	return model, nil
}

func TestModelTestSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

package model

import (
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
}

func (suite *ModelTestSuite) SetupTest() {

}

func (suite *ModelTestSuite) TeardownTest() {

}

func (suite *ModelTestSuite) TeardownSuite() {

}

func (suite *ModelTestSuite) Test_getTestModels() {
	assert.Equal(suite.T(), 2, len(suite.models), "Should be 2 model")
}

func TestModelTestSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

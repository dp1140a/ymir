package model

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"ymir/pkg/api/model/store"
	"ymir/pkg/api/model/types"
	"ymir/pkg/logger"
)

const (
	TEST_DIR = "TEST_DIR"
)

type ModelServiceTestSuite struct {
	suite.Suite
	service    ModelService
	testModels []types.Model
}

func (suite *ModelServiceTestSuite) SetupSuite() {
	fmt.Println("SetupSuite()")
	viper.SetConfigType("toml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var tomlExample = []byte(`
[logging]
logLevel="DEBUG"

[datastore]
dbFile = "test.db"

[models]
uploadsTempDir="TEST_DIR/tmp"
modelsDir="TEST_DIR/models"
`)

	err := viper.ReadConfig(bytes.NewBuffer(tomlExample))
	if err != nil {
		suite.T().Errorf("Error: %v", err)
	}

	suite.service = NewModelService().(ModelService)
	err = logger.InitLogger()
	if err != nil {
		fmt.Println(err.Error())
		suite.T().Fatal(err.Error())
	}
	suite.testModels = getTestModels()
}

func (suite *ModelServiceTestSuite) AfterTest(suiteName, testName string) {}

func (suite *ModelServiceTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestCreateModel" {
		err := suite.service.modelStore.(store.ModelStore).Truncate()
		assert.NoError(suite.T(), err)
	}
}

func (suite *ModelServiceTestSuite) SetupTest() {
	fmt.Println("SetupTest")
	err := suite.service.modelStore.(store.ModelStore).Truncate()
	assert.NoError(suite.T(), err)
	err = suite.service.modelStore.Create(suite.testModels[0])
	assert.NoError(suite.T(), err, "should be no error on test setup")
}

func (suite *ModelServiceTestSuite) TearDownTest() {
}

func (suite *ModelServiceTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite()")
	suite.T().Cleanup(func() {
	})
	err := os.RemoveAll(TEST_DIR)
	err = os.Remove("test.db")

	if err != nil {
		log.Fatal(err)
	}
}

func (suite *ModelServiceTestSuite) TestNewModelService() {
	assert.NotNil(suite.T(), suite.service, "Should not be nil")
}

func (suite *ModelServiceTestSuite) TestModelServiceGetName() {
	assert.Equal(suite.T(), "Model", suite.service.GetName())
}

func (suite *ModelServiceTestSuite) TestCreateModel() {
	id, err := suite.service.CreateModel(types.Model{})
	suite.testModels[0].Id = id
	mod, err := suite.service.GetModel(id)
	assert.Equal(suite.T(), suite.testModels[0].Id, mod.Id, "should be same")
	assert.NoError(suite.T(), err)
}

func (suite *ModelServiceTestSuite) TestUpdateModel() {
	suite.testModels[0].Description = "Updated"
	err := suite.service.UpdateModel(suite.testModels[0])
	mod, err := suite.service.GetModel(suite.testModels[0].Id)
	assert.Equal(suite.T(), suite.testModels[0].Description, mod.Description, "should be same")
	assert.NoError(suite.T(), err)
}

func (suite *ModelServiceTestSuite) TestListModels() {
	models, err := suite.service.ListModels()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), models, 1, "should be 1")
	fmt.Println(models)
	assert.IsType(suite.T(), map[string]types.Model{}, models, "should be type []model.Model{}")
	assert.IsType(suite.T(), types.Model{}, models[suite.testModels[0].Id], "should be of type model")
	assert.ObjectsAreEqual(suite.testModels[0], models[suite.testModels[0].Id])
}

func (suite *ModelServiceTestSuite) TestGetModel() {
	mod, err := suite.service.GetModel(suite.testModels[0].Id)
	assert.NoError(suite.T(), err, "should be no error")
	assert.IsType(suite.T(), types.Model{}, mod, "should be of type model")
	assert.Equal(suite.T(), suite.testModels[0].Id, mod.Id, "should be equal")
}

func (suite *ModelServiceTestSuite) TestDeleteModel() {
	err := suite.service.DeleteModel(suite.testModels[0].Id)
	assert.NoError(suite.T(), err)
	_, err = suite.service.GetModel(suite.testModels[0].Id)
	//fmt.Println(mod.Json())
	assert.Error(suite.T(), err)
}

func (suite *ModelServiceTestSuite) TestExportModel() {
	zipFile, err := os.Create(filepath.Join(TEST_DIR, "model.zip"))
	assert.NoError(suite.T(), err)
	err = suite.service.ExportModel("../../testdata/model1", zipFile)
	assert.FileExists(suite.T(), filepath.Join(TEST_DIR, "model.zip"))
	info, _ := zipFile.Stat()
	assert.NotEqual(suite.T(), 0, info.Size())
	assert.NotNil(suite.T(), zipFile)
	assert.NoError(suite.T(), err)
}

func (suite *ModelServiceTestSuite) TestUploadFilesExistingModel() {
	f, err := os.Create("TEST_DIR/test.dat")
	assert.NoError(suite.T(), err)
	key, err := suite.service.UploadFilesExistingModel(f, "test_upload.dat", "TEST_DIR")
	assert.NoError(suite.T(), err)
	if err != nil {
		log.Fatal()
	}
	assert.NotNil(suite.T(), key)
	assert.FileExists(suite.T(), "TEST_DIR/test_upload.dat")
}

func (suite *ModelServiceTestSuite) TestUploadFilesNewModel() {
	f, err := os.Create(filepath.Join(TEST_DIR, "test_new.dat"))
	assert.NoError(suite.T(), err)
	key, err := suite.service.UploadFilesNewModel(f, "test_new.dat")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), key)
	assert.FileExists(suite.T(), key)
}

func (suite *ModelServiceTestSuite) TestFetchModelImage() {
	img, err := suite.service.FetchModelImage(filepath.Join("testdata/model1", suite.testModels[0].Images[0].Path))
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), img)
	assert.Equal(suite.T(), 44000, len(img))
}
func (suite *ModelServiceTestSuite) TestFetchSTL() {
	stl, err := suite.service.FetchSTL(filepath.Join("testdata/model1", suite.testModels[0].ModelFiles[0].Path))
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stl)
	assert.Equal(suite.T(), 44284, len(stl))
}
func (suite *ModelServiceTestSuite) TestFetchSTLThumbnail() {
	thumb, err := suite.service.FetchSTLThumbnail(filepath.Join("testdata/model1", suite.testModels[0].ModelFiles[0].Path))
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), thumb, "data:image/png;base64,")
}

func (suite *ModelServiceTestSuite) TestAddNote() {
	note := types.Note{
		Text: "Test Note",
		Date: time.Now(),
	}
	var err error
	suite.testModels[0].Id, err = suite.service.CreateModel(types.Model{})
	assert.NoError(suite.T(), err)
	suite.testModels[0].Notes = append(suite.testModels[0].Notes, note)
	err = suite.service.AddNote(suite.testModels[0])
	assert.NoError(suite.T(), err)
	mod, err := suite.service.GetModel(suite.testModels[0].Id)
	assert.NotEmpty(suite.T(), mod.Notes)
	assert.Equal(suite.T(), mod.Notes[0].Text, note.Text)
}
func (suite *ModelServiceTestSuite) TestGetGCodeMetaData() {
	gcode, err := suite.service.GetGCodeMetaData(filepath.Join("testdata/model1", suite.testModels[0].PrintFiles[0].Path))
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), gcode)
	assert.Equal(suite.T(), "PRUSA", gcode.GCodeType)
}

func TestModelServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ModelServiceTestSuite))
}

package model

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/mock"
	"ymir/pkg/api/model/types"
	"ymir/pkg/gcode"
)

// MockModelService is a mock implementation of the Service interface for testing.
type MockModelService struct {
	ModelService
	mock.Mock
	models []types.Model
}

func NewMockModelService() *MockModelService {
	return &MockModelService{
		models: getTestModels(),
	}
}

func (m *MockModelService) CreateModel(types.Model) (string, error) {
	return "", nil
}

func (m *MockModelService) ListModels() (map[string]types.Model, error) {
	models := map[string]types.Model{}
	for i := 0; i < len(m.models); i++ {
		models[m.models[i].Id] = m.models[i]
	}

	return models, nil
}

func (m MockModelService) GetModel(id string) (types.Model, error) {
	// Simulate returning a list of models for testing.
	fmt.Println(m.models[0])
	return m.models[1], nil
}

func (m *MockModelService) UpdateModel(types.Model) (err error) {
	return nil
}

func (m *MockModelService) DeleteModel(id string) error {
	return nil
}

func (m *MockModelService) ExportModel(path string, writer io.Writer) error {
	return nil
}

func (m *MockModelService) UploadFilesExistingModel(file multipart.File, filename string, basePath string) (string, error) {
	return "", nil
}

func (m *MockModelService) UploadFilesNewModel(file multipart.File, filename string) (string, error) {
	return "", nil
}

func (m *MockModelService) FetchModelImage(imagePath string) (imageBytes []byte, err error) {
	return nil, nil
}

func (m *MockModelService) FetchSTL(filepath string) (stlBytes []byte, err error) {
	return nil, nil
}

func (m *MockModelService) FetchSTLThumbnail(filepath string) (string, error) {
	return "", nil
}

func (m *MockModelService) AddNote(model types.Model) error {
	return nil
}

func (m *MockModelService) GetGCodeMetaData(path string) (gcode.GCodeMetaData, error) {
	return gcode.GCodeMetaData{}, nil
}

func (m *MockModelService) GetName() string {
	return ""
}

func getTestModels() []types.Model {
	models := []types.Model{}
	cwd, _ := os.Getwd()

	d := filepath.Join(cwd, "testdata")

	dirs, err := os.ReadDir(d)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return nil
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			modelFilePath := filepath.Join(d, dir.Name(), "model.json")
			if fileExists(modelFilePath) {
				model, err := readModelFile(modelFilePath)
				if err != nil {
					fmt.Printf("Error reading model.json in %s: %v\n", dir.Name(), err)
					continue
				}

				models = append(models, model)
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

	var model types.Model
	err = json.Unmarshal(data, &model)
	if err != nil {
		return types.Model{}, err
	}

	return model, nil
}

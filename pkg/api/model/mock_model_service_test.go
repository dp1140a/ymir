package model

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/mock"
	"ymir/pkg/gcode"
)

// MockModelService is a mock implementation of the Service interface for testing.
type MockModelService struct {
	ModelService
	mock.Mock
	models []Model
}

func NewMockModelService() *MockModelService {
	return &MockModelService{
		models: getTestModels(),
	}
}

func (m *MockModelService) CreateModel(Model) error {
	return nil
}

func (m *MockModelService) ListModels() ([]Model, error) {
	// Simulate returning a list of models for testing.
	return m.models, nil
}

func (m *MockModelService) GetModel(id string) (Model, error) {
	// Simulate returning a list of models for testing.
	fmt.Println(m.models[0])
	return m.models[1], nil
}

func (m *MockModelService) UpdateModel(Model) (rev string, err error) {
	return "2", nil
}

func (m *MockModelService) DeleteModel(id string, rev string) error {
	return nil
}

func (m *MockModelService) ExportModel(path string, writer io.Writer) error {
	return nil
}

func (m *MockModelService) UploadFile(file multipart.File, filename string, basePath string, isExistingModel bool) (key string, err error) {
	return "", nil
}

func (m *MockModelService) FetchModelImage(imagePath string) (imageBytes []byte, err error) {
	return nil, nil
}

func (m *MockModelService) FetchSTL(filepath string) (stlBytes []byte, err error) {
	return nil, nil
}

func (m *MockModelService) FetchSTLThumbnail(filepath string) string {
	return ""
}

func (m *MockModelService) AddNote(model Model) error {
	return nil
}

func (m *MockModelService) GetGCodeMetaData(path string) (gcode.GCodeMetaData, error) {
	return gcode.GCodeMetaData{}, nil
}

func (m MockModelService) GetName() string {
	return ""
}

func getTestModels() []Model {
	models := []Model{}
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

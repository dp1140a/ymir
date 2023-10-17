package model

import (
	"github.com/stretchr/testify/mock"
)

// MockModelService is a mock implementation of the Service interface for testing.
type MockModelService struct {
	ModelService
	mock.Mock
}

func (m *MockModelService) ListModels() ([]Model, error) {
	// Simulate returning a list of models for testing.
	return []Model{
		// Add sample Model data here for testing.
	}, nil
}

func (m MockModelService) GetName() string {
	return ""
}

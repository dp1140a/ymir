package model

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"ymir/pkg/api" // Replace with the actual import path for your project
)

type ModelHandlerTestSuite struct {
	suite.Suite
	models []Model
}

func (suite *ModelHandlerTestSuite) SetupSuite() {

}

func (suite *ModelHandlerTestSuite) SetupTest() {

}

func (suite *ModelHandlerTestSuite) TeardownTest() {

}

func (suite *ModelHandlerTestSuite) TeardownSuite() {

}

func TestModelHandler_ListAll(t *testing.T) {
	// Create a new ModelHandler with the mock service.
	mh := ModelHandler{
		Handler: api.Handler{
			Service: MockModelService{},
		},
	}

	// Create a request for testing (You may need to adjust this based on your actual API).
	req, err := http.NewRequest("GET", "/model", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Call the listAll function with the response recorder and request.
	mh.listAll(rr, req)

	// Check the response status code.
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check the response body (you may need to adjust this based on your actual response format).
	var models []Model
	if err := json.Unmarshal(rr.Body.Bytes(), &models); err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	// Add your assertions to check the response data and any other expectations.
	// For example, you can check the length of the returned models or specific model properties.

	// Example assertion: Check if at least one model is returned.
	if len(models) == 0 {
		t.Error("Expected at least one model to be returned, but got an empty list")
	}
}

func TestModelHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ModelHandlerTestSuite))
}

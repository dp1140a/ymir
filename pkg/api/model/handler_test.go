package model

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"ymir/pkg/api" // Replace with the actual import path for your project
	"ymir/pkg/api/model/types"
)

type ModelHandlerTestSuite struct {
	suite.Suite
	handler *ModelHandler
}

func (suite *ModelHandlerTestSuite) SetupSuite() {
	suite.handler = &ModelHandler{
		Handler: api.Handler{
			Service: NewMockModelService(),
		},
	}

	suite.handler.Routes = suite.handler.addRoutes()
}

func (suite *ModelHandlerTestSuite) SetupTest() {

}

func (suite *ModelHandlerTestSuite) TeardownTest() {

}

func (suite *ModelHandlerTestSuite) TeardownSuite() {

}

func (suite *ModelHandlerTestSuite) TestModelHandler_CreateModel_Good() {
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	go func() {
		defer multipartWriter.Close()
		nameField, err := multipartWriter.CreateFormField("displayName")
		if err != nil {
			suite.T().Error(err)
		}
		nameField.Write([]byte("TestModel"))
	}()

	req := httptest.NewRequest(http.MethodPost, "/model", pipeReader)
	req.Header.Set("content-type", multipartWriter.FormDataContentType())
	rr := httptest.NewRecorder()
	suite.handler.create(rr, req)

	// Check the response status code.
	assert.Equal(suite.T(), http.StatusCreated, rr.Code, "should be status coke 200")
	assert.Equal(suite.T(), "application/json", rr.Header().Get("content-type"))
	assert.Equal(suite.T(), "\"{'status': 'ok'}\"\n", rr.Body.String(), "response should be JSON")
}

func (suite *ModelHandlerTestSuite) TestModelHandler_CreateModel_Bad() {
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	go func() {
		defer multipartWriter.Close()
		nameField, err := multipartWriter.CreateFormField("displayName")
		if err != nil {
			suite.T().Error(err)
		}
		nameField.Write([]byte("TestModel"))
	}()

	req := httptest.NewRequest(http.MethodPost, "/model", pipeReader)
	req.Header.Set("content-type", "text/plain") // <-- This will make it fail
	rr := httptest.NewRecorder()
	suite.handler.create(rr, req)

	// Check the response status code.
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code, "should be status 400")
	assert.Equal(suite.T(), "text/plain; charset=utf-8", rr.Header().Get("content-type"))
}

func (suite *ModelHandlerTestSuite) TestModelHandler_ListAll() {
	req, err := http.NewRequest("GET", "/model", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	rr := httptest.NewRecorder()
	suite.handler.listAll(rr, req)

	// Check the response status code.
	assert.Equal(suite.T(), http.StatusOK, rr.Code, "should be status coke 200")

	// Check the response body (you may need to adjust this based on your actual response format).
	var models map[string]types.Model
	if err := json.Unmarshal(rr.Body.Bytes(), &models); err != nil {
		suite.T().Errorf("Failed to unmarshal response body: %v", err)
	}

	assert.Len(suite.T(), models, 2, "should be 2 models")
	assert.NotNil(suite.T(), models, "should not be nil")
}

func (suite *ModelHandlerTestSuite) TestModelHandler_InspectModel() {
	id := "5e469849725bc276" //Should be model2 in testdata
	req, err := http.NewRequest("GET", "/model/{id}", nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	suite.handler.inspect(rr, req)

	// Check the response status code.
	assert.Equal(suite.T(), http.StatusOK, rr.Code, "should be status coke 200")

	// Check the response body (you may need to adjust this based on your actual response format).
	var model types.Model
	if err := json.Unmarshal(rr.Body.Bytes(), &model); err != nil {
		suite.T().Errorf("Failed to unmarshal response body: %v", err)
	}

	assert.NotNil(suite.T(), model, "should not be nil")
	assert.Equal(suite.T(), id, model.Id, "It should be 5e469849725bc276")
	assert.Equal(suite.T(), "Filament Grommet", model.DisplayName, "name should be Filament Grommet")
	assert.Len(suite.T(), model.Tags, 2, "Should have 2 tags")
}

// DELETE /model/{id}?rev
func (suite *ModelHandlerTestSuite) TestModelHandler_DeleteModel() {
	id := "5e469849725bc276" //Should be model2 in testdata
	req, err := http.NewRequest("DELETE", "/model/{id}", nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	suite.handler.delete(rr, req)

	assert.Equal(suite.T(), http.StatusNoContent, rr.Code, "should be status 400")

}

func TestModelHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ModelHandlerTestSuite))
}

package printer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"ymir/pkg/api"
	"ymir/pkg/api/printer/types"
)

type PrinterHandlerTestSuite struct {
	suite.Suite
	handler *PrinterHandler
}

func (suite *PrinterHandlerTestSuite) SetupSuite() {
	suite.handler = &PrinterHandler{
		Handler: api.Handler{
			Service: NewMockPrinterService(),
		},
	}
}

func (suite *PrinterHandlerTestSuite) SetupTest() {

}

func (suite *PrinterHandlerTestSuite) TeardownTest() {

}

func (suite *PrinterHandlerTestSuite) TeardownSuite() {

}

func (suite *PrinterHandlerTestSuite) TestPrinterHandler_CreatePrinter_Good() {
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	go func() {
		defer func(multipartWriter *multipart.Writer) {
			err := multipartWriter.Close()
			if err != nil {

			}
		}(multipartWriter)
		nameField, err := multipartWriter.CreateFormField("displayName")
		if err != nil {
			suite.T().Error(err)
		}
		_, err = nameField.Write([]byte("TestPrinter"))
		if err != nil {
			return
		}
	}()

	req := httptest.NewRequest(http.MethodPost, "/printer", pipeReader)
	req.Header.Set("content-type", multipartWriter.FormDataContentType())
	rr := httptest.NewRecorder()
	suite.handler.create(rr, req)

	// Check the response status code.
	assert.Equal(suite.T(), http.StatusCreated, rr.Code, "should be status code 200")
	assert.Equal(suite.T(), "application/json", rr.Header().Get("content-type"))
	assert.Equal(suite.T(), "\"{'status': 'ok'}\"\n", rr.Body.String(), "response should be JSON")
}

func (suite *PrinterHandlerTestSuite) TestPrinterHandler_CreatePrinter_Bad() {
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)

	go func() {
		defer func(multipartWriter *multipart.Writer) {
			err := multipartWriter.Close()
			if err != nil {

			}
		}(multipartWriter)
		nameField, err := multipartWriter.CreateFormField("displayName")
		if err != nil {
			suite.T().Error(err)
		}
		_, err = nameField.Write([]byte("TestPrinter"))
		if err != nil {
			return
		}
	}()

	req := httptest.NewRequest(http.MethodPost, "/printer", pipeReader)
	req.Header.Set("content-type", "text/plain") // <-- This will make it fail
	rr := httptest.NewRecorder()
	suite.handler.create(rr, req)

	// Check the response status code.
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code, "should be status 400")
	assert.Equal(suite.T(), "text/plain; charset=utf-8", rr.Header().Get("content-type"))
}

func (suite *PrinterHandlerTestSuite) TestPrinterHandler_ListAll() {
	req, err := http.NewRequest("GET", "/printer", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	rr := httptest.NewRecorder()
	suite.handler.listAll(rr, req)

	// Check the response status code.
	assert.Equal(suite.T(), http.StatusOK, rr.Code, "should be status coke 200")

	// Check the response body (you may need to adjust this based on your actual response format).
	var printers map[string]types.Printer
	fmt.Println(rr.Body.String())
	if err := json.Unmarshal(rr.Body.Bytes(), &printers); err != nil {
		suite.T().Errorf("Failed to unmarshal response body: %v", err)
	}

	assert.Len(suite.T(), printers, 2, "should be 2 printers")
	assert.NotNil(suite.T(), printers, "should not be nil")
}

func (suite *PrinterHandlerTestSuite) TestPrinterHandler_InspectPrinter() {
	id := "test-0" //Should be printer2 in testdata
	req, err := http.NewRequest("GET", "/printer/{id}", nil)
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
	var printer types.Printer
	if err := json.Unmarshal(rr.Body.Bytes(), &printer); err != nil {
		suite.T().Errorf("Failed to unmarshal response body: %v", err)
	}

	assert.NotNil(suite.T(), printer, "should not be nil")
	assert.Equal(suite.T(), id, printer.Id, "test-0")
	assert.Equal(suite.T(), "test_1", printer.PrinterName, "name should be test_1")
}

// DELETE /printer/{id}?rev
func (suite *PrinterHandlerTestSuite) TestPrinterHandler_DeletePrinter() {
	id := "5e469849725bc276" //Should be printer2 in testdata
	req, err := http.NewRequest("DELETE", "/printer/{id}", nil)
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

func TestPrinterHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PrinterHandlerTestSuite))
}

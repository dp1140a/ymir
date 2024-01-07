package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"ymir/pkg/api/printer/store"
	"ymir/pkg/api/printer/types"
	"ymir/pkg/logger"
)

const (
	TEST_DIR = "TEST_DIR"
)

type PrintersServiceTestSuite struct {
	suite.Suite
	service      PrinterService
	testPrinters []types.Printer
}

func (suite *PrintersServiceTestSuite) SetupSuite() {
	fmt.Println("SetupSuite()")
	viper.SetConfigType("toml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var tomlExample = []byte(`
[logging]
logLevel="DEBUG"

[datastore]
dbFile = "test.db"

[printers]
printersDir="TEST_DIR/printers"
`)

	err := viper.ReadConfig(bytes.NewBuffer(tomlExample))
	if err != nil {
		suite.T().Errorf("Error: %v", err)
	}

	suite.service = NewPrinterService().(PrinterService)
	err = logger.InitLogger()
	if err != nil {
		fmt.Println(err.Error())
		suite.T().Fatal(err.Error())
	}
	suite.testPrinters = getTestPrinters(2)
}

func (suite *PrintersServiceTestSuite) AfterTest(suiteName, testName string) {}

func (suite *PrintersServiceTestSuite) BeforeTest(suiteName, testName string) {
}

func (suite *PrintersServiceTestSuite) SetupTest() {
	fmt.Println("SetupTest")
	err := suite.service.printerStore.(store.PrinterStore).Truncate()
	assert.NoError(suite.T(), err)
	err = suite.service.printerStore.Create(suite.testPrinters[0])
	assert.NoError(suite.T(), err, "should be no error on test setup")
}

func (suite *PrintersServiceTestSuite) TearDownTest() {
}

func (suite *PrintersServiceTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite()")
	suite.T().Cleanup(func() {
	})
	err := os.RemoveAll(TEST_DIR)
	err = os.Remove("test.db")

	if err != nil {
		log.Fatal(err)
	}
}

func (suite *PrintersServiceTestSuite) TestNewPrinterService() {
	assert.NotNil(suite.T(), suite.service, "Should not be nil")
}

func (suite *PrintersServiceTestSuite) TestPrinterServiceGetName() {
	assert.Equal(suite.T(), "Printers", suite.service.GetName())
}

func (suite *PrintersServiceTestSuite) TestCreatePrinter() {
	id, err := suite.service.CreatePrinter(types.Printer{})
	suite.testPrinters[0].Id = id
	pri, err := suite.service.GetPrinter(id)
	assert.Equal(suite.T(), suite.testPrinters[0].Id, pri.Id, "should be same")
	assert.NoError(suite.T(), err)
}

func (suite *PrintersServiceTestSuite) TestUpdatePrinter() {
	newPri := suite.testPrinters[1]
	newName := "testUpdate"
	newPri.PrinterName = newName
	err := suite.service.UpdatePrinter(newPri)
	assert.NoError(suite.T(), err)
	pri, err := suite.service.GetPrinter(suite.testPrinters[1].Id)
	assert.NoError(suite.T(), err, "should be no error")
	assert.Equal(suite.T(), newName, pri.PrinterName, "should be equal")
}

func (suite *PrintersServiceTestSuite) TestListPrinters() {
	Printers, err := suite.service.ListPrinters()
	assert.NoError(suite.T(), err, "should be no error")
	assert.Len(suite.T(), Printers, 1, "should be 1")
	assert.IsType(suite.T(), map[string]types.Printer{}, Printers, "should be type []Printer.Printer{}")
	assert.IsType(suite.T(), types.Printer{}, Printers[suite.testPrinters[0].Id], "should be of type Printer")
	assert.Equal(suite.T(), suite.testPrinters[0], Printers[suite.testPrinters[0].Id], "should be equal")
}

func (suite *PrintersServiceTestSuite) TestGetPrinter() {
	pri, err := suite.service.GetPrinter(suite.testPrinters[0].Id)
	assert.NoError(suite.T(), err, "should be no error")
	assert.IsType(suite.T(), types.Printer{}, pri, "should be of type Printer")
	assert.Equal(suite.T(), suite.testPrinters[0], pri, "should be equal")
}

func (suite *PrintersServiceTestSuite) TestDeletePrinter() {
	err := suite.service.DeletePrinter(suite.testPrinters[1].Id)
	assert.NoError(suite.T(), err, "should be no error")
	_, err = suite.service.GetPrinter(suite.testPrinters[1].Id)
	assert.Error(suite.T(), err, "should be error")
}

/*
Utility Functions
*/
func getTestPrinters(num int) (printers []types.Printer) {
	for i := 0; i < num; i++ {
		printer := types.Printer{}
		_ = json.Unmarshal([]byte(types.TestPrinter), &printer)
		printer.PrinterName = fmt.Sprintf("test_%v", i+1)
		printer.Id = fmt.Sprintf("test-%v", i)
		printers = append(printers, printer)
	}

	return printers
}

func TestPrintersServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PrintersServiceTestSuite))
}

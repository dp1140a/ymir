package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"ymir/pkg/api/printer/types"
)

const (
	TEST_DB = "test.db"
)

type PrinterStoreTestSuite struct {
	suite.Suite
	store        PrinterStore
	testPrinters []types.Printer
}

func (suite *PrinterStoreTestSuite) SetupSuite() {
	fmt.Println("SetupSuite()")
	viper.SetConfigType("toml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var tomlExample = []byte(`
[datastore]
dbFile = "test.db"
`)

	err := viper.ReadConfig(bytes.NewBuffer(tomlExample))
	if err != nil {
		suite.T().Errorf("Error: %v", err)
	}
	suite.store = NewPrinterDataStore().(PrinterStore)
	suite.testPrinters = getTestPrinters(2)
}

func (suite *PrinterStoreTestSuite) SetupTest() {
}

func (suite *PrinterStoreTestSuite) TearDownTest() {
}

func (suite *PrinterStoreTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite()")
	if _, err := os.Stat(TEST_DB); errors.Is(err, os.ErrNotExist) {
		fmt.Println("DB Does NOT exist")
	}
	err := os.Remove(TEST_DB)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func (suite *PrinterStoreTestSuite) TestNewPrinterStore() {
	assert.NotNil(suite.T(), suite.store, "Should not be nil")
	assert.Equal(suite.T(), TEST_DB, suite.store.ds.GetDB().Path(), "should be test.db")
	assert.Equal(suite.T(), "boltDB", suite.store.ds.GetType(), "Should be boltDB")
}

func (suite *PrinterStoreTestSuite) TestPrinterCreate() {
	for i := 0; i < len(suite.testPrinters); i++ {
		err := suite.store.Create(suite.testPrinters[i])
		assert.NoError(suite.T(), err)
		numPrinters, _ := suite.store.ds.GetNumKeys(PRINTERS_BUCKET)
		assert.Equal(suite.T(), i+1, numPrinters)
	}
}

func (suite *PrinterStoreTestSuite) TestPrinterUpdate() {
	newPri := suite.testPrinters[1]
	newName := "testUpdate"
	newPri.PrinterName = newName
	err := suite.store.Update(newPri)
	assert.NoError(suite.T(), err)
	pri, err := suite.store.Inspect(suite.testPrinters[1].Id)
	assert.NoError(suite.T(), err, "should be no error")
	assert.Equal(suite.T(), newName, pri.PrinterName, "should be equal")
}

func (suite *PrinterStoreTestSuite) TestPrinterDelete() {
	err := suite.store.Delete(suite.testPrinters[1].Id)
	assert.NoError(suite.T(), err, "should be no error")
	numPrinters, _ := suite.store.ds.GetNumKeys(PRINTERS_BUCKET)
	assert.Equal(suite.T(), 1, numPrinters, "should be 0")
}

func (suite *PrinterStoreTestSuite) TestPrinterList() {
	Printers, err := suite.store.List()
	assert.NoError(suite.T(), err, "should be no error")
	assert.Len(suite.T(), Printers, 1, "should be 1")
	assert.IsType(suite.T(), map[string]types.Printer{}, Printers, "should be type []Printer.Printer{}")
	assert.IsType(suite.T(), types.Printer{}, Printers[suite.testPrinters[0].Id], "should be of type Printer")
	assert.Equal(suite.T(), suite.testPrinters[0], Printers[suite.testPrinters[0].Id], "should be equal")
}

func (suite *PrinterStoreTestSuite) TestPrinterInspect() {
	mod, err := suite.store.Inspect(suite.testPrinters[0].Id)
	assert.NoError(suite.T(), err, "should be no error")
	assert.IsType(suite.T(), types.Printer{}, mod, "should be of type Printer")
	assert.Equal(suite.T(), suite.testPrinters[0], mod, "should be equal")
}

func TestPrinterStoreTestSuite(t *testing.T) {
	suite.Run(t, new(PrinterStoreTestSuite))
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

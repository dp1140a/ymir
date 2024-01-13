package boltdatastore

import (
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoltDataStoreTestSuite struct {
	suite.Suite
	store *BoltDBDataStore
}

func (suite *BoltDataStoreTestSuite) SetupSuite() {
	fmt.Println("SetupSuite()")
	suite.store = NewBoltDBDatastore(NewBoltDBDataStoreConfig())
}

func (suite *BoltDataStoreTestSuite) SetupTest() {

}

func (suite *BoltDataStoreTestSuite) TearDownTest() {

}

func (suite *BoltDataStoreTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite()")
	suite.T().Cleanup(func() {
		suite.store.Close()
	})
	e := os.Remove("ymir.db")
	if e != nil {
		log.Fatal(e)
	}
}

func (suite *BoltDataStoreTestSuite) TestNewBoltDBDataStore_Default() {
	assert.NotNil(suite.T(), BoltDBDataStore{}, suite.store, "Should not be nil")
	assert.NotNil(suite.T(), suite.store.ds, "Datastore.db should not be nil")
	assert.Equal(suite.T(), "ymir.db", suite.store.ds.Path(), "should be ymir.db")
	assert.Equal(suite.T(), "boltDB", suite.store.GetType(), "Should be boltDB")
}

func (suite *BoltDataStoreTestSuite) TestGetDB() {
	assert.NotNil(suite.T(), BoltDBDataStore{}, suite.store, "Should not be nil")
	assert.NotNil(suite.T(), suite.store.GetDB(), "Datastore.db should not be nil")
	assert.Equal(suite.T(), "ymir.db", suite.store.GetDB().Path(), "should be ymir.db")
}

func (suite *BoltDataStoreTestSuite) TestCreateBucket() {
	err := suite.store.CreateBucket("test")
	assert.NoError(suite.T(), err, "Should have no error")
	assert.True(suite.T(), suite.store.BucketExists("test"), "bucket should exist")
}

func TestBoltDataStoreTestSuite(t *testing.T) {
	suite.Run(t, new(BoltDataStoreTestSuite))
}

package boltdatastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"ymir/pkg/db"
)

var (
	lock = &sync.Mutex{}
)

var instance *bolt.DB

type BoltDBDataStore struct {
	db.Datastore
	ds    *bolt.DB
	stats bolt.Stats
}

func NewBoltDBDatastore(config *BoltDBDataStoreConfig) (datastore *BoltDBDataStore) {
	b := &BoltDBDataStore{}
	err := getInstance(config.DBFile)
	if err != nil {
		log.Error("Could not Open DB")
	}
	b.ds = instance
	return b
}

func getInstance(dbFile string) error {
	if _, err := os.Stat(filepath.Dir(dbFile)); os.IsNotExist(err) {
		err := os.Mkdir(filepath.Dir(dbFile), 0750)
		if err != nil {
			log.Errorf("cant find or create db directory %v", err)
			return err
		}
	}
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		log.Info("Creating single instance for db ", dbFile)
		var err error
		instance, err = bolt.Open(dbFile, 0644, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal(err)
			return err
		}
	} else {
		log.Info("returning instance: ", instance)
	}
	return nil
}

/*
Satisfies Store Interface
*/
func (bds *BoltDBDataStore) Close() error {
	return bds.ds.Close()
}

func (bds *BoltDBDataStore) GetDB() *bolt.DB {
	return bds.ds
}

func (bds *BoltDBDataStore) GetType() string {
	return "boltDB"
}

func (bds *BoltDBDataStore) BucketExists(bucket string) bool {
	err := bds.ds.Batch(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte(bucket))

		if b != nil {
			return nil
		}
		return errors.New("no bucket found")
	})

	if err != nil {
		return false
	} else {
		return true
	}
}

func (bds *BoltDBDataStore) CreateBucket(bucket string) error {
	err := bds.ds.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return err
}

func (bds *BoltDBDataStore) GetNumKeys(bucket string) (int, error) {
	numKeys := 0
	err := bds.ds.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			numKeys++
		}
		return nil
	})
	return numKeys, err
}

func (bds *BoltDBDataStore) Stats() (stats []byte) {
	// Grab the initial stats.
	prev := bds.stats
	for {
		// Wait for 10s.
		time.Sleep(10 * time.Second)

		// Grab the current stats and diff them.
		stats := bds.ds.Stats()
		diff := stats.Sub(&prev)

		// Encode stats to JSON and print to STDERR.
		//json.NewEncoder(os.Stderr).Encode(diff)
		jsonStr, _ := json.Marshal(diff)

		// Save stats for the next loop.
		bds.stats = stats
		return jsonStr
	}
}

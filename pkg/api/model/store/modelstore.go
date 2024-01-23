package store

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"ymir/pkg/api/model/types"
	db "ymir/pkg/db/boltdatastore"
)

const (
	MODELS_BUCKET = "models"
)

type ModelStoreIFace interface {
	Create(model types.Model) (err error)
	Update(model types.Model) (err error)
	Delete(id string) (err error)
	List() (models map[string]types.Model, err error)
	Inspect(id string) (model types.Model, err error)
	Truncate() (err error)
}

type ModelStore struct {
	ModelStoreIFace
	bucket string
	ds     db.BoltDBDataStore
}

func NewModelDataStore() (store ModelStoreIFace) {
	config := db.NewBoltDBDataStoreConfig()
	fmt.Println(config.DBFile)
	d := ModelStore{
		ds: *db.NewBoltDBDatastore(config),
	}
	err := createBucketIfNotExists(&d.ds)
	if err != nil {
		log.Error("could not create bucket:")
		return nil
	}
	return d
}

func createBucketIfNotExists(ds *db.BoltDBDataStore) error {
	//Create Bucket if  not exists
	err := ds.CreateBucket(MODELS_BUCKET)
	if err != nil {
		return nil
	}
	return nil
}

func (ms ModelStore) Create(model types.Model) (err error) {
	return ms.Update(model)
}

func (ms ModelStore) Update(model types.Model) (err error) {
	err = ms.ds.GetDB().Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(MODELS_BUCKET))
		mJson, err := json.Marshal(model)
		err = b.Put([]byte(model.Id), mJson)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	return err
}

func (ms ModelStore) Delete(id string) (err error) {
	err = ms.ds.GetDB().Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(MODELS_BUCKET))
		err = b.Delete([]byte(id))
		return err
	})
	return err
}

func (ms ModelStore) Truncate() (err error) {
	err = ms.ds.GetDB().Update(func(tx *bolt.Tx) (err error) {
		err = tx.DeleteBucket([]byte(MODELS_BUCKET))
		return err
	})

	err = createBucketIfNotExists(&ms.ds)
	return
}

func (ms ModelStore) List() (map[string]types.Model, error) {
	models := map[string]types.Model{}
	err := ms.ds.GetDB().View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(MODELS_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			m := types.Model{}
			err := json.Unmarshal(v, &m)
			if err != nil {
				log.Error("error unmarshalling model")
				return err
			}
			models[string(k)] = m
		}
		return err
	})
	return models, err
}

func (ms ModelStore) Inspect(id string) (mod types.Model, err error) {
	err = ms.ds.GetDB().View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(MODELS_BUCKET))
		m := b.Get([]byte(id))
		mod = types.Model{}
		err = json.Unmarshal(m, &mod)
		if err != nil {
			return err
		}
		return err
	})
	return mod, nil
}

func (ms ModelStore) NumModels() int {
	numModels, _ := ms.ds.GetNumKeys(MODELS_BUCKET)
	return numModels
}

func (ms ModelStore) Close() error {
	return ms.ds.Close()
}

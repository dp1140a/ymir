package store

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"ymir/pkg/api/printer/types"
	db "ymir/pkg/db/boltdatastore"
)

const (
	PRINTERS_BUCKET = "printers"
)

type PrinterStoreIFace interface {
	Create(printer types.Printer) (err error)
	Update(printer types.Printer) (err error)
	Delete(id string) (err error)
	List() (printers map[string]types.Printer, err error)
	Inspect(id string) (printer types.Printer, err error)
}

type PrinterStore struct {
	PrinterStoreIFace
	bucket string
	ds     db.BoltDBDataStore
}

func NewPrinterDataStore() (store PrinterStoreIFace) {
	config := db.NewBoltDBDataStoreConfig()
	d := PrinterStore{
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
	err := ds.CreateBucket(PRINTERS_BUCKET)
	if err != nil {
		return nil
	}
	return nil
}

func (ms PrinterStore) Create(printer types.Printer) (err error) {
	return ms.Update(printer)
}

func (ms PrinterStore) Update(printer types.Printer) (err error) {
	err = ms.ds.GetDB().Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(PRINTERS_BUCKET))
		mJson, err := json.Marshal(printer)
		err = b.Put([]byte(printer.Id), mJson)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	return err
}

func (ms PrinterStore) Delete(id string) (err error) {
	err = ms.ds.GetDB().Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(PRINTERS_BUCKET))
		err = b.Delete([]byte(id))
		return err
	})
	return err
}

func (ms PrinterStore) Truncate() (err error) {
	err = ms.ds.GetDB().Update(func(tx *bolt.Tx) (err error) {
		err = tx.DeleteBucket([]byte(PRINTERS_BUCKET))
		return err
	})

	err = createBucketIfNotExists(&ms.ds)
	return
}

func (ms PrinterStore) List() (map[string]types.Printer, error) {
	printers := map[string]types.Printer{}
	err := ms.ds.GetDB().View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(PRINTERS_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			m := types.Printer{}
			err := json.Unmarshal(v, &m)
			if err != nil {
				log.Error("error unmarshalling printer")
				return err
			}
			printers[string(k)] = m
		}
		return err
	})
	return printers, err
}

func (ms PrinterStore) Inspect(id string) (mod types.Printer, err error) {
	err = ms.ds.GetDB().View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(PRINTERS_BUCKET))
		m := b.Get([]byte(id))
		mod = types.Printer{}
		err = json.Unmarshal(m, &mod)
		if err != nil {
			return err
		}
		return err
	})
	return mod, nil
}

func (ms PrinterStore) NumPrinters() int {
	numPrinters, _ := ms.ds.GetNumKeys(PRINTERS_BUCKET)
	return numPrinters
}

func (ms PrinterStore) Close() error {
	return ms.ds.Close()
}

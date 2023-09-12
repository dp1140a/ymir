package db

import (
	"context"
	"fmt"

	"github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
	log "github.com/sirupsen/logrus"
)

type DB struct {
	config *DBConfig
	Client *kivik.Client //https://github.com/go-kivik
}

func NewDB() *DB {
	db := &DB{
		config: NewDBConfig(),
	}
	//fmt.Println(db.config.StringToml())
	return db
}

func (db *DB) GetDB() *kivik.DB {
	return db.Client.DB(context.TODO(), db.config.DBName)
}

func (db *DB) Drop() {
	err := db.Client.DestroyDB(context.TODO(), db.config.DBName)
	if err != nil {
		log.Error(err.Error())
	}
	log.Infof("Dropped DB %s", db.config.DBName)
}

func (db *DB) Truncate() {
	db.Drop()
	db.CreateDB()
}

func (db *DB) CreateDB() {
	err := db.Client.CreateDB(context.TODO(), db.config.DBName)
	if err != nil {
		log.Error(err.Error())
	}
	log.Infof("Created DB %s", db.config.DBName)
}

func (db *DB) Connect() {
	client, err := kivik.New("couch", fmt.Sprintf("http://%s:%s", db.config.Host, db.config.Port))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Authenticate(context.TODO(), couchdb.BasicAuth(db.config.Username, db.config.Password))
	if err != nil {
		log.Fatal(err)
	}

	up, err := client.Ping(context.TODO())
	if err != nil {
		log.Fatalf("DB Connection Error: %v", err)
	} else {
		log.Infof("DB at %s is up: %v", client.DSN(), up)
		db.Client = client
	}
}

package model

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
	"ymir/pkg/db"
)

type ModelService struct {
	name      string
	DataStore *db.DB
}

func NewModelService() (modelService api.Service) {
	ms := ModelService{
		name: "Model",
	}

	ds := db.NewDB()
	ds.Connect()
	ms.DataStore = ds

	return ms
}

func (ms ModelService) GetName() (name string) {
	return ms.name
}

func (ms ModelService) GetModel(id string) (model Model, err error) {
	row := ms.DataStore.GetDB().Get(context.TODO(), id)
	model = Model{}
	if err = row.ScanDoc(&model); err != nil {
		return model, err
	}

	return model, nil
}

func (ms ModelService) CreateModel(model Model) (err error) {
	ctx := context.TODO()
	docId := uuid.New().String()

	rev, err := ms.DataStore.GetDB().Put(ctx, docId, model)
	if err != nil {
		return err
	} else {
		log.Infof("created model %v in db with rev %v", model.Id, rev)
		return nil
	}
}

func (ms ModelService) UpdateModel(model Model) (err error) {
	ctx := context.TODO()

	rev, err := ms.DataStore.GetDB().Put(ctx, model.Id, model)
	if err != nil {
		return err
	} else {
		log.Infof("updated model %v in db with rev %v", model.Id, rev)
		return nil
	}
}

func (ms ModelService) DeleteModel(id string, rev string) error {
	newRev, err := ms.DataStore.GetDB().Delete(context.TODO(), id, rev)
	if err != nil {
		return err
	}
	log.Infof("deleted model %v in db with newRev %v", id, newRev)
	return nil
}

func (ms ModelService) ListModels() ([]Model, error) {
	query := `{
		"selector": {
			        "displayName": {"$regex": ".+"}
		},
		"fields": ["_id", "_rev", "displayName"]
	}`
	var q interface{}
	_ = json.Unmarshal([]byte(query), q)
	rows, err := ms.DataStore.GetDB().Find(context.TODO(), query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	docs := []Model{}
	for rows.Next() {
		doc := &Model{}
		err = rows.ScanDoc(doc)
		if err != nil {
			log.Errorf("rows error: %v\n", err)
			return nil, err
		}
		docs = append(docs, *doc)
	}

	return docs, nil
}

func (ms ModelService) GetModelsByTag(tag string) ([]Model, error) {

	query := `{
	   "selector": {
		  "tags": {
			 "$elemMatch": {
				"$eq": "---"
			 }
		  }
	   }
	}`

	query = strings.Replace(query, "---", tag, 1)
	//fmt.Println(query)
	var q interface{}
	_ = json.Unmarshal([]byte(query), q)
	rows, err := ms.DataStore.GetDB().Find(context.TODO(), query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	docs := []Model{}
	for rows.Next() {
		doc := &Model{}
		err = rows.ScanDoc(doc)
		if err != nil {
			log.Errorf("rows error: %v\n", err)
			return nil, err
		}
		docs = append(docs, *doc)
	}

	return docs, nil
}

package printer

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
	"ymir/pkg/db"
	"ymir/pkg/utils"
)

type PrintersService struct {
	name      string
	DataStore *db.DB
	config    *PrintersConfig
}

func NewPrinterService() (printersService api.Service) {
	ps := PrintersService{
		name:   "Printers",
		config: NewPrintersConfig(),
	}

	ds := db.NewDB()
	ds.Connect()
	ps.DataStore = ds

	utils.MakeDirIfNotExists(ps.config.PrintersDir)

	return ps
}

func (ps PrintersService) GetName() (name string) {
	return ps.name
}

func (ps PrintersService) GetConfig() *PrintersConfig {
	return ps.config
}

func (ps PrintersService) ListPrinters() ([]Printer, error) {

	query := `{
		"selector": {
			"printer": {"$regex": ".+"}
		},
		"fields": ["_id", "_rev", "url", "description", "displayName", "tags" ]
	}`
	var q interface{}
	_ = json.Unmarshal([]byte(query), q)
	rows, err := ps.DataStore.GetDB().Find(context.TODO(), query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var docs []Printer
	for rows.Next() {
		doc := &Printer{}
		err = rows.ScanDoc(doc)
		if err != nil {
			log.Errorf("rows error: %v\n", err)
			return nil, err
		}
		docs = append(docs, *doc)
	}

	return docs, nil
}

func (ps PrintersService) GetPrinter(id string) (printer Printer, err error) {
	return Printer{}, nil
}

func (ps PrintersService) CreatePrinter(Printer Printer) (err error) {

	return nil
}

func (ps PrintersService) UpdatePrinter(Printer Printer) (err error) {

	return nil
}

func (ps PrintersService) DeletePrinter(id string, rev string) error {

	return nil
}

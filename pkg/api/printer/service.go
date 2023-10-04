package printer

import (
	"context"
	"encoding/json"
	"fmt"

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
			"printerName": {"$regex": ".+"}
		},
		"fields": ["_id", "_rev", "url", "printerName", "tags", "location", "apiKey", "type"]
	}`
	var q interface{}
	_ = json.Unmarshal([]byte(query), q)
	rows, err := ps.DataStore.GetDB().Find(context.TODO(), query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	docs := []Printer{}
	for rows.Next() {
		doc := &Printer{}
		err = rows.ScanDoc(doc)
		if err != nil {
			log.Errorf("rows error: %v\n", err)
			return nil, err
		}
		docs = append(docs, *doc)
	}

	if log.GetLevel() == log.DebugLevel {
		fmt.Printf("%v printers returned\n", len(docs))
	}

	return docs, nil
}

func (ps PrintersService) GetPrinter(id string) (printer Printer, err error) {
	row := ps.DataStore.GetDB().Get(context.TODO(), id)
	printer = Printer{}
	if err = row.ScanDoc(&printer); err != nil {
		log.Error(err)
		return printer, err
	}

	return printer, nil
}

func (ps PrintersService) CreatePrinter(printer Printer) (err error) {
	ctx := context.TODO()
	printer.Id = utils.GenId()
	if log.GetLevel() == log.DebugLevel {
		fmt.Println(printer.Json())
	}

	docId, rev, err := ps.DataStore.GetDB().CreateDoc(ctx, printer)
	if err != nil {
		log.Errorf("error adding printer: %v", err)
		return err
	}

	log.Infof("added printer %v in db with docid %v and rev %v", printer.Id, docId, rev)
	return nil
}

func (ps PrintersService) UpdatePrinter(printer Printer) (rev string, err error) {
	ctx := context.TODO()

	rev, err = ps.DataStore.GetDB().Put(ctx, printer.Id, printer)
	if err != nil {
		log.Errorf("error updating printer: %v", err)
		return "", err
	}
	log.Infof("updated printer %v in db with rev %v", printer.Id, rev)
	return rev, nil
}

func (ps PrintersService) DeletePrinter(id string, rev string) error {
	_, err := ps.DataStore.GetDB().Delete(context.TODO(), id, rev)
	if err != nil {
		log.Errorf("error deleting printer: %v", err)
		return err
	}

	return nil
}

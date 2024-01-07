package printer

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
	"ymir/pkg/api/printer/store"
	"ymir/pkg/api/printer/types"
	"ymir/pkg/utils"
)

type PrinterServiceIface interface {
	api.Service
	CreatePrinter(printer types.Printer) (id string, err error)
	UpdatePrinter(printer types.Printer) (err error)
	DeletePrinter(id string) error
	GetPrinter(id string) (types.Printer, error)
	ListPrinters() (map[string]types.Printer, error)
}

type PrinterService struct {
	PrinterServiceIface
	name         string
	printerStore store.PrinterStoreIFace
	config       *PrintersConfig
}

func NewPrinterService() api.Service {
	ps := PrinterService{
		name:         "Printers",
		config:       NewPrintersConfig(),
		printerStore: store.NewPrinterDataStore(),
	}

	err := utils.MakeDirIfNotExists(ps.config.PrintersDir)
	if err != nil {
		return nil

	}
	return ps
}

func (ps PrinterService) GetName() (name string) {
	return ps.name
}

func (ps PrinterService) GetConfig() *PrintersConfig {
	return ps.config
}

func (ps PrinterService) ListPrinters() (printers map[string]types.Printer, err error) {
	printers, err = ps.printerStore.List()
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (ps PrinterService) GetPrinter(id string) (printer types.Printer, err error) {
	printer, err = ps.printerStore.Inspect(id)
	if err != nil {
		log.Errorf("error retrieving printerl: %v", err)
		return
	} else if printer.Id == "" {
		err = errors.New(fmt.Sprintf("printer with id: %v does not exist", id))
	}
	return
}

func (ps PrinterService) CreatePrinter(printer types.Printer) (id string, err error) {
	printer.Id = utils.GenId()
	if log.GetLevel() == log.DebugLevel {
		fmt.Println(printer.Json())
	}

	err = ps.printerStore.Create(printer)
	if err != nil {
		log.Error(err)
		return
	} else {
		log.Infof("created printer %v in db", printer.Id)
		return printer.Id, nil
	}
}

func (ps PrinterService) UpdatePrinter(printer types.Printer) (err error) {
	err = ps.printerStore.Update(printer)
	if err != nil {
		log.Error(err)
		return
	} else {
		log.Infof("updated printer %v in db", printer.Id)
		return
	}
}

func (ps PrinterService) DeletePrinter(id string) (err error) {
	err = ps.printerStore.Delete(id)
	if err != nil {
		return err
	}
	log.Infof("deleted printer %v in db", id)
	return
}

package admin

import (
	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
	ms "ymir/pkg/api/model/store"
	model "ymir/pkg/api/model/types"
	ps "ymir/pkg/api/printer/store"
	printer "ymir/pkg/api/printer/types"
)

type AdminServiceIface interface {
	api.Service
	ListModels() (models map[string]model.Model, err error)
	ListPrinters() (printers map[string]printer.Printer, err error)
	TruncateModels() error
	TruncatePrinters() error
}

type AdminService struct {
	AdminServiceIface
	name         string
	modelStore   ms.ModelStoreIFace
	printerStore ps.PrinterStoreIFace
	//config     *AdminConfig
}

func NewAdminService() api.Service {
	as := AdminService{
		name:         "Admin",
		modelStore:   ms.NewModelDataStore(),
		printerStore: ps.NewPrinterDataStore(),
	}
	return as
}

func (as AdminService) GetName() (name string) {
	return as.name
}

func (as AdminService) ListModels() (models map[string]model.Model, err error) {
	return as.modelStore.List()
}

func (as AdminService) ListPrinters() (printers map[string]printer.Printer, err error) {
	return as.printerStore.List()
}

func (as AdminService) TruncateModels() error {
	err := as.modelStore.Truncate()
	if err != nil {
		return err
	}
	log.Info("truncating models")
	return nil
}

func (as AdminService) TruncatePrinters() error {
	err := as.printerStore.Truncate()
	if err != nil {
		return err
	}
	log.Info("truncating printers")
	return nil
}

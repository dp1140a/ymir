package admin

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
)

type AdminHandler struct {
	api.Handler
}

func NewAdminHandler() api.HandlerIFace {
	ah := AdminHandler{
		Handler: api.Handler{
			Prefix:  "/admin",
			Service: NewAdminService(),
		},
	}

	ah.Routes = ah.addRoutes()
	return ah
}

func (ah AdminHandler) addRoutes() []api.Route {
	return []api.Route{
		{
			"listModelsAdmin",
			http.MethodGet,
			"/models",
			false,
			ah.listModels,
		},
		{
			"truncateModels",
			http.MethodDelete,
			"/models",
			false,
			ah.truncateModels,
		},
		{
			"listPrintersAdmin",
			http.MethodGet,
			"/printers",
			false,
			ah.listPrinters,
		},
		{
			"truncatePrinters",
			http.MethodDelete,
			"/printers",
			false,
			ah.truncatePrinters,
		},
	}
}

func (ah AdminHandler) GetRoutes() []api.Route {
	return ah.Routes
}

func (ah AdminHandler) GetService() api.Service {
	return ah.Service
}

func (ah AdminHandler) GetPrefix() string {
	return ah.Prefix
}

/*
GET /models (200, 500) -- get all models
*/
func (ah AdminHandler) listModels(w http.ResponseWriter, r *http.Request) {
	models, err := ah.Service.(AdminServiceIface).ListModels()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(models)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Errorf("http write error: %v", err)
	}
}

/*
GET /model (200, 500) -- get all models
*/
func (ah AdminHandler) listPrinters(w http.ResponseWriter, r *http.Request) {
	printers, err := ah.Service.(AdminServiceIface).ListPrinters()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(printers)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Errorf("http write error: %v", err)
	}
}

func (ah AdminHandler) truncateModels(w http.ResponseWriter, r *http.Request) {
	err := ah.Service.(AdminServiceIface).TruncateModels()
	log.Info("truncating models")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(`{"truncate": "ok"}`)
}

func (ah AdminHandler) truncatePrinters(w http.ResponseWriter, r *http.Request) {
	err := ah.Service.(AdminServiceIface).TruncatePrinters()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(`{"truncate": "ok"}`)
}

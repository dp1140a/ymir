package printer

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
)

type PrinterHandler struct {
	api.Handler
}

func NewPrinterHandler() api.HandlerIFace {
	ph := PrinterHandler{
		Handler: api.Handler{
			Prefix:  "/printer",
			Service: NewPrinterService(),
		},
	}

	ph.Routes = []api.Route{
		{
			"corsPreflight",
			http.MethodOptions,
			"/",
			false,
			ph.corsPreflightHandler,
		},
		{
			"createPrinter",
			http.MethodPost,
			"/",
			false,
			ph.create,
		},
		{
			"updatePrinter",
			http.MethodPut,
			"/{id}",
			false,
			ph.update,
		},
		{
			"deletePrinter",
			http.MethodDelete,
			"/{id}",
			false,
			ph.delete,
		},
		{
			"listAllPrinters",
			http.MethodGet,
			"/",
			false,
			ph.listAll,
		},
		{
			"inspectPrinter",
			http.MethodGet,
			"/{id}",
			false,
			ph.inspect,
		},
	}

	return ph
}

func (ph PrinterHandler) GetRoutes() []api.Route {
	return ph.Routes
}

func (ph PrinterHandler) GetService() api.Service {
	return ph.Service
}

func (ph PrinterHandler) GetPrefix() string {
	return ph.Prefix
}

/*
POST /Printer [Printer{}] (201, 400, 500) -- adds a Printer
*/
func (ph PrinterHandler) create(w http.ResponseWriter, r *http.Request) {}

/*
PUT /Printer/{id} [Printer{}] (201, 400, 404, 409, 500) -- updates the Printer with {id} as [Printer{}]
*/
func (ph PrinterHandler) update(w http.ResponseWriter, r *http.Request) {}

/*
DELETE /Printer/{id} (200, 404, 500) -- deletes the Printer with {id}
*/
func (ph PrinterHandler) delete(w http.ResponseWriter, r *http.Request) {}

/*
GET /Printer (200, 500) -- get all Printers
*/
func (ph PrinterHandler) listAll(w http.ResponseWriter, r *http.Request) {
	printers, err := ph.Service.(PrintersService).ListPrinters()
	if err != nil {
		log.Errorf("list all models service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(printers)
	if err != nil {
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
GET /Printer/{id} (200, 401, 404, 500) -- gets Printer with {id}
*/
func (ph PrinterHandler) inspect(w http.ResponseWriter, r *http.Request) {}

func (ph PrinterHandler) corsPreflightHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("CORS Request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, X-Sveltekit-Action")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
	w.WriteHeader(http.StatusOK)
}

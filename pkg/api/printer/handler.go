package printer

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
)

const (
	_PRINTER_NAME  = "displayName"
	_URL           = "url"
	_API_KEY       = "apiKey"
	_API_TYPE      = "apiType"
	_LOCATION      = "location"
	_PRINTER_MAKE  = "printerMake"
	_PRINTER_MODEL = "printerModel"
	_TAGS          = "tags"
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
func (ph PrinterHandler) create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB is the maximum file size
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	printer := Printer{
		DisplayName: "",
		URL:         "",
		APIType:     "octoprint",
		APIKey:      "",
		Location: Location{
			Name: "",
		},
		Type: PrinterType{
			Make:  "",
			Model: "",
		},
		DateAdded: time.Now(),
		Tags:      []string{},
	}

	for k, v := range r.MultipartForm.Value {
		switch k {
		case _PRINTER_NAME:
			printer.DisplayName = v[0]
		case _URL:
			printer.URL = v[0]
		case _API_KEY:
			printer.APIKey = v[0]
		case _LOCATION:
			printer.Location.Name = v[0]
		case _PRINTER_MAKE:
			printer.Type.Make = v[0]
		case _PRINTER_MODEL:
			printer.Type.Model = v[0]
		case _TAGS:
			printer.Tags = v
		}
	}

	err = ph.Service.(PrintersService).CreatePrinter(printer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("{'status': 'ok'}")
}

/*
PUT /Printer/{id} [Printer{}] (201, 400, 404, 409, 500) -- updates the Printer with {id} as [Printer{}]
*/
func (ph PrinterHandler) update(w http.ResponseWriter, r *http.Request) {
	var printer = &Printer{}
	err := json.NewDecoder(r.Body).Decode(&printer)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

/*
DELETE /Printer/{id} (200, 404, 500) -- deletes the Printer with {id}
*/
func (ph PrinterHandler) delete(w http.ResponseWriter, r *http.Request) {
	printerId := chi.URLParam(r, "id")
	rev := chi.URLParam(r, "rev")
	err := ph.Service.(PrintersService).DeletePrinter(printerId, rev)
	if err != nil {
		log.Errorf("delet model handler error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

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

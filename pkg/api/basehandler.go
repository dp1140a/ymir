package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"ymir/pkg/version"

	"net/http"
)

type BaseHandler struct {
	Handler
	mux *chi.Mux
}

func NewBaseHandler(logger *log.Logger, m *chi.Mux) HandlerIFace {
	bh := BaseHandler{
		mux: m,
		Handler: Handler{
			Prefix:  "/",
			Service: nil,
		},
	}

	bh.Routes = []Route{
		Route{
			"ping",
			http.MethodGet,
			"/ping",
			false,
			bh.pingHandler,
		},
		{
			"api",
			http.MethodGet,
			"/api",
			false,
			bh.getAPI,
		},
		{
			"version",
			http.MethodGet,
			"/version",
			false,
			bh.getVersion,
		},
		{
			"metrics",
			http.MethodGet,
			"/metrics",
			false,
			promhttp.Handler().ServeHTTP,
		},
	}
	return bh
}

func (bh BaseHandler) GetRoutes() []Route {
	return bh.Routes
}

func (bh BaseHandler) GetService() Service {
	return bh.Service
}

func (bh BaseHandler) GetPrefix() string {
	return bh.Prefix
}

/*
*
GET /  (200) -- Returns "OK"
*/
func RootHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I am Root\n"))
}

/*
*
GET /ping  (200) -- Returns "OK"
*/
func (bh *BaseHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

/*
*
GET /api  (200) -- Returns JSON of API
*/
func (bh *BaseHandler) getAPI(w http.ResponseWriter, _ *http.Request) {
	var routes []string
	chi.Walk(bh.mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routes = append(routes, fmt.Sprintf("%s %s", method, route))
		return nil
	})

	routeMap := make(map[string][]string)
	for _, v := range routes {
		// This will split "{METHOD} {ROUTE}" into ["{METHOD}", "{ROUTE}"]
		route := strings.Split(v, " ")
		if _, exists := routeMap[route[1]]; exists {
			//if route exists append method to route entry in the map
			routeMap[route[1]] = append(routeMap[route[1]], route[0])
		} else {
			// if route does not exist add new route to map with method
			routeMap[route[1]] = []string{route[0]}
		}
	}

	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(routeMap)
}

/*
*
GET /version  (200) -- Returns JSON of the current version
*/
func (bh *BaseHandler) getVersion(w http.ResponseWriter, _ *http.Request) {
	versionInfo := version.NewVersionInfo()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(versionInfo)
}

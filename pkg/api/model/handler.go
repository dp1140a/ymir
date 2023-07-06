package model

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
)

type ModelHandler struct {
	api.Handler
}

/*
*
200 OK – Request completed successfully
201 Created – Document created and stored on disk
202 Accepted – Document data accepted, but not yet stored on disk
304 Not Modified – Document wasn’t modified since specified revision
400 Bad Request – Invalid request body or parameters
401 Unauthorized – Write privileges required
404 Not Found – Specified database or document ID doesn’t exists
409 Conflict – Document with the specified ID already exists or specified revision is not latest for target document
*/
func NewModelHandler() api.HandlerIFace {
	mh := ModelHandler{
		Handler: api.Handler{
			Prefix:  "/model",
			Service: NewModelService(),
		},
	}

	mh.Routes = []api.Route{
		{
			"createModel",
			http.MethodPost,
			"/",
			false,
			mh.create,
		},
		{
			"updateModel",
			http.MethodPut,
			"/{id}",
			false,
			mh.update,
		},
		{
			"deleteModel",
			http.MethodDelete,
			"/{id}",
			false,
			mh.delete,
		},
		{
			"listAllModels",
			http.MethodGet,
			"/",
			false,
			mh.listAll,
		},
		{
			"listModelsByTag",
			http.MethodGet,
			"/{tag}",
			false,
			mh.listByTag,
		},
		{
			"inspectModel",
			http.MethodGet,
			"/{id}",
			false,
			mh.inspect,
		},
	}

	return mh
}

func (mh ModelHandler) GetRoutes() []api.Route {
	return mh.Routes
}

func (mh ModelHandler) GetService() api.Service {
	return mh.Service
}

func (mh ModelHandler) GetPrefix() string {
	return mh.Prefix
}

/*
POST /model [Model{}] (201, 400, 500) -- adds a model
*/
func (mh ModelHandler) create(w http.ResponseWriter, r *http.Request) {

}

/*
PUT /model/{id} [Model{}] (201, 400, 404, 409, 500) -- updates the model with {id} as [Model{}]
*/
func (mh ModelHandler) update(w http.ResponseWriter, r *http.Request) {

}

/*
DELETE /model/{id} (200, 404, 500) -- deletes the model with {id}
*/
func (mh ModelHandler) delete(w http.ResponseWriter, r *http.Request) {

}

/*
GET /model (200, 500) -- get all models
*/
func (mh ModelHandler) listAll(w http.ResponseWriter, r *http.Request) {
	models, err := mh.Service.(ModelService).ListModels()
	if err != nil {
		log.Errorf("list all models service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(models)
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
GET /model/{tag} (200, 401, 500) -- get []model{} that has {tag}
*/
func (mh ModelHandler) listByTag(w http.ResponseWriter, r *http.Request) {

}

/*
GET /model/{id} (200, 401, 404, 500) -- gets model with {id}
*/
func (mh ModelHandler) inspect(w http.ResponseWriter, r *http.Request) {

}

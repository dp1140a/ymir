package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"unsafe"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
)

type ModelHandler struct {
	api.Handler
}

const (
	_MODEL_NAME  = "modelName"
	_DESCRIPTION = "description"
	_SUMMARY     = "summary"
	_TAGS        = "tags"

	_IMAGE_FILES = "Image_Files"
	_MODEL_FILES = "Model_Files"
	_OTHER_FILES = "Other_Files"
	_PRINT_FILES = "Print_Files"
)

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
			"corsPreflight",
			http.MethodOptions,
			"/",
			false,
			mh.corsPreflightHandler,
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
		{
			"uploadFile",
			http.MethodPost,
			"/file",
			false,
			mh.uploadHandler,
		},
		{
			"fetchImage",
			http.MethodGet,
			"/image",
			false,
			mh.fetchImage,
		},
		{
			"addNote",
			http.MethodPost,
			"/note",
			false,
			mh.addNote,
		},
		{
			"parseGCode",
			http.MethodGet,
			"/gcode",
			false,
			mh.parseGCode,
		},
		{
			"fetchSTL",
			http.MethodGet,
			"/stl",
			false,
			mh.fetchSTL,
		},
		{
			"fetchSTLThumbnail",
			http.MethodGet,
			"/stl/image",
			false,
			mh.fetchSTLThumbnail,
		},
		{
			"exportModel",
			http.MethodGet,
			"/export",
			false,
			mh.exportModelHandler,
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
	err := r.ParseMultipartForm(32 << 20) // 32 MB is the maximum file size
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := Model{
		DisplayName: "",
		Tags:        []Tags{},
		Images:      []FileType{},
		ModelFiles:  []FileType{},
		OtherFiles:  []FileType{},
		PrintFiles:  []FileType{},
		DateCreated: time.Now(),
		VersionLog:  []ModelVersion{},
		Description: "",
		Summary:     "",
		Notes:       []Note{},
	}

	//fmt.Printf("Form Values: %v\n", r.MultipartForm.Value)
	for k, v := range r.MultipartForm.Value {
		switch k {
		case _DESCRIPTION:
			model.Description = v[0]
		case _SUMMARY:
			model.Summary = v[0]
		case _IMAGE_FILES:
			if len(model.Images) > 0 {
				model.Images = append(model.Images, makeFileType(v)...)
			}
		case _MODEL_FILES:
			if len(model.Images) > 0 {
				log.Debugf("appending Model File %v", v)
				model.ModelFiles = append(model.ModelFiles, makeFileType(v)...)
			}
		case _MODEL_NAME:
			model.DisplayName = v[0]
		case _OTHER_FILES:
			if len(model.Images) > 0 {
				log.Debugf("appending Other File %v", v)
				model.OtherFiles = append(model.OtherFiles, makeFileType(v)...)
			}
		case _PRINT_FILES:
			if len(model.Images) > 0 {
				log.Debugf("appending Print File %v", v)
				model.PrintFiles = append(model.PrintFiles, makeFileType(v)...)
			}
		case _TAGS:
			model.Tags = *(*[]Tags)(unsafe.Pointer(&v))
		}
	}

	err = mh.Service.(ModelService).CreateModel(model)
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
PUT /model/{id} [Model{}] (201, 400, 404, 409, 500) -- updates the model with {id} as [Model{}]
*/
func (mh ModelHandler) update(w http.ResponseWriter, r *http.Request) {
	var model = &Model{}
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(model.Json())

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
	modelId := chi.URLParam(r, "id")
	log.Debugf("inspecting modelId: %v", modelId)
	model, err := mh.Service.(ModelService).GetModel(modelId)
	if err != nil {
		log.Errorf("list all models service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println(string(js))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Errorf("http write error: %v", err)
	}
}

func (mh ModelHandler) exportModelHandler(w http.ResponseWriter, r *http.Request) {
	dirPath := r.URL.Query().Get("path")
	if dirPath == "" {
		http.Error(w, "Missing 'dir' parameter", http.StatusBadRequest)
		return
	}

	zipFileName := fmt.Sprintf("%s.zip", r.URL.Query().Get("filename"))
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))

	if err := mh.Service.(ModelService).ExportModel(dirPath, w); err != nil {
		http.Error(w, fmt.Sprintf("Error zipping directory: %s", err), http.StatusInternalServerError)
		return
	}
}

func (mh ModelHandler) uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB is the maximum file size
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Printf("FILES: %v\n", r.MultipartForm.Value)
	fType := ""
	if _, ok := r.MultipartForm.Value[_MODEL_FILES]; ok {
		fType = _MODEL_FILES
	} else if _, ok = r.MultipartForm.Value[_PRINT_FILES]; ok {
		fType = _PRINT_FILES
	} else if _, ok = r.MultipartForm.Value[_IMAGE_FILES]; ok {
		fType = _IMAGE_FILES
	} else {
		fType = _OTHER_FILES
	}

	//fmt.Println(fType)
	// Get the file from the request
	file, handler, err := r.FormFile(fType)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := mh.Service.(ModelService).UploadFiles(file, handler.Filename)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
}

func (mh ModelHandler) fetchImage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	buf, err := mh.Service.(ModelService).FetchModelImage(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(buf)
}

func (mh ModelHandler) fetchSTL(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	buf, err := mh.Service.(ModelService).FetchSTL(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/sla")
	w.Write(buf)
}

func (mh ModelHandler) fetchSTLThumbnail(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	log.Info(path)
	imgStr := mh.Service.(ModelService).FetchSTLThumbnail(path)
	w.Header().Set("Content-Type", "image/png")
	w.Write([]byte(imgStr))
}

func (mh ModelHandler) addNote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := Model{
		Id:  r.FormValue("_id"),
		Rev: r.FormValue("_rev"),
		Notes: []Note{
			{
				Text: r.FormValue("text"),
				Date: time.Now(),
			},
		},
	}

	err = mh.Service.(ModelService).AddNote(model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("{'status': 'ok'}")

}

func (mh ModelHandler) parseGCode(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	gcode, err := mh.Service.(ModelService).ParseGCode(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(gcode)
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

func (mh ModelHandler) corsPreflightHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("CORS Request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, X-Sveltekit-Action")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
	w.WriteHeader(http.StatusOK)
}

func makeFileType(files []string) (fileTypes []FileType) {
	fileTypes = []FileType{}
	for _, v := range files {
		fileTypes = append(fileTypes, FileType{
			Path: v,
		})
	}

	return
}

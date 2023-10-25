package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unsafe"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
	"ymir/pkg/api/printer"
)

type ModelHandler struct {
	api.Handler
}

const (
	_MODEL_NAME      = "modelName"
	_MODEL_BASE_PATH = "basePath"
	_DESCRIPTION     = "description"
	_SUMMARY         = "summary"
	_TAGS            = "tags"

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

	mh.Routes = mh.addRoutes()
	return mh
}

func (mh ModelHandler) addRoutes() []api.Route {
	return []api.Route{
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
			"uploadFileToPrinter",
			http.MethodPost,
			"/file/printer",
			false,
			mh.UploadFileToPrinter,
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

	if log.GetLevel() == log.DebugLevel {
		fmt.Printf("Form Values: %q\n", r.MultipartForm.Value)
		fmt.Printf("%q, %v , %T\n", r.MultipartForm.Value["Other_Files"], len(r.MultipartForm.Value["Other_Files"]), r.MultipartForm.Value["Other_Files"])
	}

	for k, v := range r.MultipartForm.Value {
		switch k {
		case _DESCRIPTION:
			model.Description = v[0]
		case _SUMMARY:
			model.Summary = v[0]
		case _IMAGE_FILES:
			if len(v) > 0 && v[0] != "" {
				model.Images = append(model.Images, makeFileType(v)...)
			}
		case _MODEL_FILES:
			if len(v) > 0 && v[0] != "" {
				log.Debugf("appending Model File %v", v)
				model.ModelFiles = append(model.ModelFiles, makeFileType(v)...)
			}
		case _MODEL_NAME:
			model.DisplayName = v[0]
		case _OTHER_FILES:
			if len(v) > 0 && v[0] != "" {
				log.Debugf("appending Other File %v", v)
				model.OtherFiles = append(model.OtherFiles, makeFileType(v)...)
			}
		case _PRINT_FILES:
			if len(v) > 0 && v[0] != "" {
				log.Debugf("appending Print File %v", v)
				model.PrintFiles = append(model.PrintFiles, makeFileType(v)...)
			}
		case _TAGS:
			model.Tags = *(*[]Tags)(unsafe.Pointer(&v))
		}
	}

	err = mh.Service.(ModelServiceIface).CreateModel(model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("{'status': 'ok'}")
}

/*
PUT /model/{id} [Model{}] (200, 400, 500) -- updates the model with {id} as [Model{}]
*/
func (mh ModelHandler) update(w http.ResponseWriter, r *http.Request) {
	var model = Model{}
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if log.GetLevel() == log.DebugLevel {
		fmt.Println(model.Json())
	}
	rev, err := mh.Service.(ModelService).UpdateModel(model)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody := map[string]string{"status": "ok", "rev": rev}
	json.NewEncoder(w).Encode(respBody)
}

/*
DELETE /model/{id}?rev (204, 500) -- deletes the model with {id}
*/
func (mh ModelHandler) delete(w http.ResponseWriter, r *http.Request) {
	modelId := chi.URLParam(r, "id")
	rev := r.URL.Query().Get("rev")
	if log.GetLevel() == log.DebugLevel {
		fmt.Printf("id : %v / rev: %v\n", modelId, rev)
	}
	err := mh.Service.(ModelServiceIface).DeleteModel(modelId, rev)
	if err != nil {
		log.Errorf("delete printer handler error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

/*
GET /model (200, 500) -- get all models
*/
func (mh ModelHandler) listAll(w http.ResponseWriter, r *http.Request) {
	models, err := mh.Service.(ModelServiceIface).ListModels()
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
func (mh ModelHandler) listByTag(w http.ResponseWriter, r *http.Request) {}

/*
GET /model/{id} (200, 400, 500) -- gets model with {id}
*/
func (mh ModelHandler) inspect(w http.ResponseWriter, r *http.Request) {
	modelId := chi.URLParam(r, "id")
	if modelId == "" {
		log.Error("modelId is missing or bad")
		http.Error(w, "modelId is missing or bad", http.StatusBadRequest)
	}
	log.Debugf("inspecting modelId: %v", modelId)
	model, err := mh.Service.(ModelServiceIface).GetModel(modelId)
	if err != nil {
		log.Errorf("inspect model service error: %v", err)
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

/*
GET /model/export?path (204, 400, 500) -- gets model with {id}
*/
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
	}
	w.WriteHeader(http.StatusNoContent)
}

/*
POST /model/{id} (202, 400, 500) -- gets model with {id}
*/
func (mh ModelHandler) uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB is the maximum file size

	fmt.Println(r.MultipartForm.Value)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	// Get the file from the request
	file, handler, err := r.FormFile(fType)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/**
	If we have a basePath param this is not for a new model.
	*/
	basePath := r.URL.Query().Get(_MODEL_BASE_PATH)
	key := ""
	if basePath != "" {
		log.Infof("base path: %v", basePath)
		key, err = mh.Service.(ModelService).UploadFilesExistingModel(file, handler.Filename, basePath)

	} else { // New Model
		key, err = mh.Service.(ModelService).UploadFilesNewModel(file, handler.Filename)

	}
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(key))
}

/*
POST /model/file/printer?file:<string>&print:<bool> (201, 400, 401, 500) -- Uploads file to printer
Body: {Printer}
https://github.com/mcuadros/go-octoprint/blob/master/files.go#L106
*/
func (mh ModelHandler) UploadFileToPrinter(w http.ResponseWriter, r *http.Request) {
	//Get the file
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		return
	}

	printFile, err := strconv.ParseBool(r.URL.Query().Get("print"))
	if err != nil {
		log.Errorf("malformed print parameter.  setting to false")
		printFile = false
	}
	//Decode the printer
	var p printer.Printer
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Errorf("printer decode error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ur := UploadFileRequest{
		Location: "local",
		Select:   true,
		Print:    printFile,
	}
	file, _ := os.Open(filePath)
	defer file.Close()
	err = ur.AddFile(filepath.Base(file.Name()), file)
	if err != nil {
		log.Errorf("error adding file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ur.writer().WriteField("path", "ymir")
	ur.addSelectPrintAndClose()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/files/%s", p.URL, ur.Location), ur.b)
	if err != nil {
		log.Errorf("request error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Header.Add("Host", "localhost:5000")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", ur.w.FormDataContentType())
	req.Header.Add("X-Api-Key", p.APIKey)

	c := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Errorf("response error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "application/json")

	if resp.StatusCode == 401 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(`{"error": "Missing or invalid API key"}`)
	} else if resp.StatusCode == 201 {
		w.WriteHeader(http.StatusCreated)
		log.Debug(resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		if log.GetLevel() == log.DebugLevel {
			fmt.Println(string(body))
		}
		json.NewEncoder(w).Encode(string(body))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("{error: %v, code: %v, body: %v}", resp.Status, resp.StatusCode, resp.Body))
	}
}

type UploadFileRequest struct {
	// Location is the target location to which to upload the file. Currently
	// only `local` and `sdcard` are supported here, with local referring to
	// OctoPrint’s `uploads` folder and `sdcard` referring to the printer’s
	// SD card. If an upload targets the SD card, it will also be stored
	// locally first.
	Location string
	// Select whether to select the file directly after upload (true) or not
	// (false). Optional, defaults to false. Ignored when creating a folder.
	Select bool
	//Print whether to start printing the file directly after upload (true) or
	// not (false). If set, select is implicitely true as well. Optional,
	// defaults to false. Ignored when creating a folder.
	Print bool
	b     *bytes.Buffer
	w     *multipart.Writer
}

// AddFile adds a new file to be uploaded from a given reader.
func (req *UploadFileRequest) AddFile(filename string, r io.Reader) error {
	w, err := req.writer().CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	return err

}

func (req *UploadFileRequest) writer() *multipart.Writer {
	if req.w == nil {
		req.b = bytes.NewBuffer(nil)
		req.w = multipart.NewWriter(req.b)
	}

	return req.w
}

func (req *UploadFileRequest) addSelectPrintAndClose() error {
	err := req.writer().WriteField("select", fmt.Sprintf("%t", req.Select))
	if err != nil {
		return err
	}

	err = req.writer().WriteField("print", fmt.Sprintf("%t", req.Print))
	if err != nil {
		return err
	}

	return req.writer().Close()
}

/*
GET /image (200, 500) -- Fetches Model Image
*/
func (mh ModelHandler) fetchImage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	buf, err := mh.Service.(ModelService).FetchModelImage(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

/*
GET /stl (200, 500) -- Fetches Model STL
*/
func (mh ModelHandler) fetchSTL(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	buf, err := mh.Service.(ModelService).FetchSTL(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/sla")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

/*
GET /stl/image (200, 500) -- Fetches Model STL Thumbnail
*/
func (mh ModelHandler) fetchSTLThumbnail(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	log.Info(path)
	imgStr := mh.Service.(ModelService).FetchSTLThumbnail(path)
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(imgStr))
}

/*
POST /note (201, 400, 500) -- Add note to Model
*/
func (mh ModelHandler) addNote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB is the maximum file size
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
				Text: r.FormValue("noteText"),
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("{'status': 'ok'}")

}

/*
GET /gcode (200, 500) -- Fetches Model STL Thumbnail
*/
func (mh ModelHandler) parseGCode(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	gcode, err := mh.Service.(ModelService).GetGCodeMetaData(path)
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

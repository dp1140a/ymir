package model

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
	"ymir/pkg/db"
	"ymir/pkg/gcode"
	"ymir/pkg/stl"
)

type ModelService struct {
	name      string
	DataStore *db.DB
	config    *ModelsConfig
}

func NewModelService() (modelService api.Service) {
	ms := ModelService{
		name:   "Model",
		config: NewModelsConfig(),
	}

	ds := db.NewDB()
	ds.Connect()
	ms.DataStore = ds

	err := makeDirIfNotExists(ms.config.UploadsTempDir)
	if err != nil {
		return nil

	}
	makeDirIfNotExists(ms.config.ModelsDir)
	if err != nil {
		return nil
	}
	return ms
}

func (ms ModelService) GetName() (name string) {
	return ms.name
}

func makeDirIfNotExists(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Errorf("error creating dir %v: %v", path, err)
			return err
		}
	}
	return nil
}

func (ms ModelService) GetModel(id string) (model Model, err error) {
	row := ms.DataStore.GetDB().Get(context.TODO(), id)
	model = Model{}
	if err = row.ScanDoc(&model); err != nil {
		log.Error(err)
		return model, err
	}

	return model, nil
}

func (ms ModelService) ExportModel(path string, writer io.Writer) error {
	zipWriterObj := zip.NewWriter(writer)
	defer zipWriterObj.Close()
	modelPath := fmt.Sprintf("%s/%s", ms.config.ModelsDir, path)
	err := filepath.Walk(modelPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			log.Error(err)
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(modelPath, filePath)
		if err != nil {
			log.Error()
			return err
		}

		zipFile, err := zipWriterObj.Create(relPath)
		if err != nil {
			log.Error()
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			log.Error()
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipFile, file)
		return err
	})

	return err
}

func addFilesToZip(w *zip.Writer, basePath, baseInZip string) error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullfilepath := filepath.Join(basePath, file.Name())
		if _, err := os.Stat(fullfilepath); os.IsNotExist(err) {
			// ensure the file exists. For example a symlink pointing to a non-existing location might be listed but not actually exist
			continue
		}

		if file.Mode()&os.ModeSymlink != 0 {
			// ignore symlinks alltogether
			continue
		}

		if file.IsDir() {
			if err := addFilesToZip(w, fullfilepath, filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		} else if file.Mode().IsRegular() {
			dat, err := ioutil.ReadFile(fullfilepath)
			if err != nil {
				return err
			}
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else {
			// we ignore non-regular files because they are scary
		}
	}
	return nil
}

func (ms ModelService) CreateModel(model Model) (err error) {
	ctx := context.TODO()
	model.Id = GenId()
	//fmt.Println(model.Json())
	// We want to organize all the files and the model
	err = ms.organize(&model)
	if err != nil {
		log.Error(err)
		return err
	}

	docId, rev, err := ms.DataStore.GetDB().CreateDoc(ctx, model)
	if err != nil {
		log.Error(err)
		return err
	} else {
		log.Infof("created model %v in db with docid %v and rev %v", model.Id, docId, rev)
		return nil
	}
}

func (ms ModelService) UpdateModel(model Model) (rev string, err error) {
	ctx := context.TODO()

	rev, err = ms.DataStore.GetDB().Put(ctx, model.Id, model)
	if err != nil {
		log.Error(err)
		return "", err
	} else {
		log.Infof("updated model %v in db with rev %v", model.Id, rev)
		return rev, nil
	}
}

func (ms ModelService) DeleteModel(id string, rev string) error {
	newRev, err := ms.DataStore.GetDB().Delete(context.TODO(), id, rev)
	if err != nil {
		return err
	}
	log.Infof("deleted model %v in db with newRev %v", id, newRev)
	return nil
}

func (ms ModelService) ListModels() ([]Model, error) {
	query := `{
		"selector": {
			"displayName": {"$regex": ".+"}
		},
		"fields": ["_id", "_rev", "basePath", "description", "displayName", "images", "summary", "tags" ]
	}`
	var q interface{}
	_ = json.Unmarshal([]byte(query), q)
	rows, err := ms.DataStore.GetDB().Find(context.TODO(), query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	docs := []Model{}
	for rows.Next() {
		doc := &Model{}
		err = rows.ScanDoc(doc)
		if err != nil {
			log.Errorf("rows error: %v\n", err)
			return nil, err
		}
		docs = append(docs, *doc)
	}

	return docs, nil
}

func (ms ModelService) GetModelsByTag(tag string) ([]Model, error) {

	query := `{
	   "selector": {
		  "tags": {
			 "$elemMatch": {
				"$eq": "---"
			 }
		  }
	   }
	}`

	query = strings.Replace(query, "---", tag, 1)
	//fmt.Println(query)
	var q interface{}
	_ = json.Unmarshal([]byte(query), q)
	rows, err := ms.DataStore.GetDB().Find(context.TODO(), query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	docs := []Model{}
	for rows.Next() {
		doc := &Model{}
		err = rows.ScanDoc(doc)
		if err != nil {
			log.Errorf("rows error: %v\n", err)
			return nil, err
		}
		docs = append(docs, *doc)
	}

	return docs, nil
}

func writeFile(file *multipart.File, path string) error {
	// Create a new file in the uploads directory
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
		return err
	}
	defer f.Close()

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, *file)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (ms ModelService) UploadFilesExistingModel(file multipart.File, filename string, basePath string) (key string, err error) {
	defer file.Close()

	//Look a folder named "files" in the model folder
	entries, err := os.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}

	hasFilesDir := false
	for _, dir := range entries {
		if dir.Name() == "files" {
			hasFilesDir = true
			break
		}
	}

	var path string
	if hasFilesDir {
		//If has files dir then move the new file there
		key = filepath.Join("files", filename)
		path = filepath.Join(basePath, key)
	} else {
		//else move to basePath
		key = filename
		path = filepath.Join(basePath, key)
	}
	log.Debugf("upload filepath: %v", path)
	err = writeFile(&file, path)
	if err != nil {
		log.Errorf("error writing file: %v", err)
		return key, err
	}
	return key, nil
}

func (ms ModelService) UploadFilesNewModel(file multipart.File, filename string) (key string, err error) {
	defer file.Close()

	//Generate key
	tK := GenId()
	err = makeDirIfNotExists(filepath.Join(ms.config.UploadsTempDir, tK))
	if err != nil {
		log.Debugf("error making upload dir: %v", err)
		return "", err
	}
	key = filepath.Join(tK, filename)
	path := filepath.Join(ms.config.UploadsTempDir, key)
	log.Debugf("upload filepath: %v", path)

	err = writeFile(&file, path)
	if err != nil {
		log.Errorf("error writing file: %v", err)
		return "", err
	}
	return path, nil
}

func (ms ModelService) FetchModelImage(imagePath string) (imageBytes []byte, err error) {
	imageBytes, err = os.ReadFile(imagePath)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return imageBytes, nil
}

func (ms ModelService) FetchSTL(filepath string) (stlBytes []byte, err error) {
	stlBytes, err = os.ReadFile(filepath)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return stlBytes, nil
}

func (ms ModelService) FetchSTLThumbnail(filepath string) string {
	img := stl.Image(filepath)
	imgStr := stl.ThumbnailBase64(img, 128, 128)
	//fmt.Println(imgStr)
	return imgStr
}

func (ms ModelService) AddNote(model Model) (err error) {
	fmt.Println(model.Json())

	ctx := context.TODO()

	existingModel, err := ms.GetModel(model.Id)
	if err != nil {
		log.Error(err)
		return err
	}

	existingModel.Notes = append(existingModel.Notes, model.Notes...)

	rev, err := ms.DataStore.GetDB().Put(ctx, existingModel.Id, existingModel)
	if err != nil {
		log.Error(err)
		return err
	} else {
		log.Infof("updated model %v in db with rev %v", model.Id, rev)
		return nil
	}
}

func (ms ModelService) ParseGCode(path string) (gcode.GCodeMetaData, error) {
	g := gcode.NewGCode(path)
	err := g.ParseGCode(false)
	if err != nil {
		log.Error(err)
		return gcode.GCodeMetaData{}, err
	}

	return g.MetaData, nil
}

func (ms ModelService) organize(model *Model) (err error) {
	mDir := filepath.Join(ms.config.ModelsDir, model.Id)
	log.Debugf("model dir: %v", mDir)
	err = makeDirIfNotExists(mDir)
	if err != nil {
		return err
	}
	model.BasePath = mDir

	moveAndClean := func(itemPath string) error {
		from := itemPath
		fDir, fName := filepath.Split(itemPath)
		to := filepath.Join(mDir, fName)
		log.Debugf("moving model from: %v to: %v ", from, to)

		if err := os.Rename(from, to); err != nil {
			log.Error(err)
			return err
		}

		if err := os.RemoveAll(filepath.Join(ms.config.UploadsTempDir, fDir)); err != nil {
			log.Error(err)
			return err
		}

		return nil
	}

	updatePaths := func(files []FileType) error {
		log.Debugf("num files: %v", len(files))
		for i, file := range files {
			err := moveAndClean(file.Path)
			if err != nil {
				return err
			}
			files[i].Path = filepath.Base(file.Path)
		}
		return nil
	}

	log.Debugf("ModelFiles:")
	if err := updatePaths(model.ModelFiles); err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("OtherFiles:")
	if err := updatePaths(model.OtherFiles); err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("PrintFiles:")
	if err := updatePaths(model.PrintFiles); err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("Images:")
	if err := updatePaths(model.Images); err != nil {
		log.Error(err)
		return err
	}

	if err := model.WriteModel(mDir); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

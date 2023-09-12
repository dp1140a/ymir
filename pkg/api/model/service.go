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

	makeDirIfNotExists(ms.config.UploadsTempDir)
	makeDirIfNotExists(ms.config.ModelsDir)

	return ms
}

func (ms ModelService) GetName() (name string) {
	return ms.name
}

func makeDirIfNotExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Errorf("error creating dir %v: %v", path, err)
		}
	}
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

func (ms ModelService) UpdateModel(model Model) (err error) {
	ctx := context.TODO()

	rev, err := ms.DataStore.GetDB().Put(ctx, model.Id, model)
	if err != nil {
		return err
	} else {
		log.Infof("updated model %v in db with rev %v", model.Id, rev)
		return nil
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

func (ms ModelService) UploadFiles(file multipart.File, filename string) (key string, err error) {

	defer file.Close()

	//Generate key
	tK := GenId()
	makeDirIfNotExists(fmt.Sprintf("%s/%s", ms.config.UploadsTempDir, tK))
	key = fmt.Sprintf("%s/%s", tK, filename)
	path := fmt.Sprintf("%s/%s", ms.config.UploadsTempDir, key)
	log.Debugf("upload filepath: %v", path)
	// Create a new file in the uploads directory
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer f.Close()

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return key, nil
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

func (ms ModelService) organize(model *Model) error {
	mDir := filepath.Join(ms.config.ModelsDir, model.Id)
	makeDirIfNotExists(mDir)

	moveAndClean := func(itemPath string) error {
		from := itemPath
		fDir, fName := filepath.Split(itemPath)
		to := filepath.Join(mDir, fName)
		log.Debugf("moving model from:\n\t%v\nto:\n\t%v\n", from, to)

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
		log.Debugf("files: %v\n", len(files))
		for i, file := range files {
			err := moveAndClean(file.Path)
			if err != nil {
				return err
			}
			files[i].Path = filepath.Base(file.Path)
		}
		return nil
	}

	if err := updatePaths(model.ModelFiles); err != nil {
		log.Error(err)
		return err
	}

	if err := updatePaths(model.OtherFiles); err != nil {
		log.Error(err)
		return err
	}

	if err := updatePaths(model.PrintFiles); err != nil {
		log.Error(err)
		return err
	}

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
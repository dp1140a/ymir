package model

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"ymir/pkg/api"
	"ymir/pkg/api/model/store"
	"ymir/pkg/api/model/types"
	"ymir/pkg/gcode"
	"ymir/pkg/stl"
	"ymir/pkg/utils"
)

type ModelServiceIface interface {
	api.Service
	CreateModel(model types.Model) (id string, err error)
	UpdateModel(model types.Model) (err error)
	DeleteModel(id string) error
	GetModel(id string) (types.Model, error)
	ListModels() (models map[string]types.Model, err error)
	ExportModel(path string, writer io.Writer) error
	UploadFile(file multipart.File, filename string, basePath string, isExistingModel bool) (key string, err error)
	FetchModelImage(imagePath string) (imageBytes []byte, err error)
	FetchSTL(filepath string) (stlBytes []byte, err error)
	FetchSTLThumbnail(filepath string) string
	AddNote(model types.Model) error
	GetGCodeMetaData(path string) (gcode.GCodeMetaData, error)
}

type ModelService struct {
	ModelServiceIface
	name       string
	modelStore store.ModelStoreIFace
	config     *ModelsConfig
}

func NewModelService() (modelService api.Service) {
	ms := ModelService{
		name:       "Model",
		config:     NewModelsConfig(),
		modelStore: store.NewModelDataStore(),
	}

	err := utils.MakeDirIfNotExists(ms.config.UploadsTempDir)
	if err != nil {
		return nil

	}
	err = utils.MakeDirIfNotExists(ms.config.ModelsDir)
	if err != nil {
		return nil
	}
	return ms
}

func (ms ModelService) GetName() (name string) {
	return ms.name
}

func (ms ModelService) GetModel(id string) (model types.Model, err error) {
	model, err = ms.modelStore.Inspect(id)
	if err != nil {
		log.Errorf("error retrieving model: %v", err)
		return
	} else if model.Id == "" {
		err = errors.New(fmt.Sprintf("model with id: %v does not exist", id))
	}
	return
}

func (ms ModelService) ExportModel(path string, writer io.Writer) error {
	zipWriterObj := zip.NewWriter(writer)
	defer func(zipWriterObj *zip.Writer) {
		err := zipWriterObj.Close()
		if err != nil {
			log.Error(err)
		}
	}(zipWriterObj)
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
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Error(err)
			}
		}(file)

		_, err = io.Copy(zipFile, file)
		return err
	})

	return err
}

func addFilesToZip(w *zip.Writer, basePath, baseInZip string) error {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	/*
			entries, err := os.ReadDir(dirname)
			if err != nil { ... }
			infos := make([]fs.FileInfo, 0, len(entries))
			for _, entry := range entries {
				info, err := entry.Info()
				if err != nil { ... }
				infos = append(infos, info)
		}
	*/

	for _, entry := range files {
		file, err := entry.Info()
		if err != nil {
			log.Error(err)
		}
		fullFilePath := filepath.Join(basePath, file.Name())
		if _, err := os.Stat(fullFilePath); os.IsNotExist(err) {
			// ensure the file exists. For example a symlink pointing to a non-existing location might be listed but not actually exist
			continue
		}

		if file.Mode()&os.ModeSymlink != 0 {
			// ignore symlinks alltogether
			continue
		}

		if file.IsDir() {
			if err := addFilesToZip(w, fullFilePath, filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		} else if file.Mode().IsRegular() {
			dat, err := os.ReadFile(fullFilePath)
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

func (ms ModelService) CreateModel(model types.Model) (id string, err error) {
	model.Id = utils.GenId()
	//fmt.Println(model.Json())
	// We want to organize all the files and the model
	err = ms.organize(&model)
	if err != nil {
		log.Error(err)
		return
	}

	err = ms.modelStore.Create(model)
	if err != nil {
		log.Error(err)
		return
	} else {
		log.Infof("created model %v in db", model.Id)
		return model.Id, nil
	}
}

func (ms ModelService) UpdateModel(model types.Model) (err error) {
	err = ms.modelStore.Update(model)
	if err != nil {
		log.Error(err)
		return
	} else {
		log.Infof("updated model %v in db", model.Id)
		return
	}
}

func (ms ModelService) DeleteModel(id string) (err error) {
	err = ms.modelStore.Delete(id)
	if err != nil {
		return err
	}
	log.Infof("deleted model %v in db", id)
	return nil
}

func (ms ModelService) ListModels() (models map[string]types.Model, err error) {
	models, err = ms.modelStore.List()
	if err != nil {
		log.Error(err)
		return
	}
	return
}

/*
*
@TODO Move to utils package
*/
func writeFile(file multipart.File, path string) error {
	// Create a new file in the uploads directory
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Error(err)
		}
	}(f)

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// uploadFileHelper encapsulates the common operations for uploading a file
func (ms ModelService) uploadFileHelper(file multipart.File, path string) (string, error) {
	defer file.Close()
	log.Debugf("upload filepath: %v", path)
	err := writeFile(file, path)
	if err != nil {
		log.Errorf("error writing file: %v", err)
	}

	// returning path and error to be used in the calling function for error handling
	return path, err
}

func (ms ModelService) UploadFilesExistingModel(file multipart.File, filename string, basePath string) (string, error) {
	// Look for a "files" folder in the model folder
	entries, err := os.ReadDir(basePath)
	if err != nil {
		log.Error(err)
		return "", err
	}

	hasFilesDir := false
	for _, dir := range entries {
		if dir.Name() == "files" {
			hasFilesDir = true
			break
		}
	}

	key := filename
	if hasFilesDir {
		// If "files" dir exists, move the new file there
		key = filepath.Join("files", filename)
	}
	// Call the refactored function
	path, err := ms.uploadFileHelper(file, filepath.Join(basePath, key))
	if err != nil {
		return key, err
	}
	return path, nil
}

func (ms ModelService) UploadFilesNewModel(file multipart.File, filename string) (string, error) {
	// Generate key
	tK := utils.GenId()
	err := utils.MakeDirIfNotExists(filepath.Join(ms.config.UploadsTempDir, tK))
	if err != nil {
		log.Debugf("error making upload dir: %v", err)
		return "", err
	}

	key := filepath.Join(tK, filename)
	// Call the refactored function
	path, err := ms.uploadFileHelper(file, filepath.Join(ms.config.UploadsTempDir, key))
	if err != nil {
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

func (ms ModelService) FetchSTLThumbnail(filepath string) (string, error) {
	img := stl.Image(filepath)
	imgStr, err := stl.ThumbnailBase64(img, 128, 128)
	if err != nil {
		log.Error(err)
		return "", err
	}
	//fmt.Println(imgStr)
	return imgStr, nil
}

func (ms ModelService) AddNote(model types.Model) (err error) {
	//fmt.Println(model.Json())
	existingModel, err := ms.GetModel(model.Id)
	if err != nil {
		log.Error(err)
		return err
	}
	existingModel.Notes = append(existingModel.Notes, model.Notes...)
	err = ms.modelStore.Update(existingModel)
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("updated model %v in db", model.Id)
	}
	return
}

func (ms ModelService) GetGCodeMetaData(path string) (gcode.GCodeMetaData, error) {
	g := gcode.NewGCode(path)
	err := g.ParseGCode(false)
	if err != nil {
		log.Error(err)
		return gcode.GCodeMetaData{}, err
	}

	return g.MetaData, nil
}

func (ms ModelService) organize(model *types.Model) (err error) {
	mDir := filepath.Join(ms.config.ModelsDir, model.Id)
	log.Debugf("model dir: %v", mDir)
	err = utils.MakeDirIfNotExists(mDir)
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

	updatePaths := func(files []types.FileType) error {
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

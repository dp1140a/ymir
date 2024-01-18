package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"ymir/pkg/api/model/store"
	"ymir/pkg/api/model/types"
	"ymir/pkg/utils"
)

type Importer struct {
	config     *ImporterConfig
	modelStore store.ModelStore
	baseDir    string
	Models     []types.Model
	Tags       []types.Tags
}

func NewImporter(base string, inDB bool) *Importer {
	importer := &Importer{
		config:  NewImporterConfig(),
		baseDir: base,
	}
	if inDB {
		fmt.Printf("DB Flag is set will write to DB")
		importer.modelStore = store.NewModelDataStore().(store.ModelStore)
	}

	return importer
}

func (i *Importer) PutInDB() error {
	for _, model := range i.Models {
		err := i.modelStore.Create(model)
		if err != nil {
			log.Errorf("error writing model %v %v", model.BasePath, err)
			return err
		}
		fmt.Printf("created model %v in db with docid %v\n", model.BasePath, model.Id)
	}
	return nil
}

var depth = 0

func (i *Importer) walk(path string, m *types.Model) error {
	fmt.Printf("Scanning: %v [%v]\n", path, depth)
	dirs, _ := os.ReadDir(path)
	fmt.Printf("  %v entries\n", len(dirs))
	for _, f := range dirs {
		relPath, _ := filepath.Rel(i.baseDir, path)
		fName := fmt.Sprintf("%v/%v", relPath, f.Name())
		if f.IsDir() {
			if onlyDirectories(path) {
				fmt.Printf("   %v is only directories moving on\n", path)
				i.walk(filepath.Join(path, f.Name()), m)
			} else {
				depth++
				i.walk(filepath.Join(path, f.Name()), m)
				depth--
			}

		} else {
			if isExecutable(f.Type()) || !isModelFile(f.Name()) {
				continue
			}
			if m == nil {
				fmt.Println("Creating Model:")
				m = &types.Model{
					Id:          utils.GenId(),
					DisplayName: cleanDisplayName(filepath.Base(i.baseDir)),
					BasePath:    path,
					Tags:        i.Tags,
					ModelFiles:  []types.FileType{},
					PrintFiles:  []types.FileType{},
					OtherFiles:  []types.FileType{},
					Images:      []types.FileType{},
					DateCreated: time.Now(),
					Notes:       []types.Note{},
				}
			}

			/**
			@TODO If the files is README or README.md then set the contents as the file description
			*/
			if strings.Contains(strings.ToLower(f.Name()), strings.ToLower("README")) {
				fmt.Printf("  Found a README File for model %v. Setting contents as model description\n", fName)
			}

			if slices.Contains(MODEL_TYPES, filepath.Ext(f.Name())[1:]) {
				fmt.Printf("  Found a model: %v\n", f.Name())
				m.ModelFiles = append(m.ModelFiles, types.FileType{Path: fName})
				continue
			}
			if slices.Contains(PRINT_TYPES, filepath.Ext(f.Name())[1:]) {
				fmt.Printf("  Adding Print File: %v\n", f.Name())
				m.PrintFiles = append(m.PrintFiles, types.FileType{Path: fName})
				continue
			}
			if slices.Contains(IMAGE_TYPES, filepath.Ext(f.Name())[1:]) {
				fmt.Printf("  Adding Image: %v\n", f.Name())
				m.Images = append(m.Images, types.FileType{Path: fName})
				continue
			}
			if slices.Contains(OTHER_TYPES, filepath.Ext(f.Name())[1:]) {
				if f.Name() == "model.json" {
					continue
				}
				fmt.Printf("  Adding Other File: %v\n", f.Name())
				m.OtherFiles = append(m.OtherFiles, types.FileType{Path: f.Name()})
				continue
			}
		}
	}

	if depth == 0 && m != nil {
		err := writeModel(path, m)
		if err != nil {
			return err
		}
		i.Models = append(i.Models, *m)
	}
	return nil
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9.\- ]+`)

func cleanDisplayName(str string) string {
	fmt.Printf("MODEL NAME: %v\n", nonAlphanumericRegex.ReplaceAllString(str, " "))
	return nonAlphanumericRegex.ReplaceAllString(str, " ")
}
func writeModel(path string, model *types.Model) error {
	fmt.Printf("Writng model %v to %v\n", model.DisplayName, path)
	err := model.WriteModel(path)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func onlyDirectories(dirPath string) bool {
	dirContents, _ := os.ReadDir(dirPath)
	for _, entry := range dirContents {
		if entry.IsDir() {
			continue // It's a directory, so continue checking the next entry.
		} else {
			return false // Found a non-directory, return false.
		}
	}
	return true
}

func isExecutable(mode os.FileMode) bool {
	return mode&0111 != 0
}

func isModelFile(name string) bool {
	fileTypes := [][]string{MODEL_TYPES, PRINT_TYPES, IMAGE_TYPES, OTHER_TYPES}

	for _, fileType := range fileTypes {
		if slices.Contains(fileType, filepath.Ext(name)[1:]) {
			return true
		}
	}

	return false
}

/*
This function is just a conveniience wrapper
*/
func (i *Importer) FindModels() error {
	base, err := filepath.Abs(i.baseDir)
	if err != nil {
		log.Fatal(err)
	}
	i.walk(base, nil)
	if len(i.Models) == 0 {
		fmt.Println("NO MODELS FOUND.  TRY ANOTHER DIRECTORY!")
	} else {
		fmt.Printf("%v MODELS FOUND.", len(i.Models))
	}
	return nil
}

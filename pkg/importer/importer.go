package importer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"ymir/pkg/api/model"
	"ymir/pkg/db"
)

type Importer struct {
	config  *ImporterConfig
	baseDir string
	Models  []model.Model
	Tags    []model.Tags
}

func NewImporter(base string) *Importer {
	importer := &Importer{
		config:  NewImporterConfig(),
		baseDir: base,
	}
	return importer
}

func (i *Importer) PutInDB() error {
	//if put in DB put them all in the DB
	ctx := context.TODO()
	ds := db.NewDB()
	ds.Connect()

	for _, model := range i.Models {
		docId, rev, err := ds.GetDB().CreateDoc(ctx, model)
		if err != nil {
			log.Errorf("error writing model %v %v", model.BasePath, err)
			return err
		}
		fmt.Printf("created model %v in db with docid %v and rev %v\n", model.BasePath, docId, rev)
	}
	return nil
}

var depth = 0

func (i *Importer) walk(path string, m *model.Model) error {
	fmt.Printf("Scanning: %v [%v]\n", path, depth)
	dirs, _ := os.ReadDir(path)
	fmt.Printf("  %v entries\n", len(dirs))
	for _, f := range dirs {
		fName := fmt.Sprintf("%v/%v", filepath.Base(path), f.Name())
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
			if m == nil {
				m = &model.Model{
					Id:          model.GenId(),
					DisplayName: cleanDisplayName(filepath.Base(filepath.Dir(fName))),
					BasePath:    path,
					Tags:        i.Tags,
					ModelFiles:  []model.FileType{},
					PrintFiles:  []model.FileType{},
					OtherFiles:  []model.FileType{},
					Images:      []model.FileType{},
					DateCreated: time.Now(),
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
				m.ModelFiles = append(m.ModelFiles, model.FileType{Path: fName})
				continue
			}
			if slices.Contains(PRINT_TYPES, filepath.Ext(f.Name())[1:]) {
				m.PrintFiles = append(m.PrintFiles, model.FileType{Path: fName})
				continue
			}
			if slices.Contains(IMAGE_TYPES, filepath.Ext(f.Name())[1:]) {
				m.Images = append(m.Images, model.FileType{Path: fName})
				continue
			}
			if slices.Contains(OTHER_TYPES, filepath.Ext(f.Name())[1:]) {
				if f.Name() == "model.json" {
					continue
				}
				m.OtherFiles = append(m.OtherFiles, model.FileType{Path: f.Name()})
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
	return nonAlphanumericRegex.ReplaceAllString(str, " ")
}
func writeModel(path string, model *model.Model) error {
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

/*
This function is just a conveniience wrapper
*/
func (i *Importer) FindModels() error {
	base, err := filepath.Abs(i.baseDir)
	if err != nil {
		log.Fatal(err)
	}
	i.walk(base, nil)
	return nil
}

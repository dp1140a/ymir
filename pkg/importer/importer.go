package importer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	cp "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"

	"ymir/pkg/api/model/store"
	"ymir/pkg/api/model/types"
	"ymir/pkg/utils"
)

type Importer struct {
	config     *ImporterConfig
	modelStore store.ModelStore
	baseDir    string
	copyModels bool
	modelsDir  string
	ymirHost   string
	Models     []types.Model
	Tags       []types.Tags
}

var method string

func NewImporter(base string, inDB bool, copy bool, modelsDir string, ymirHost string) *Importer {
	importer := &Importer{
		config:     NewImporterConfig(),
		baseDir:    base,
		copyModels: copy,
	}
	if modelsDir != "" {
		importer.modelsDir = modelsDir
	} else {
		importer.modelsDir = importer.config.modelsBase
	}

	if ymirHost != "" {
		importer.ymirHost = ymirHost
	} else {
		importer.ymirHost = fmt.Sprintf("%v:%v", viper.GetString("http.hostname"), viper.GetString("http.port"))
	}

	fmt.Printf("base (-p): %v\ninDB (-d): %v\n", base, inDB)
	fmt.Printf("modelsDir: %v\n", importer.config.modelsBase)
	return importer
}

func pingYmir(host string) bool {
	response, err := http.Get(fmt.Sprintf("http://%v/v1/ping", host))
	if err != nil {
		fmt.Printf("Connection Error: %v\n", err.Error())
	}
	if response != nil && response.StatusCode == 200 {
		return true
	}
	return false
}

func (i *Importer) PutInDB() {
	fmt.Printf("DB Flag is set will write to DB. Determining best method.")
	if pingYmir(i.ymirHost) {
		fmt.Printf("It's Alive. Found ymir host at %v.\n I shall import via the API\n", i.ymirHost)
		/**
		Ive chosen to make one import call per model.
		Its a bit more expensive but I can treat them atomically rather than have to handle transaction logic if some succeed and some fail
		*/
		for x, model := range i.Models {
			fmt.Printf("Importing model %v\n", x)
			body, _ := json.Marshal(model)
			response, err := http.Post(fmt.Sprintf("http://%v/v1/model/import", i.ymirHost), "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Printf("Import failed: %v\n   %v\nTrying next model.", model.BasePath, err.Error())
			}

			defer response.Body.Close()
			b, err := io.ReadAll(response.Body)
			fmt.Printf("   %v\n", string(b))
		}
	} else {
		fmt.Println("No response from the server. I shall try to import via the DB")
		fmt.Printf("Looking for db at %v\n", i.config.storeConfig.DBFile)
		ds := store.NewModelDataStore()
		if ds == nil {
			fmt.Printf("Could not find a ymir db file to use.")
			return
		} else {
			i.modelStore = ds.(store.ModelStore)
		}
		for x, model := range i.Models {
			model.Id = utils.GenId() //Since we go direct to db we bypass the id gen in the service so need it here
			fmt.Printf("Importing model %v\n", x)
			err := i.modelStore.Create(model)
			if err != nil {
				fmt.Printf("Import failed: %v\n   %v\nTrying next model.", model.BasePath, err.Error())
			}
			fmt.Printf("   Import model %v in db with docid %v\n", model.BasePath, model.Id)
		}
	}
}

var depth = 0
var curModelPath = ""

func (i *Importer) walk(path string, m *types.Model) error {
	fmt.Println("\n==============================================================")
	fmt.Printf("Scanning: %v [depth: %v]\n", path, depth)
	fmt.Println("==============================================================")
	dirs, _ := os.ReadDir(path)
	fmt.Printf("%v entry found: %v\n", len(dirs), dirs)
	for _, f := range dirs {
		fmt.Printf("\nAbs Path: %v\n", path)
		relPath, _ := filepath.Rel(i.baseDir, path)
		//relPath, _ = filepath.Rel(relPath, path)
		fmt.Printf("Rel Path: %v\n", relPath)
		fName := fmt.Sprintf("%v/%v", filepath.Base(relPath), f.Name())
		fmt.Printf("   Name(isDir): %v(%v)\n", fName, f.Type().IsDir())
		if f.IsDir() {
			if onlyDirectories(path) {
				fmt.Printf("%v is only directories moving on\n", path)
			}
			depth++
			i.walk(filepath.Join(path, f.Name()), m)
			depth--

		} else {
			if isExecutable(f.Type()) || !isModelFile(f.Name()) {
				fmt.Printf("File %v is either an executable or not a model file or image\n", filepath.Join(relPath, f.Name()))
				continue
			}
			if m == nil {
				curModelPath = path
				fmt.Println("Creating Model:")
				m = &types.Model{
					Id:          "",
					DisplayName: cleanDisplayName(filepath.Base(relPath)),
					BasePath:    curModelPath,
					Tags:        i.Tags,
					ModelFiles:  []types.FileType{},
					PrintFiles:  []types.FileType{},
					OtherFiles:  []types.FileType{},
					Images:      []types.FileType{},
					DateCreated: time.Now(),
					Notes:       []types.Note{},
				}
			}

			_, curPath := filepath.Split(filepath.Clean(curModelPath))
			fPath := filepath.Dir(filepath.Clean(fName))
			//fmt.Printf("relPath: %v / fName: %v \n", curPath, fPath)
			if curPath == fPath {
				fName = filepath.Base(fName)
			}

			/**
			If the files is README, README.txt or README.md then set the contents as the model description
			*/
			if strings.Contains(strings.ToLower(f.Name()), strings.ToLower("README")) {
				fmt.Printf("   Found a README File for model %v. Setting contents as model description\n", fName)

				m.OtherFiles = append(m.OtherFiles, types.FileType{Path: f.Name()})
				rm, err := os.ReadFile(fName)
				if err != nil {
					fmt.Printf("Cant open the README file")
				}
				m.Description = strconv.Quote(string(rm))
				continue
			}
			if slices.Contains(MODEL_TYPES, filepath.Ext(f.Name())[1:]) {
				fmt.Printf("  Adding model file: %v\n", f.Name())
				m.ModelFiles = append(m.ModelFiles, types.FileType{Path: fName})
				continue
			}
			if slices.Contains(PRINT_TYPES, filepath.Ext(f.Name())[1:]) {
				fmt.Printf("  Adding Print File: %v\n", f.Name())
				m.PrintFiles = append(m.PrintFiles, types.FileType{Path: fName})
				continue
			}
			if slices.Contains(IMAGE_TYPES, filepath.Ext(f.Name())[1:]) {
				fmt.Printf("   Adding Image: %v\n", f.Name())
				m.Images = append(m.Images, types.FileType{Path: fName})
				continue
			}
			if slices.Contains(OTHER_TYPES, filepath.Ext(f.Name())[1:]) {
				if f.Name() == "model.json" {
					fmt.Println("   Ignoring model.json as recursive (Dont cross the streams Egon!!)")
					continue
				}
				//fmt.Printf("  Adding Other File: %v\n", f.Name())
				m.OtherFiles = append(m.OtherFiles, types.FileType{Path: f.Name()})
				continue
			}
		}
	}

	//fmt.Printf("Depth: %v | mNull: %v\n", depth, m == nil)
	if path == curModelPath && m != nil {
		err := writeModel(path, m)
		if err != nil {
			return err
		}
		if i.copyModels {
			fmt.Printf("Copying model %v to %v\n", curModelPath, i.modelsDir)
			err = copyModel(curModelPath, i.modelsDir)
			if err != nil {
				fmt.Printf("Could not copy model %v to %v.  Skipping this model from import.\n", curModelPath, i.modelsDir)
				return err
			}
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
	fmt.Println("\n-----------------------------------------------------------")
	fmt.Printf("Writng model %v to %v/model.json\n", model.DisplayName, path)
	fmt.Println("-----------------------------------------------------------\n")
	err := model.WriteModel(path)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func copyModel(src string, dst string) error {
	return cp.Copy(src, filepath.Join(dst, src))
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
		ext := filepath.Ext(name)
		if ext != "" {
			ext = ext[1:]
		}
		if slices.Contains(fileType, ext) {
			fmt.Printf("   %v is a model file type.\n", name)
			return true
		}
	}
	return false
}

/*
This function is just a conveniience wrapper
*/
func (i *Importer) FindModels() error {
	i.walk(i.baseDir, nil)
	if len(i.Models) == 0 {
		fmt.Println("NO MODELS FOUND.  TRY ANOTHER DIRECTORY!")
	} else {
		fmt.Printf("%v MODELS FOUND.\n", len(i.Models))
	}
	return nil
}

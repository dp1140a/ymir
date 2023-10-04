package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Tags string

type Model struct {
	Id          string         `json:"_id,omitempty"`
	Rev         string         `json:"_rev,omitempty"`
	DisplayName string         `json:"displayName,omitempty"`
	Tags        []Tags         `json:"tags"`
	BasePath    string         `json:"basePath,omitempty"`
	ModelFiles  []FileType     `json:"modelFiles"`
	PrintFiles  []FileType     `json:"printFiles"`
	OtherFiles  []FileType     `json:"otherFiles"`
	Images      []FileType     `json:"images"`
	DateCreated time.Time      `json:"dateCreated,omitempty"`
	VersionLog  []ModelVersion `json:"versionLog,omitempty"`
	Description string         `json:"description,omitempty"`
	Summary     string         `json:"summary,omitempty"`
	Notes       []Note         `json:"notes"`
}

type FileType struct {
	Path string `json:"path,omitempty"`
}

type ModelVersion struct {
	DateModified time.Time `json:"dateModified"`
	ByWho        string    `json:"byWho"`
	VersionFlag  int       `json:"versionFlag,omitempty"`
	Comment      string    `json:"comment,omitempty"`
	Checksum     string    `json:",omitempty"`
}

type Note struct {
	Text string    `json:"text"`
	Date time.Time `json:"date"`
}

func testModel(num int) string {
	m := Model{
		Id:          uuid.New().String(),
		Rev:         "",
		DisplayName: fmt.Sprintf("test%v", num),
		Tags:        []Tags{"tag1", "tag2"},
		//Files:       []FileType{{"file1", "./file1"}, {"file2", "./file2"}},
		DateCreated: time.Now(),
		VersionLog:  nil,
		Description: "A test model for testing",
		Notes:       []Note{{"Note1", time.Now()}},
	}

	bytes, _ := json.MarshalIndent(m, "", "\t")
	return string(bytes)
}

func (m *Model) WriteModel(dir string) error {
	modelJSON := []byte(m.Json())
	if err := os.WriteFile(filepath.Join(dir, "model.json"), modelJSON, 0664); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (m *Model) Json() string {
	data, _ := json.MarshalIndent(m, "", "\t")
	return string(data)
}

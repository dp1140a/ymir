package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Tags string

type Model struct {
	Id          string         `json:"_id"`
	Rev         string         `json:"_rev"`
	DisplayName string         `json:"displayName"`
	Tags        []Tags         `json:"tags,omitempty"`
	Files       []FileType     `json:"files,omitempty"`
	DateCreated time.Time      `json:"dateCreated"`
	VersionLog  []ModelVersion `json:"versionLog,omitempty"`
	Description string         `json:"description,omitempty"`
	Notes       []Note         `json:"notes,omitempty"`
}

type FileType struct {
	Name string `json:"name"`
	Path string `json:"path"`
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
		Files:       []FileType{{"file1", "./file1"}, {"file2", "./file2"}},
		DateCreated: time.Now(),
		VersionLog:  nil,
		Description: "A test model for testing",
		Notes:       []Note{{"Note1", time.Now()}},
	}

	bytes, _ := json.MarshalIndent(m, "", "\t")
	return string(bytes)
}

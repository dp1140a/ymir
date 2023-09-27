package printer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

/*
*
host="127.0.0.1"
type="Prusa Mk3S+"
apiType="octoprint"
apikey="ABC123"
*/
type Printer struct {
	Id          string      `json:"_id,omitempty"`
	Rev         string      `json:"_rev,omitempty"`
	DisplayName string      `json:"printerName,omitempty"`
	URL         string      `json:"url"`
	APIType     string      `json:"apiType"`
	APIKey      string      `json:"apiKey"`
	Location    Location    `json:"location"`
	Type        PrinterType `json:"type"`
	DateAdded   time.Time   `json:"dateAdded"`
	Tags        []string    `json:"tags"`
}

type PrinterType struct {
	Make    string
	Model   string
	Version string
}

/*
*
Holding struct for future expansion
*/
type Location struct {
	Name string `json:"name"`
}

func (p *Printer) WriteModel(dir string) error {
	modelJSON := []byte(p.Json())
	if err := os.WriteFile(filepath.Join(dir, "model.json"), modelJSON, 0664); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (p *Printer) Json() string {
	data, _ := json.MarshalIndent(p, "", "\t")
	return string(data)
}

var TestPrinter = `{
	"_id": "4d3e3476-d7e8-4f34-957c-60e5fe1e29f3",
	"displayName": "test",
	"url": "http://myPrinter:8081",
	"apiType": "OctoPrint",
	"apiKey": "ABC123",
	"location": {
		"name": "Home"
	},
	"type": {
		"Make": "Prusa",
		"Model": "Mk3S+",
		"Version": "1.0"
	},
	"dateAdded": "2023-09-25T15:53:15.5652276-07:00",
	"tags": [
		"tag1",
		"tag2"
	]
}`

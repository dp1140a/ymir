package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	tNow, _ = time.Parse(time.RFC3339, "2023-09-25T15:53:15.5652276-07:00")
	t       = time.Now()

	testPrinter = Printer{
		Id:          "4d3e3476-d7e8-4f34-957c-60e5fe1e29f3",
		Rev:         "",
		PrinterName: "test",
		URL:         "http://myPrinter:8081",
		APIType:     "OctoPrint",
		APIKey:      "ABC123",
		Location: Location{
			Name: "Home",
		},
		Type: PrinterType{
			Make:    "Prusa",
			Model:   "Mk3S+",
			Version: "1.0",
		},
		DateAdded:   tNow,
		Tags:        []string{"tag1", "tag2"},
		AutoConnect: false,
	}
)

func TestPrinter_Json(t *testing.T) {
	tests := []struct {
		name    string
		printer Printer
		want    string
	}{
		{
			"testPrinter",
			testPrinter,
			`{
	"_id": "4d3e3476-d7e8-4f34-957c-60e5fe1e29f3",
	"printerName": "test",
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
	],
	"autoConnect": false
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := testPrinter
			assert.Equalf(t, tt.want, p.Json(), "Json()")
		})
	}
}

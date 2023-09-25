package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ymir/pkg/api"
)

func TestNewPrinterService(t *testing.T) {
	tests := []struct {
		name     string
		expected api.Service
	}{
		{
			"Default",
			PrintersService{
				name:   "Printers",
				config: NewPrintersConfig(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPrinterService()
			assert.NotNil(t, got, "Should Not be nil")
			assert.Equal(t, tt.expected.(PrintersService).name, got.GetName(), "Should be equal")
			assert.Equal(t, tt.expected.(PrintersService).config, got.(PrintersService).GetConfig(), "Should be equal")
		})
	}
}

func TestPrintersService_ListPrinters(t *testing.T) {
	tests := []struct {
		name string
		want []Printer
	}{
		{
			"Test Printer",
			[]Printer{
				Printer{
					Id:          "f7e5231e550c5f5bd36d4ccb2603a6d6",
					Rev:         "1-b352edf5ef13fea61218d1e877ca7b8f",
					DisplayName: "test",
					URL:         "http://myPrinter:8081",
					Tags:        []string{"tag1", "tag2"},
				},
			},
		},
	}
	for _, tt := range tests {
		ps := NewPrinterService()
		got, err := ps.(PrintersService).ListPrinters()
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, err, "Should be nil")
			assert.Len(t, got, 1, "Should be 1")
			assert.Equal(t, tt.want, got, "ListPrinters()")
		})
	}
}

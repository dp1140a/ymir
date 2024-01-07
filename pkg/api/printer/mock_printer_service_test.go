package printer

import (
	"github.com/stretchr/testify/mock"
	"ymir/pkg/api/printer/types"
)

// MockPrinterService is a mock implementation of the Service interface for testing.
type MockPrinterService struct {
	PrinterService
	mock.Mock
	printers []types.Printer
}

func NewMockPrinterService() *MockPrinterService {
	return &MockPrinterService{
		printers: getTestPrinters(2),
	}
}

func (m *MockPrinterService) CreatePrinter(types.Printer) (string, error) {
	return "", nil
}

func (m *MockPrinterService) ListPrinters() (map[string]types.Printer, error) {
	// Simulate returning a list of printers for testing.
	printers := map[string]types.Printer{}
	for i := 0; i < len(m.printers); i++ {
		printers[m.printers[i].Id] = m.printers[i]
	}
	return printers, nil
}

func (m *MockPrinterService) GetPrinter(id string) (types.Printer, error) {
	return m.printers[0], nil
}

func (m *MockPrinterService) UpdatePrinter(types.Printer) (err error) {
	return nil
}

func (m *MockPrinterService) DeletePrinter(id string) error {
	return nil
}

func (m *MockPrinterService) GetName() string {
	return ""
}

package printer

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewPrintersConfig(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		want       *PrintersConfig
	}{
		{
			"No Config File",
			"",
			&PrintersConfig{
				PrintersDir: "uploads/printers",
			},
		},
		{
			"With Good Config File",
			"testdata/goodConfig.toml",
			&PrintersConfig{
				PrintersDir: "uploads/printers2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.configFile != "" {
				//run viper first to get valid config
				viper.SetConfigFile(tt.configFile)
				if err := viper.ReadInConfig(); err != nil {
					t.Errorf("Error reading config file %v: %v\n", tt.configFile, err)
				}
			}
			assert.Equalf(t, tt.want, NewPrintersConfig(), "Should Be Equal")
		})
	}
}

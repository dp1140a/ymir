package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"ymir/pkg"
)

var defaultConfig = []byte(`printConfig = false

[logging]
logFile = "/var/log/ymir/ymir.log"
logLevel="INFO"
stdOut = true
fileOut = true

[datastore]
dbFile="~/.ymir/ymir.db"

[models]
uploadsTempDir="~/.ymir/uploads/tmp"
modelsDir="~/.ymir/models"

[printers]
printersDir="~/.ymir/printers"

[http]
hostname = "0.0.0.0"
port = "8081"
usehttps = false
TLSMinVersion = "1.2"
HttpTLSStrictCiphers = false
TLSCert = "ymir.crt"
TLSKey = "ymir.key"
enableCORS = true
JWTSecret = "abc123"

[http.logging]
enabled = true
stdOut = false
fileOut = true
logFile = "/var/log/ymir/ymir_http.log""`)

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	cfgFile := viper.GetString("cfgFile")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigType("toml")
		viper.SetConfigName("ymir")
		viper.AddConfigPath(path.Join(home, pkg.APP_NAME))    // Looks in ~/.APP_NAME
		viper.AddConfigPath(".")                              // Local dir
		viper.AddConfigPath(pkg.APP_NAME)                     // Looks in ./APP_NAME
		viper.AddConfigPath(path.Join("/etc/", pkg.APP_NAME)) // Looks in /etc/APP_NAME
		viper.AddConfigPath(home)                             // Looks in HOME
	}
	viper.SetEnvPrefix(pkg.APP_NAME)
	viper.AutomaticEnv() // read in environment variables that match
	loadDefaults()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	} else {
		if e, ok := err.(viper.ConfigParseError); ok { //Config Parsing Error
			fmt.Errorf("error parsing config file: %v\n", e)
			//log.Fatal("Exiting!")
		} else if _, ok := err.(*fs.PathError); ok {
			fmt.Printf("Config File Specified at %v Not Found.  Continuing with defaults.\n", cfgFile)
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok { // Config file not found; Use defaults
			fmt.Printf("Config File Specified at %v Not Found.  Continuing with defaults.\n", cfgFile)
		}
	}

	if viper.GetBool("printConfig") {
		fmt.Println(Toml())
	}
}

func loadDefaults() {
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		fmt.Errorf("%v", err)
	}
}

func Json() (config string) {
	cb, err := json.MarshalIndent(viper.AllSettings(), "", "   ")
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(-1)
	}
	return string(cb)
}

func Toml() (config string) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(viper.AllSettings())
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(-1)
	}
	return buf.String()
}

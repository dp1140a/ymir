package cmd

import (
	"unsafe"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"ymir/pkg/api/model"
	"ymir/pkg/importer"
)

var (
	Path string
	Tags []string
	inDb bool
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import command",
	Long:  `Model Import Command.  All params necessary are taken form the config used.  Either the default config or the config file specified as a global flag`,
	Run: func(cmd *cobra.Command, args []string) {
		runImport()
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&Path, "path", "p", "create", "Base directory path of models to import.  Will walk all directories from here")
	importCmd.Flags().StringSliceVarP(&Tags, "tagsr", "t", []string{}, "Tags for each model found")
	importCmd.Flags().BoolVarP(&inDb, "db", "d", false, "If true enters each model into the db")
}

func runImport() {
	imp := importer.NewImporter(Path)
	imp.Tags = *(*[]model.Tags)(unsafe.Pointer(&Tags))
	err := imp.FindModels()
	if err != nil {
		log.Fatal(err)
	}
	if inDb {
		err = imp.PutInDB()
	}
	if err != nil {
		log.Fatal(err)
	}

}

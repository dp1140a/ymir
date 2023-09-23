package cmd

import (
	"unsafe"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"ymir/pkg/api/model"
	"ymir/pkg/config"
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
	Long: `Model Import Command.  This command will import any models it finds into ymir.  without the db flag it will create a model.json file for each model.  
With the db flag it will make entries into the db.  Some assumptions were made that each model is in its own directory.  Ideally this looks like:

	.
	├── files
	│	├── MPR-1_0.2mm_PLA_MK3S_2h13m.gcode
	│	├── MPR-1.skp
	│	├── MPR-1.stl
	│	├── Unnamed-MPR_1001.step
	├── images
	│	├── MPR-1.jpg
	│	├── MPR-2.jpg
	│	├── MPR-3a.png
	│	└── MPR-3b.png
	├── LICENSE.txt
	└── README.txt
`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		runImport()
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&Path, "path", "p", ".", "Base directory path of models to import.  Will walk all sub-directories from there")
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

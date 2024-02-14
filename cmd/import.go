package cmd

import (
	"fmt"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"ymir/pkg/api/model/types"
	"ymir/pkg/config"
	"ymir/pkg/importer"
)

var (
	Path      string
	Tags      []string
	inDb      bool
	copy      bool
	modelsDir string
	ymirHost  string
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import command",
	Long: `Model Import Command.  This command will import any models it finds into ymir.  without the db flag it will create a model.json file for each model.  
With the db flag it will also make entries into the db.  Some assumptions were made that each model is in its own directory.  Ideally this looks like:

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
	importCmd.Flags().StringSliceVarP(&Tags, "tags", "t", []string{}, "Tags for each model found")
	importCmd.Flags().BoolVarP(&inDb, "db", "d", false, "If true enters each model into the db")
	importCmd.Flags().BoolVarP(&copy, "cp", "c", true, "Copy the models to ymir base models dir")
	importCmd.Flags().StringVarP(&modelsDir, "modelsDir", "m", "", "Path to the Ymir models dir as defined in the ymir config.  "+
		"This will override any value in the ymir config file. Use with caution.")
	importCmd.Flags().StringVarP(&ymirHost, "ymirHost", "y", "", "The host and port of the ymir server.")
}

func runImport() {
	imp := importer.NewImporter(Path, inDb, copy, modelsDir, ymirHost)
	imp.Tags = *(*[]types.Tags)(unsafe.Pointer(&Tags))
	err := imp.FindModels()
	if err != nil {
		log.Fatal(err)
	}
	if inDb {
		imp.PutInDB()
	}
	fmt.Println("IMPORT COMPLETE.")

}

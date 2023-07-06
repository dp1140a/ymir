/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	db2 "ymir/pkg/db"
)

var (
	Action       string
	actionValues = []string{"create", "drop", "export", "import", "truncate"}
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "db command",
	Long:  `Database Command.  All params necessary are taken form the config used.  Either the default config or the config file specified as a global flag`,
	Run: func(cmd *cobra.Command, args []string) {
		if slices.Contains(actionValues, Action) {
			run()

		} else {
			fmt.Printf("%s is not a valid action.\n", Action)
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.Flags().StringVarP(&Action, "action", "a", "create", "Database action to take [\"create\", \"drop\", \"export\", \"import\", \"truncate\"].  default: create")
}

func run() {
	db := db2.NewDB()
	db.Connect()

	switch Action {
	case "create":
		fmt.Println("Create Database")
		db.CreateDB()
	case "drop":
		fmt.Println("Drop Database")
		db.Drop()
	case "export":
		fmt.Println("Export Database Not Yet Implemented")
	case "import":
		fmt.Println("Import Database Not Yet Implemented")
	case "truncate":
		fmt.Println("Truncate Database")
		db.Truncate()

	}
}

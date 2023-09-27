package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ymir",
	Short: "A 3D File Manager for STL and GCode",
	Long:  `Ymir is a 3D model manager. In a nutshell it is a light and local version of the printables.com website`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//Global Flags
	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "ymir.toml", "config file (default is ./ymir.toml)")
	_ = viper.BindPFlag("cfgFile", rootCmd.PersistentFlags().Lookup("config"))
}

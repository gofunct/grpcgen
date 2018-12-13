package source

var RootTemplate = `package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/gofunct/gotasks/runtime/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               {{ .appname }},
	Short:             "write a short description of your app here",
	Version:           "0.1",
	PersistentPreRunE: viper.Viperize(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
`

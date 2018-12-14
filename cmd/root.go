package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
}

var (
	rootCmd = &cobra.Command{
		Use:   viper.GetString("root-usage"),
		Short: viper.GetString("root-summary"),
	}
)

// Execute executes the root command.
func Execute() {
	rootCmd.Execute()
}

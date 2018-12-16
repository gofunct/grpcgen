package cmd

import (
	"github.com/spf13/cobra"
)

var (
	service, packageName, parentName string
)

func init() {
	rootCmd.Flags().StringVar(&service, "service", "", "The protobuf message used for this configuration")
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
}

var (
	rootCmd = &cobra.Command{
		Use:   "gen",
		Short: "gen is a utility for easily creating highly configurable golang microservices",
	}
)

// Execute executes the root command.
func Execute() {
	rootCmd.Execute()
}

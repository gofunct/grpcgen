package cmd

import (
	"fmt"
	"github.com/gofunct/gen/viperizer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfg     *viperizer.Viperizer
	service, packageName, parentName string

)

func init() {
	var err error
	rootCmd.Flags().StringVar(&service, "service", "", "The protobuf message used for this configuration")
	cfg, err = viperizer.NewViperizer(service)
	if err != nil {
		logrus.Fatalf("failed to create config %v", err)
	}
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

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

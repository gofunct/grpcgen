package viperizer

import (
	"path"
	"path/filepath"
)

func (v *Viperizer) CreateMainFile() {
	mainTemplate := `package main

import "{{ .importpath }}"

func main() {
	cmd.Execute()
}
`
	data := make(map[string]interface{})
	data["importpath"] = path.Join(v.Project.GetName(), filepath.Base(v.Project.GetCmd()))

	mainScript, err := executeTemplate(mainTemplate, data)
	if err != nil {
		er(err)
	}

	err = WriteStringToFile(filepath.Join(v.Project.GetAbsPath(), "main.go"), mainScript)
	if err != nil {
		er(err)
	}
}

func (v *Viperizer) CreateRootCmdFile() {
	template := `package cmd

import (
	"github.com/gofunct/gen/viperizer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

var (
	viper *viperizer.Viperizer
	service string
)

func init() {
	var err error
	rootCmd.Flags().StringVar(&service, "service", "", "The protobuf message used for this configuration")
	viper, err = viperizer.NewViperizer(service)
	if err != nil {
		logrus.Fatalf("failed to create config %v", err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{.appName}}",
	Short: "A brief description of your application",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
`

	data := make(map[string]interface{})
	data["viper"] = true
	data["appName"] = path.Base(v.Project.GetName())

	rootCmdScript, err := executeTemplate(template, data)
	if err != nil {
		er(err)
	}

	err = WriteStringToFile(filepath.Join(v.Project.GetCmd(), "root.go"), rootCmdScript)
	if err != nil {
		er(err)
	}

}

func CreateCmdFile(path, cmdName, parent string) {
	template := `package {{.cmdPackage}}

import (
	"fmt"

	"github.com/spf13/cobra"
)

// {{.cmdName}}Cmd represents the {{.cmdName}} command
var {{.cmdName}}Cmd = &cobra.Command{
	Use:   "{{.cmdName}}",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("{{.cmdName}} called")
	},
}

func init() {{{.parentName}}.AddCommand({{.cmdName}}Cmd)}
`

	data := make(map[string]interface{})
	data["cmdPackage"] = filepath.Base(filepath.Dir(path)) // last dir of path
	data["parentName"] = parent
	data["cmdName"] = cmdName

	cmdScript, err := executeTemplate(template, data)
	if err != nil {
		er(err)
	}
	err = WriteStringToFile(path, cmdScript)
	if err != nil {
		er(err)
	}
}

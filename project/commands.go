package project

import (
	"github.com/gofunct/grpcgen/errors"
	"path"
	"path/filepath"
)

func (p *Project) CreateMainFile() {
	mainTemplate := `package main

import "{{ .importpath }}"

func main() {
	cmd.Execute()
}
`
	data := make(map[string]interface{})
	data["importpath"] = path.Join(p.GetName(), filepath.Base(p.GetCmd()))

	mainScript, err := ExecTemplate(mainTemplate, data)
	errors.IfErr("failed to execute template", err)

	err = WriteStringToFile(filepath.Join(p.GetAbsPath(), "main.go"), mainScript)
	errors.IfErr("failed to write file", err)

}

func (p *Project) CreateRootCmdFile() {
	template := `package cmd

import (
	"github.com/gofunct/grpcgen/viperizer"
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
		logrus.Fatalf("failed to create config %p", err)
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
	data["appName"] = path.Base(p.GetName())

	rootCmdScript, err := ExecTemplate(template, data)
	errors.IfErr("failed to execute template", err)
	err = WriteStringToFile(filepath.Join(p.GetCmd(), "root.go"), rootCmdScript)
	errors.IfErr("failed to write file", err)

}

func (p *Project) CreateCmdFile(path, cmdName, parent string) {
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

	cmdScript, err := ExecTemplate(template, data)
	errors.IfNoErr("failed to execute template", err)

	err = WriteStringToFile(path, cmdScript)
	errors.IfNoErr("failed to write file", err)
}

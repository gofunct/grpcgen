package cmd

import (
	"flag"
	"github.com/gofunct/grpcgen/logging"
	"github.com/gofunct/grpcgen/project"
	"github.com/gofunct/grpcgen/project/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"path/filepath"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		newProject.CreateProxyCmdFile()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

}


func CreateProxyCmdFile(g *Generator) {
	template := pr
	var service string
	pflag.StringVar(&service, "service", "", "The gRPC backend service to proxy.")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	data := make(map[string]interface{})
	data["viper"] = true
	data["service"] = service

	proxyScript, err := utils.ExecTemplate(template, data)
	logging.IfErr("failed to execute template", err)
	err = utils.WriteStringToFile(filepath.Join(g.Project.GetCmd(), "proxy.go"), proxyScript)
	logging.IfErr("failed to write file", err)

}




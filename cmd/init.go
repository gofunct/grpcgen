package cmd

import (
	"fmt"
	"github.com/gofunct/grpcgen/errors"
	"github.com/gofunct/grpcgen/project"
	"os"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init [name]",
	Aliases: []string{"initialize", "initialise", "create"},
	Short:   "Initialize a Cobra Application",
	Long: `Initialize (gen init) will create a new application
with the appropriate structure for a gen-based CLI application.`,

	Run: func(cmd *cobra.Command, args []string) {
		var p *project.Project
		wd, err := os.Getwd()
		errors.IfErr("failed to get working directory", err)

		if len(args) == 0 {
			p = project.NewProjectFromPath(wd)
		}

		project.InitializeProject(p)

		fmt.Fprintln(cmd.OutOrStdout(), `Your gen application is ready at
`+p.GetAbsPath()+`

Give it a try by going there and running `+"`go run main.go`."+`
Add commands to it by running `+"`gen add [cmdname]`.")
	},
}

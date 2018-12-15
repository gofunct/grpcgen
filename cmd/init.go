package cmd

import (
	"fmt"
	"github.com/gofunct/gen/viperizer"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init [name]",
	Aliases: []string{"initialize", "initialise", "create"},
	Short:   "Initialize a Cobra Application",
	Long: `Initialize (gen init) will create a new application
with the appropriate structure for a gen-based CLI application.`,

	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			er(err)
		}

		if len(args) == 0 {
			cfg.Project = viperizer.NewProjectFromPath(wd)
		} else if len(args) == 1 {
			arg := args[0]
			if arg[0] == '.' {
				arg = filepath.Join(wd, arg)
			}
			if filepath.IsAbs(arg) {
				cfg.Project = viperizer.NewProjectFromPath(arg)
			} else {
				cfg.Project = viperizer.NewProject(arg)
			}
		} else {
			er("please provide only one argument")
		}

		cfg.InitializeProject()

		fmt.Fprintln(cmd.OutOrStdout(), `Your gen application is ready at
`+cfg.Project.GetAbsPath()+`

Give it a try by going there and running `+"`go run main.go`."+`
Add commands to it by running `+"`gen add [cmdname]`.")
	},
}

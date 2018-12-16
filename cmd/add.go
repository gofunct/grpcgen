package cmd

import (
	"fmt"
	"github.com/gofunct/grpcgen/errors"
	"github.com/gofunct/grpcgen/project"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func init() {
	addCmd.Flags().StringVarP(&packageName, "package", "t", "", "target package name (e.g. github.com/spf13/hugo)")
	addCmd.Flags().StringVarP(&parentName, "parent", "p", "rootCmd", "variable name of parent command for this command")
}

var addCmd = &cobra.Command{
	Use:     "add [command name]",
	Aliases: []string{"command"},
	Short:   "Add a command to a Gen Application",
	Long: `Add (gen add) will create a new command, with the appropriate structure for a gen-based CLI application,
and register it to its parent (default rootCmd).

If you want your command to be public, pass in the command name
with an initial uppercase letter.

Example: gen add server -> resulting in a new cmd/server.go`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			errors.Exit("add needs a name for the command")
		}

		wd, err := os.Getwd()
		errors.IfErr("failed to get working directory", err)
		p := project.NewProjectFromPath(wd)


		cmdName := project.ValidateCmdName(args[0])
		cmdPath := filepath.Join(p.GetCmd(), cmdName+".go")
		pack := project.Project{
			CmdPath: cmdPath,
		}

		pack.CreateCmdFile(cmdPath, cmdName, parentName)

		fmt.Fprintln(cmd.OutOrStdout(), cmdName, "created at", cmdPath)
	},
}

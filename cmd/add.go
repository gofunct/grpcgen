package cmd

import (
	"fmt"
	"github.com/gofunct/gen/viperizer"
	"os"
	"path/filepath"
	"unicode"

	"github.com/spf13/cobra"
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
			er("add needs a name for the command")
		}

		if packageName != "" {
			cfg.Project = viperizer.NewProject(packageName)
		} else {
			wd, err := os.Getwd()
			if err != nil {
				er(err)
			}
			cfg.Project = viperizer.NewProjectFromPath(wd)
		}

		cmdName := validateCmdName(args[0])
		cmdPath := filepath.Join(cfg.Project.GetCmd(), cmdName+".go")
		viperizer.CreateCmdFile(cmdPath, cmdName, parentName)

		fmt.Fprintln(cmd.OutOrStdout(), cmdName, "created at", cmdPath)
	},
}

// validateCmdName returns source without any dashes and underscore.
// If there will be dash or underscore, next letter will be uppered.
// It supports only ASCII (1-byte character) strings.
// https://github.com/spf13/cobra/issues/269
func validateCmdName(source string) string {
	i := 0
	l := len(source)
	// The output is initialized on demand, then first dash or underscore
	// occurs.
	var output string

	for i < l {
		if source[i] == '-' || source[i] == '_' {
			if output == "" {
				output = source[:i]
			}

			// If it's last rune and it's dash or underscore,
			// don't add it output and break the loop.
			if i == l-1 {
				break
			}

			// If next character is dash or underscore,
			// just skip the current character.
			if source[i+1] == '-' || source[i+1] == '_' {
				i++
				continue
			}

			// If the current character is dash or underscore,
			// upper next letter and add to output.
			output += string(unicode.ToUpper(rune(source[i+1])))
			// We know, what source[i] is dash or underscore and source[i+1] is
			// uppered character, so make i = i+2.
			i += 2
			continue
		}

		// If the current character isn't dash or underscore,
		// just add it.
		if output != "" {
			output += string(source[i])
		}
		i++
	}

	if output == "" {
		return source // source is initially valid name.
	}
	return output
}

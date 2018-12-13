package source

var CmdTemplate = `package {{.cmdPackage}}

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

func init() {
	RootCmd.AddCommand({{.cmdName}}Cmd)
}
`

package cmd

import (
	"context"
	"github.com/gofunct/grpcgen/errors"
	"github.com/gofunct/grpcgen/example/todo"
	"github.com/gofunct/grpcgen/proxy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
var proxyViper = viper.New()

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
		prox := proxy.NewProxy(proxyViper, "todo_gateway")
		err := todo.RegisterTodoServiceHandlerFromEndpoint(ctx, prox.Gateway, viper.GetString("proxy.backend"), prox.DialOpts)
		errors.IfErr("failed to register gateway", err)
		prox.Listen(ctx)
	},
}

func init() {
	RootCmd.AddCommand(gatewayCmd)
}


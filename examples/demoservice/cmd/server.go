package cmd

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gofunct/grpcgen/example/services/session"
	"github.com/gofunct/grpcgen/example/services/session/endpoints"
	"github.com/gofunct/grpcgen/example/services/session/transports/grpc"
	"github.com/gofunct/grpcgen/example/services/session/transports/http"
	"github.com/gofunct/grpcgen/example/services/user"
	"github.com/gofunct/grpcgen/example/services/user/endpoints"
	"github.com/gofunct/grpcgen/example/services/user/transports/grpc"
	"github.com/gofunct/grpcgen/example/services/user/transports/http"
	"github.com/gorilla/handlers"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		mux := http.NewServeMux()
		errc := make(chan error)
		s := grpc.NewServer()
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			logger = log.With(logger, "ts", log.DefaultTimestampUTC)
			logger = log.With(logger, "caller", log.DefaultCaller)
		}

		// initialize services
		{
			svc := session.New()
			endpoints := session_endpoints.MakeEndpoints(svc)
			srv := session_grpctransport.MakeGRPCServer(endpoints)
			session.RegisterSessionServiceServer(s, srv)
			session_httptransport.RegisterHandlers(svc, mux, endpoints)
		}
		{
			svc := user.New()
			points := user_endpoints.MakeEndpoints(svc)
			srv := user_grpctransport.MakeGRPCServer(points)
			user.RegisterUserServiceServer(s, srv)
			user_httptransport.RegisterHandlers(svc, mux, points)
		}

		// start servers
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			errc <- fmt.Errorf("%s", <-c)
		}()

		go func() {
			logger := log.With(logger, "transport", "HTTP")
			logger.Log("addr", ":8000")
			errc <- http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stderr, mux))
		}()

		go func() {
			logger := log.With(logger, "transport", "gRPC")
			ln, err := net.Listen("tcp", ":9000")
			if err != nil {
				errc <- err
				return
			}
			logger.Log("addr", ":9000")
			errc <- s.Serve(ln)
		}()

		logger.Log("exit", <-errc)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}


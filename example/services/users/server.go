package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/handlers"
	"google.golang.org/grpc"

	session_svc "github.com/gofunct/grpcgen/services/session"
	session_endpoints "github.com/gofunct/grpcgen/services/session/gen/endpoints"
	session_pb "github.com/gofunct/grpcgen/services/session/gen/pb"
	session_grpctransport "github.com/gofunct/grpcgen/services/session/gen/transports/grpc"
	session_httptransport "github.com/gofunct/grpcgen/services/session/gen/transports/http"

	sprint_svc "github.com/gofunct/grpcgen/services/sprint"
	sprint_endpoints "github.com/gofunct/grpcgen/services/sprint/gen/endpoints"
	sprint_pb "github.com/gofunct/grpcgen/services/sprint/gen/pb"
	sprint_grpctransport "github.com/gofunct/grpcgen/services/sprint/gen/transports/grpc"
	sprint_httptransport "github.com/gofunct/grpcgen/services/sprint/gen/transports/http"

	user_svc "github.com/gofunct/grpcgen/services/user"
	user_endpoints "github.com/gofunct/grpcgen/services/user/gen/endpoints"
	user_pb "github.com/gofunct/grpcgen/services/user/gen/pb"
	user_grpctransport "github.com/gofunct/grpcgen/services/user/gen/transports/grpc"
	user_httptransport "github.com/gofunct/grpcgen/services/user/gen/transports/http"
)

var httpPort, grpcPort string

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&httpPort, "httpPort", "h", ":8000", "http transport port")
	serverCmd.Flags().StringVarP(&grpcPort, "grpcPort", "g", ":9000", "grpc transport port")
}

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "start a grpc server",
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
			svc := session_svc.New()
			endpoints := session_endpoints.MakeEndpoints(svc)
			srv := session_grpctransport.MakeGRPCServer(endpoints)
			session_pb.RegisterSessionServiceServer(s, srv)
			session_httptransport.RegisterHandlers(svc, mux, endpoints)
		}
		{
			svc := sprint_svc.New()
			endpoints := sprint_endpoints.MakeEndpoints(svc)
			srv := sprint_grpctransport.MakeGRPCServer(endpoints)
			sprint_pb.RegisterSprintServiceServer(s, srv)
			sprint_httptransport.RegisterHandlers(svc, mux, endpoints)
		}
		{
			svc := user_svc.New()
			endpoints := user_endpoints.MakeEndpoints(svc)
			srv := user_grpctransport.MakeGRPCServer(endpoints)
			user_pb.RegisterUserServiceServer(s, srv)
			user_httptransport.RegisterHandlers(svc, mux, endpoints)
		}

		// start servers
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			errc <- fmt.Errorf("%s", <-c)
		}()

		go func() {
			logger := log.With(logger, "transport", "HTTP")
			logger.Log("addr", httpPort)
			errc <- http.ListenAndServe(httpPort, handlers.LoggingHandler(os.Stderr, mux))
		}()

		go func() {
			logger := log.With(logger, "transport", "gRPC")
			ln, err := net.Listen("tcp", grpcPort)
			if err != nil {
				errc <- err
				return
			}
			logger.Log("addr", grpcPort)
			errc <- s.Serve(ln)
		}()

		logger.Log("exit", <-errc)
	},
}

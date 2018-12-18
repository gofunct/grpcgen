package project

import (
	"github.com/gofunct/grpcgen/logging"
	"os"
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
	logging.IfErr("failed to execute template", err)

	err = WriteStringToFile(filepath.Join(p.GetAbsPath(), "main.go"), mainScript)
	logging.IfErr("failed to write file", err)

}

func (p *Project) CreateRootCmdFile() {
	template := `package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

var (
	service string
)

func init() {
	rootCmd.Flags().StringVar(&service, "service", "", "The protobuf message used for this configuration")
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
	logging.IfErr("failed to execute template", err)
	err = WriteStringToFile(filepath.Join(p.GetCmd(), "root.go"), rootCmdScript)
	logging.IfErr("failed to write file", err)

}

func (p *Project) CreateServerCmdFile(path, cmdName, parent string) {
	template := `package cmd

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

func init() {rootCmd.AddCommand(serverCmd)}
`
	data := make(map[string]interface{})
	data["cmdPackage"] = filepath.Base(filepath.Dir(path)) // last dir of path
	data["parentName"] = parent

	cmdScript, err := ExecTemplate(template, data)
	logging.IfErr("failed to execute template", err)

	err = WriteStringToFile(path, cmdScript)
	logging.IfErr("failed to write file", err)
}


func (p *Project) CreateMakeFile() {
	mainTemplate := `SOURCES :=	$(shell find . -name "*.proto" -not -path ./vendor/\*)
TARGETS_GO :=	$(foreach source, $(SOURCES), $(source)_go)
TARGETS_TMPL :=	$(foreach source, $(SOURCES), $(source)_tmpl)

service_name =	$(word 2,$(subst /, ,$1))

.PHONY: setup
setup: ## download dependencies and tls certificates
	brew install prototool
	go get -u \
		google.golang.org/grpc \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/gogo/protobuf/protoc-gen-gogo \
		github.com/gogo/protobuf/protoc-gen-gogofast \
		github.com/ckaznocha/protoc-gen-lint \
		github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc \
		github.com/golang/protobuf/{proto,protoc-gen-go} \
		moul.io/protoc-gen-gotemplate

.PHONY: build
build: do ## install a grpc server

do: $(TARGETS_GO) $(TARGETS_TMPL)
	go install

$(TARGETS_GO): %_go:
	protoc --go_out=plugins=grpc:. "$*"
	@mkdir -p services/$(call service_name,$*)/gen/pb
	@mv ./services/$(call service_name,$*)/$(call service_name,$*).pb.go ./services/$(call service_name,$*)/gen/pb/pb.go

$(TARGETS_TMPL): %_tmpl:
	@mkdir -p $(dir $*)gen
	protoc -I. --gotemplate_out=destination_dir=services/$(call service_name,$*)/gen,template_dir=vendor/github.com/gofunct/grpcgen/templates:services "$*"
	@rm -rf services/services  # need to investigate why this directory is created
	gofmt -w $(dir $*)gen

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
`
	data := make(map[string]interface{})
	data["importpath"] = path.Join(p.GetName(), filepath.Base(p.GetCmd()))

	mainScript, err := ExecTemplate(mainTemplate, data)
	logging.IfErr("failed to execute template", err)

	err = WriteStringToFile(filepath.Join(p.GetAbsPath(), "Makefile"), mainScript)
	logging.IfErr("failed to write file", err)

}

func (p *Project) CreateDockerfile() {
	mainTemplate := `FROM golang
COPY    . "{{ .importpath }}"
WORKDIR "{{ .importpath }}"
CMD     ["{{ .appName }}"]
EXPOSE  8000 9000
RUN     make install
`
	data := make(map[string]interface{})
	data["importpath"] = path.Join(p.GetName(), filepath.Base(p.GetCmd()))
	data["appName"] = path.Base(p.GetName())

	mainScript, err := ExecTemplate(mainTemplate, data)
	logging.IfErr("failed to execute template", err)

	err = WriteStringToFile(filepath.Join(p.GetAbsPath(), "Dockerfile"), mainScript)
	logging.IfErr("failed to write file", err)

}


func (p *Project) CreateProtofile() {
	mainTemplate := `syntax = "proto3";
option go_package = "{{ .importpath }}";
import "google/protobuf/empty.proto";

package account_service;

message Account {
    string id = 1;
    string name = 2;
    string email = 3;
    string confirm_token = 5;
    string password_reset_token = 6;
    map<string, string> metadata = 7;
}

message ListAccountsRequest {
    int32 page_size = 1;
    string page_token = 2;
}

message ListAccountsResponse {
    repeated Account accounts = 1;
    string next_page_token = 2;
}

message GetByIdRequest {
    string id = 1;
}

message GetByEmailRequest {
    string email = 1;
}

message AuthenticateByEmailRequest {
    string email = 1;
    string password = 2;
}

message GeneratePasswordTokenRequest {
    string email = 1;
}

message GeneratePasswordTokenResponse {
    string token = 1;
}

message ResetPasswordRequest {
    string token = 1;
    string password = 2;
}

message ConfirmAccountRequest {
    string token = 1;
}

message CreateAccountRequest {
    Account account = 1;
    string password = 2;
}

message UpdateAccountRequest {
    string id = 1;
    string password = 2;
    Account account = 4;
}

message DeleteAccountRequest {
    string id = 1;
}

service AccountService {
    rpc List (ListAccountsRequest) returns (ListAccountsResponse) {}
    rpc GetById (GetByIdRequest) returns (Account) {}
    rpc GetByEmail (GetByEmailRequest) returns (Account) {}
    rpc AuthenticateByEmail (AuthenticateByEmailRequest) returns (Account) {}
    rpc GeneratePasswordToken (GeneratePasswordTokenRequest) returns (GeneratePasswordTokenResponse) {}
    rpc ResetPassword (ResetPasswordRequest) returns (Account) {}
    rpc ConfirmAccount (ConfirmAccountRequest) returns (Account) {}
    rpc Create (CreateAccountRequest) returns (Account) {}
    rpc Update (UpdateAccountRequest) returns (Account) {}
    rpc Delete (DeleteAccountRequest) returns (google.protobuf.Empty) {}
}
`
	data := make(map[string]interface{})
	data["importpath"] = path.Join(p.GetName(), filepath.Base(p.GetCmd()))

	mainScript, err := ExecTemplate(mainTemplate, data)
	logging.IfErr("failed to execute template", err)
	os.MkdirAll("services/accounts", os.ModePerm)
	err = WriteStringToFile(filepath.Join(p.GetAbsPath()+"/services/accounts", "accounts.proto"), mainScript)
	logging.IfErr("failed to write file", err)

}

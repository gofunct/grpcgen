package project

import (
	"github.com/gofunct/grpcgen/logging"
	"github.com/gofunct/grpcgen/project/utils"
	"path"
	"path/filepath"
)

func (p *Project) CreateMakeFile() {
	mainTemplate := `SOURCES :=	$(shell find . -name "*.proto" -not -path ./vendor/\*)
TARGETS_GO :=	$(foreach source, $(SOURCES), $(source)_go)
TARGETS_TMPL :=	$(foreach source, $(SOURCES), $(source)_tmpl)
import_path := {{ .importpath }}
app_name = {{ .appName }}
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

.PHONY: sessions
sessions: services/sessions/sessions.pb.go

services/sessions/sessions.pb.go:	services/sessions/sessions.proto
	@protoc --gotemplate_out=destination_dir=services/sessions,template_dir=$(GOPATH)/src/github.com/gofunct/grpcgen/templates:services/sessions services/sessions/sessions.proto
	gofmt -w services/sessions
	@protoc --gogo_out=plugins=grpc:. services/sessions/sessions.proto

.PHONY: users
users: services/users/users.pb.go

services/users/users.pb.go:	services/users/users.proto
	@protoc --gotemplate_out=destination_dir=services/users,template_dir=$(GOPATH)/src/github.com/gofunct/grpcgen/templates:services/users services/users/users.proto
	gofmt -w services/users
	@protoc --gogo_out=plugins=grpc:. services/users/users.proto

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
`
	data := make(map[string]interface{})
	data["importpath"] = path.Join(p.GetName(), filepath.Base(p.GetCmd()))
	data["appName"] = path.Base(p.GetName())
	mainScript, err := utils.ExecTemplate(mainTemplate, data)
	logging.IfErr("failed to execute template", err)

	err = utils.WriteStringToFile(filepath.Join(p.GetAbsPath(), "Makefile"), mainScript)
	logging.IfErr("failed to write file", err)

}

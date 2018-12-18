SOURCES :=	$(shell find . -name "*.proto" -not -path ./vendor/\*)
DOCKER_IMAGE ?=	moul/kafkagw
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

.PHONY: install
install: server ## install the service servers

server: $(TARGETS_GO) $(TARGETS_TMPL)
	go install .

$(TARGETS_GO): %_go:
	protoc --go_out=plugins=grpc:. "$*"
	@mkdir -p services/$(call service_name,$*)/gen/pb
	@mv ./services/$(call service_name,$*)/$(call service_name,$*).pb.go ./services/$(call service_name,$*)/gen/pb/pb.go

$(TARGETS_TMPL): %_tmpl:
	@mkdir -p $(dir $*)gen
	protoc -I. --gotemplate_out=destination_dir=services/$(call service_name,$*)/gen,template_dir=vendor/github.com/gofunct/grpcgen/templates:services "$*"
	@rm -rf services/services  # need to investigate why this directory is created
	gofmt -w $(dir $*)gen

.PHONY: stats
stats: ## stats
	wc -l service/service.go cmd/*/*.go pb/*.proto
	wc -l $(shell find gen -name "*.go")


.PHONY: test
test: ## run all unit tests
	go test -v $(shell go list ./... | grep -v /vendor/)

.PHONY: docker.build
docker.build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker.run
docker.run:
	docker run -p 8000:8000 -p 9000:9000 $(DOCKER_IMAGE)

.PHONY: docker.test
docker.test: docker.build
	docker run $(DOCKER_IMAGE) make test

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
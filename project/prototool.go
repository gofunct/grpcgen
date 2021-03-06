package project

import (
	"github.com/gofunct/grpcgen/logging"
	"github.com/gofunct/grpcgen/project/utils"
	"path"
	"path/filepath"
)

func (p *Project) CreatePrototoolfile() {
	mainTemplate := prototool
	data := make(map[string]interface{})
	data["importpath"] = path.Join(p.GetName(), filepath.Base(p.GetCmd()))
	data["appName"] = path.Base(p.GetName())

	mainScript, err := utils.ExecTemplate(mainTemplate, data)
	logging.IfErr("failed to execute template", err)

	err = utils.WriteStringToFile(filepath.Join(p.GetAbsPath(), "prototool.yaml"), mainScript)
	logging.IfErr("failed to write file", err)

}

var prototool = `# if starting from scratch in development, run "make init example" to get/build the .gitignored files
protoc:
  includes:
  # note vendor is .gitignored
  - ../../../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
lint:
  # run "prototool lint --list-linters" to see the currently configured linters
  # add "exclude_ids" to ignore specific linter IDs for all files
  ignores:
    - id: REQUEST_RESPONSE_TYPES_IN_SAME_FILE
      files:
        - foo/foo.proto
    - id: REQUEST_RESPONSE_TYPES_UNIQUE
      files:
        - foo/foo.proto
generate:
  go_options:
    import_path: github.com/uber/prototool/example/idl/uber
    extra_modifiers:
      google/api/annotations.proto: google.golang.org/genproto/googleapis/api/annotations
      google/api/http.proto: google.golang.org/genproto/googleapis/api/annotations
  plugins:
    - name: gogoslick
      type: gogo
      flags: plugins=grpc
      output: ../../gen/proto/go
    - name: yarpc-go
      type: gogo
      output: ../../gen/proto/go
      # note ../../gen/proto/java is .gitignored
    - name: java
      output: ../../gen/proto/java
    - name: grpc-gateway
      type: go
      output: ../../gen/proto/go
`

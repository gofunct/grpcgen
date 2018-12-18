package project

import (
	"github.com/gofunct/grpcgen/logging"
	"github.com/gofunct/grpcgen/project/utils"
	"os"
	"path/filepath"
)

func (p *Project) CreateSessionsProto() {
	mainTemplate := session
	data := make(map[string]interface{})

	mainScript, err := utils.ExecTemplate(mainTemplate, data)
	logging.IfErr("failed to execute template", err)
	os.MkdirAll("services/sessions", os.ModePerm)

	err = utils.WriteStringToFile(filepath.Join(p.GetAbsPath()+"/services/sessions", "sessions.proto"), mainScript)
	logging.IfErr("failed to write file", err)

}

var session = `syntax = "proto3";

package session;

service SessionService {
rpc Login(LoginRequest) returns (LoginResponse) {}
}

message LoginRequest {
string username = 1;
string password = 2;
}

message LoginResponse {
string token = 1;
string err_msg = 2;
}
`

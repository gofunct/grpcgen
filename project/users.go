package project


import (
	"github.com/gofunct/grpcgen/logging"
	"github.com/gofunct/grpcgen/project/utils"
	"os"
	"path/filepath"
)


func (p *Project) CreateUsersProto() {
	mainTemplate := users
	data := make(map[string]interface{})

	mainScript, err := utils.ExecTemplate(mainTemplate, data)
	logging.IfErr("failed to execute template", err)
	os.MkdirAll("services/users", os.ModePerm)

	err = utils.WriteStringToFile(filepath.Join(p.GetAbsPath()+"/services/users", "users.proto"), mainScript)
	logging.IfErr("failed to write file", err)

}

var users = `syntax = "proto3";

package user;

service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {}
}

message CreateUserRequest {
  string name = 1;
}
message CreateUserResponse {
  User user = 1;
  string err_msg = 2;
}

message GetUserRequest {
  string id = 1;
}
message GetUserResponse {
  User user = 1;
  string err_msg = 2;
}

message User {
  string id = 1;
  string name = 2;
}
`

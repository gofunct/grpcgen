package source

var ConfigTemplate = `app_name: {{ .appname }}
grpc_host: "localhost"
grpc_port: ":8443"
grpc_debug_port: ":8444"
db_port: ":5432"
db_host: "localhost"
db_pass: "admin"
db_name: "postgresdb"
db_user: "admin"
backend: ":8443"
cors:
  allow-origin:
  allow-credentials:
  allow-methods:
  allow-headers:
proxy:
  port: 8080
  api-prefix: "/"
swagger:
  file: "${SWAGGER_FILE_NAME}"
`

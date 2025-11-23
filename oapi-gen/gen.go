package oapi_gen

//go:generate oapi-codegen -o ../internal/gen/types.go -generate types -package gen openapi.yml
//go:generate oapi-codegen  -o ../internal/gen/server.go -generate gin-server,strict-server -package gen openapi.yml
//go:generate oapi-codegen  -o ../internal/gen/spec.go -generate spec -package gen openapi.yml

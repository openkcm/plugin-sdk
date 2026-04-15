module github.com/openkcm/plugin-sdk

go 1.25.4

tool (
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	github.com/openkcm/plugin-sdk/cmd/protoc-gen-go-extension
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260415201107-50325440f8f2.1
	buf.build/go/protovalidate v1.2.0
	github.com/hashicorp/go-hclog v1.6.3
	github.com/hashicorp/go-plugin v1.7.0
	github.com/stretchr/testify v1.11.1
	github.com/zeebo/errs/v2 v2.0.5
	golang.org/x/sys v0.43.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260420184626-e10c466a9529
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	cel.dev/expr v0.25.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/cel-go v0.28.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.21 // indirect
	github.com/oklog/run v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp v0.0.0-20260410095643-746e56fc9e2f // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260420184626-e10c466a9529 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.6.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

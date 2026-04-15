.PHONY: generate
generate: fetch-protos format-proto go-gen internal-go-gen

.PHONY: clean-proto
clean-proto:
	@find ./proto -type f -name '*.go' -exec rm {} +

.PHONY: clean-proto-internal
clean-proto-internal:
	@find ./internal/proto -type f -name '*.go' -exec rm {} +

.PHONY: fetch-protos
fetch-protos:
	@protofetch -o vendor-proto fetch

.PHONY: go-gen
go-gen: clean-proto
	@find ./proto -type f -iname '*.proto' -exec \
		protoc -I./proto -I./vendor-proto \
		--go_out=./proto \
		--go_opt=paths=import \
		--go_opt=module=github.com/openkcm/plugin-sdk/proto \
		--go-grpc_out=./proto \
		--go-grpc_opt=paths=import \
		--go-grpc_opt=module=github.com/openkcm/plugin-sdk/proto \
		--go-extension_out=./proto \
		--go-extension_opt=module=github.com/openkcm/plugin-sdk/proto \
		--go-extension_opt=submodule=github.com/openkcm/plugin-sdk/proto/service/common \
		--go-extension_opt=kind=service \
		--grpc-gateway_out=./proto \
		--grpc-gateway_opt=paths=import \
		--grpc-gateway_opt=module=github.com/openkcm/plugin-sdk/proto \
		--grpc-gateway_opt=logtostderr=true \
		{} +

	@find ./proto -type f -iname '*.proto' -exec \
		protoc -I./proto -I./vendor-proto \
		--go-extension_out=./proto \
		--go-extension_opt=module=github.com/openkcm/plugin-sdk/proto \
		--go-extension_opt=submodule=github.com/openkcm/plugin-sdk/proto/plugin \
		--go-extension_opt=kind=plugin \
		{} +

.PHONY: internal-go-gen
internal-go-gen: clean-proto-internal
	@find ./internal/proto -type f -iname '*.proto' -exec \
		protoc -I./internal/proto -I./proto -I./vendor-proto \
		--go_out=./internal/proto \
		--go_opt=paths=import \
		--go_opt=module=github.com/openkcm/plugin-sdk/internal/proto \
		--go-grpc_out=./internal/proto \
		--go-grpc_opt=paths=import \
		--go-grpc_opt=module=github.com/openkcm/plugin-sdk/internal/proto \
		--go-extension_out=./internal/proto \
		--go-extension_opt=module=github.com/openkcm/plugin-sdk/internal/proto \
		--go-extension_opt=submodule=github.com/openkcm/plugin-sdk/internal/proto/service \
		--go-extension_opt=kind=service \
		--grpc-gateway_out=./internal/proto \
		--grpc-gateway_opt=paths=import \
		--grpc-gateway_opt=module=github.com/openkcm/plugin-sdk/internal/proto \
		--grpc-gateway_opt=logtostderr=true \
		{} +

.PHONY: format-proto
format-proto:
	@buf format -w

.PHONY: format
format: format-proto

.PHONY: validate-proto
validate-proto: format-proto lint-proto breaking

.PHONY: lint-proto
lint-proto:
	@buf lint

.PHONY: breaking
breaking:
	@buf breaking --against https://github.com/openkcm/plugin-sdk.git#branch=main

.PHONY: install-proto-tools
install-proto-tools:
	brew install protobuf
	go install \
		google.golang.org/protobuf/cmd/protoc-gen-go@latest \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
		./cmd/protoc-gen-go-extension
	brew install bufbuild/buf/buf
	npm install -g @coralogix/protofetch

.PHONY: lint-install
lint-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: lint
lint: lint-proto
	golangci-lint run -v --fix

.PHONY: goimports
goimports:
	find ./ -name \*.go -not -path '.git/*' -exec goimports -local github.com/openkcm/plugin-sdk -w {} +

.PHONY: test
test:
	go test -race -coverprofile cover.out ./...
	# On a Mac, you can use the following command to open the coverage report in the browser
	# go tool cover -html=cover.out -o cover.html && open cover.html

.PHONY: reuse-lint
reuse-lint:
	docker run --rm --volume $(PWD):/data fsfe/reuse lint

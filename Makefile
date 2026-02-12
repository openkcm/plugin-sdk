.PHONY: generate

generate: patch-swagger-doc format
	go mod tidy
	go mod vendor

patch-swagger-doc: buf-gen internal-buf-gen
	#./scripts/update_swagger.sh docs/openapiv2/apidocs.swagger.json

init-git-hooks:
	git config --local core.hooksPath .githooks/

buf-gen: init-git-hooks
	buf dep update
	./buf.gen.yaml

internal-buf-gen:
	buf dep update
	./internal-buf.gen.yaml

format: buf-gen
	buf format -w

lint-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

lint:
	golangci-lint run -v --fix

goimports:
	find . -name \*.go -not -path '.git/*' -exec goimports -local github.com/openkcm/plugin-sdk -w {} \;

.PHONY: test
test:
	go test -race -coverprofile cover.out ./...
	# On a Mac, you can use the following command to open the coverage report in the browser
	# go tool cover -html=cover.out -o cover.html && open cover.html

.PHONY: reuse-lint
reuse-lint:
	docker run --rm --volume $(PWD):/data fsfe/reuse lint

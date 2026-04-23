[![REUSE status](https://api.reuse.software/badge/github.com/openkcm/plugin-sdk)](https://api.reuse.software/info/github.com/openkcm/plugin-sdk)

# Plugin SDK

## About this project

Service definitions, code generated stubs and infrastructure for running and testing KCM plugins.

## Overview

External plugins are separate processes and use
[go-plugin](https://github.com/hashicorp/go-plugin) under the covers.

KMS communicates with plugins over gRPC. As such, the various interfaces are defined via gRPC service definitions.

## Pre-requisites

Several tools are required to generate the code:

1. **`protoc compiler`**: see the instruction on the official [web site](https://protobuf.dev/installation) or install using homebrew `brew install protobuf`.
2. **`protoc-gen-go`**: install via `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`.
3. **`protoc-gen-go-grpc`**: install via `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`.
4. **`protoc-gen-grpc-gateway`**: install via `go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway`.
5. **`protoc-gen-openapiv2`**: install via `go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2`.
6. **`protofetch`**: download from the [releases page](https://github.com/coralogix/protofetch/releases/latest) of [the GitHub repository](https://github.com/coralogix/protofetch) or install via `cargo install protofetch` or `npm install -g @coralogix/protofetch`, or using any similar tool compatible with the npm registry.

## Optional tools

Optionally, the [`buf` CLI](https://github.com/bufbuild/buf) tool can be used to validate, detect breaking changes, and format `.proto` files:

- **`buf breaking --against https://github.com/openkcm/api-sdk.git#branch=main`**: detect breaking changes against the main branch of the remote repository.
- **`buf format -w`**: format `.proto` files.
- **`buf lint`**: lint `.proto` files.

## Makefile

There are several `make` targets defined in the `Makefile`:

- **`fetch-protos`**: download `.proto` dependencies using `protofetch`.
- **`generate`**: fetches `.proto` dependencies, formats `.proto` files, and generates Go code.
- **`install-proto-tools`**: installs the tools (including optional) from the following sources: Homebrew, Go registry (via `go install`), NPM registry (via `npm install -g`). See the target definition for the details.
- **`validate-proto`**: formats and lints `.proto` files, detects breaking changes.

For the rest `make` targets see `Makefile`.

## Dependencies

`.proto` dependencies are managed with the [`protofetch`](https://github.com/coralogix/protofetch) tool. This tool downloads `.proto` files from a specified location of a git repository and places them into the `vendor-proto` directory. The dependencies are specified in the `protofetch.toml` file.

For instance, a dependency on the [`protovalidate`](https://github.com/bufbuild/protovalidate) proto definitions can be specified as follow:

``` toml
name = "github.com/openkcm/plugin-sdk"
description = "Plugins SDK of the OpenKCM project"

[protovalidate]
url = "github.com/bufbuild/protovalidate"
revision = "v1.1.1"
content_roots = ["/proto/protovalidate"]
allow_policies = ["buf/validate/*"]
```

In order to fetch dependencies, execute:

``` sh
$ protofetch -o vendor-proto fetch
```

## Generate Go code from the .proto definitions

The code can be generated with executing the following Make target

```sh 
$ make generate
```

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/openkcm/plugin-sdk/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/openkcm/plugin-sdk/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/openkcm/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2025 SAP SE or an SAP affiliate company and OpenKCM contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/openkcm/plugin-sdk).

<p align="center"><img alt="Bundesministerium für Wirtschaft und Klimaschutz (BMWK)-EU funding logo" src="https://apeirora.eu/assets/img/BMWK-EU.png" width="400"/></p>

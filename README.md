# protocel

[![Coverage](https://coveralls.io/repos/Neakxs/protocel/badge.svg?branch=main&service=github)](https://coveralls.io/github/Neakxs/protocel?branch=main) [![GoReportCard](https://goreportcard.com/badge/github.com/Neakxs/protocel)](https://goreportcard.com/badge/github.com/Neakxs/protocel) ![GitHub](https://img.shields.io/github/license/Neakxs/protocel)

## About

This repository is a collection of protoc plugins based on the Common Expression Language :

- `protoc-gen-cel-authorize`, a plugin for writing authorization rules using gRPC metadata against protobuf messages
- `protoc-gen-cel-validate`, a plugin for writing validation rules on protobuf messages

The only language supported is [Go](https://go.dev/).

An example is located at [protocel-example](https://github.com/Neakxs/protocel-example) repository.

## Installation

For installating the plugin, you can simply run the `go install` command :

```shell
go install github.com/Neakxs/protocel/cmd/protoc-gen-go-cel-authorize
go install github.com/Neakxs/protocel/cmd/protoc-gen-go-cel-validate
```

The binary will be placed in your $GOBIN location.

## Configuration

For every protoc plugin defined here, it is possible to define global constats and functions for your proto definitions. It is possible using a protobuf option or through a configuration file, using the `config=path/to/config.yml` option.

Here is an example of a configuration file :

```yaml
options:
    globals:
        constants:
            sub: subject
        functions:
            isAdmin: "x-admin" in context.metadata
```

> When the same function is defined inside a protobuf file and in the configuration, the protobuf one is used.

## Usage

### protoc-gen-cel-authorize

For writing authorization rules, two variables are defined :

- `headers` (**map[string][]string** type), corresponding to gRPC metadata
- `request` (declared request message type), corresponding to incoming request

With the `headers` variable comes, two receiver functions are available for easier rules writing :
- `headers.get(string)`, equivalent to the go `func (http.Header) Get(string)` function
- `headers.values(string)`, equivalent to the go `func (http.Header) Values(string)`

#### Example

1. Create protobuf definition

```protobuf
syntax = "proto3";

package testdata.basic;
option go_package = "github.com/Neakxs/protocel/testdata/authorize/basic";

import "authorize/authorize.proto";
import "google/protobuf/empty.proto";

service BasicService {
    rpc Basic(BasicRequest) returns (google.protobuf.Empty) {
        option (protocel.authorize.method).expr = 'request.name == "name"';
    };
}

message BasicRequest {
    string name = 1;
}
```

2. Generate protobuf service
3. Implement gRPC service
4. Add interceptors to your gRPC server
5. Profit

### protoc-gen-cel-validate

For writing validation rules, everything contained in the message can be used.

#### Example

1. Create protobuf definition

```protobuf
syntax = "proto3";

package testdata.basic;
option go_package = "github.com/Neakxs/protocel/testdata/validate/basic";

import "validate/validate.proto";
import "google/api/field_behavior.proto";
import "google/protobuf/empty.proto";

service BasicService {
    rpc Basic(BasicRequest) returns (google.protobuf.Empty) {};
}

message BasicRequest {
    string name = 1 [
        (google.api.field_behavior) = REQUIRED,
        (protocel.validate.field).expr = 'name.startswith("names/")'
    ];
}
```

2. Generate protobuf code
3. Call `Validate` or `ValidateWithMask` functions when needed

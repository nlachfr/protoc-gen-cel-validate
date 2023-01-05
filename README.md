# protocel

[![Coverage](https://coveralls.io/repos/Neakxs/protocel/badge.svg?branch=main&service=github)](https://coveralls.io/github/Neakxs/protocel?branch=main) [![GoReportCard](https://goreportcard.com/badge/github.com/Neakxs/protocel)](https://goreportcard.com/badge/github.com/Neakxs/protocel) ![GitHub](https://img.shields.io/github/license/Neakxs/protocel)

## About

This repository is a collection of protoc plugins based on the Common Expression Language :

- `protoc-gen-cel-authorize`, a plugin for writing authorization rules (deprecated)
- `protoc-gen-cel-validate`, a plugin for writing validation rules

The supported language is [Go](https://go.dev/). More languages may be added in the future, depending on available CEL implementations.

With v0.2.0, the `protoc-gen-cel-validate` plugin features :

- complete support of the [CEL specification](https://github.com/google/cel-spec)
- support of multiple rules
- extensive configuration, with user defined constants, functions and overloads
- method validation based on RPC context and request
- message validation with cross fields reference and `google.protobuf.FieldMask` support
- recursive message validation using build-in validation functions
- support of the `google.api.field_behavior` REQUIRED annotation for enforcing a non default value ([AIP-203](https://google.aip.dev/203))
- support of the `google.api.resource_reference` annotations for enforcing matching patterns ([AIP-122](https://google.aip.dev/122))

Since `protoc-gen-cel-validate` is able to fully replace `protoc-gen-cel-authorize`, this plugin is now longer maintained.

> An example is located at [protocel-example](https://github.com/Neakxs/protocel-example) repository.
## Installation

For installating the plugin, you can simply run the `go install` command :

```shell
go install github.com/Neakxs/protocel/cmd/protoc-gen-go-cel-validate
```

The binary will be placed in your $GOBIN location.

## Configuration

The plugin is highly using protobuf options, configuration file (using the `config=path/to/config.yml` option) or both.
The configuration can be defined at various levels of the protobuf specification, with the following loading orders :

- configuration file > file option > service option > method option
- configuration file > file option > message option > field option

Here is an example :

- config.yml
```yaml
options:
    globals:
        constants:
            banned: my_banned_name
```

- example.proto
```protobuf
syntax = "proto3";

package example;
import "validate/validate.proto";
import "google/protobuf/empty.proto";

service Example {
    option (protocel.validate.service) = {
        options: {
            globals: {
                functions: [{
                    key: 'validateMessage'
                    value: 'request.validate()'
                }]
            }
        }
        expr: 'validateMessage()'
    };
    rpc Rpc(RpcRequest) returns (google.protobuf.Empty) {};
}

message RpcRequest {
    option (protocel.validate.message) = {
        options: {
            globals: {
                constants: [{
                    key: 'admin'
                    value: 'my_admin_name'
                }]
            }
        }
    };
    string name = 1 [(protocel.validate.field) = {
        expr: 'name != banned && name == admin'
    }];
}
```

## Usage of `protoc-gen-cel-validate`

For writing validation rules, some variables are defined, depending on the scope of the rule.

- for service and method rules, two variables are defined
  - `attribute_context` (google.rpc.context.AttributeContext), containing transport related metadata
  - `request` (declared request message type), corresponding to incoming request
- for message and field rules, all the fields of the message are defined

Furthermore, every message including validation rules provides the `validate()` and `validateWithMask(google.protobuf.FieldMask)` methods, allowing nested validation calls.

## Example

1. Create protobuf definition

```protobuf
syntax = "proto3";

package testdata.basic;
option go_package = "github.com/Neakxs/protocel/testdata/validate/basic";

import "validate/validate.proto";
import "google/api/field_behavior.proto";
import "google/protobuf/empty.proto";

service BasicService {
    rpc Basic(BasicRequest) returns (google.protobuf.Empty) {
        option (protocel.validate.method).expr = 'request.validate()';
    };
}

message BasicRequest {
    string name = 1 [
        (google.api.field_behavior) = REQUIRED,
        (protocel.validate.field).expr = 'name.startswith("names/")'
    ];
}
```

2. Generate protobuf code
3. For validating message, just call the `Validate` or `ValidateWithMask` methods on the corresponding messages. For validating methods of a service, build a `ServiceValidateProgram` using the generated builder and call the `Validate` method.
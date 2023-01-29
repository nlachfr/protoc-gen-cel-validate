<div align="center">
<h1>protocel</h1>
<p>Enforcing CEL validation rules with protobuf annotations</p>
<a href="https://coveralls.io/github/nlachfr/protocel?branch=main"><img src="https://coveralls.io/repos/nlachfr/protocel/badge.svg?branch=main&service=github"/></a>
<a href="https://goreportcard.com/badge/github.com/nlachfr/protocel"><img src="https://goreportcard.com/badge/github.com/nlachfr/protocel"/></a>
<a href="https://img.shields.io/github/license/nlachfr/protocel"><img src="https://img.shields.io/github/license/nlachfr/protocel"></a>
</div>

## About

*This project is still in alpha: APIs should be considered unstable and likely to change.*

Protocel is a plugin for the protocol buffers compiler. With the help of the Common Expression Language, this plugin reads user-defined rules on service and message definitions and generate ready to use validation functions.

It features :

- complete support of the [CEL specification](https://github.com/google/cel-spec)
- support of multiple rules
- extensive configuration, with user defined constants, functions and overloads
- method validation based on RPC context and request
- message validation with cross fields reference and `google.protobuf.FieldMask` support
- recursive message validation using build-in validation functions
- support of the `google.api.field_behavior` REQUIRED annotation for enforcing a non default value ([AIP-203](https://google.aip.dev/203))
- support of the `google.api.resource_reference` annotations for enforcing matching patterns ([AIP-122](https://google.aip.dev/122))

This repository contains two utilities :

- `protoc-gen-cel-validate`, the protoc plugin for writing and generating validation rules
- `protocel-gateway`, a small reverse proxy for handling validation rules without code generation

For now, the plugin is dedicated for the [Go](https://go.dev/) language. More languages may be added in the future, depending on available CEL implementations (and time). 

If you would like to use this project in another language, the `protocel-gateway` is made for you.

## Installation

For installating the plugin or the gateway, you can simply run the `go install` command :

```shell
go install github.com/nlachfr/protocel/cmd/protoc-gen-go-cel-validate
go install github.com/nlachfr/protocel/cmd/protocel-gateway
```

The binary will be placed in your $GOBIN location.

## Rules configuration

Protocel is highly configurable: validation rules can be written using protobuf options, a configuration file or both. 
Furthermore, options can be defined at various levels of your specification, with the following loading orders :

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

For more information on configuration fields, have a look at the [`protocel.validate.Options`](./validate/validate.proto) message specification.

## `protoc-gen-cel-validate`

### Configuration file

Even if rules can be written in the protobuf definition, you might not be able to edit the files for your use case. This is why you can use an external file for adding custom validation using the **config=/path/to/config.yml** parameter.

The configuration file will be loaded as a global `protocel.validate.Options` and will be used in all the generated files.
### Writing rules

For writing validation rules, some variables are defined, depending on the scope of the rule.

- for service and method rules, two variables are defined
  - `attribute_context` (google.rpc.context.AttributeContext), containing transport related metadata
  - `request` (declared request message type), corresponding to incoming request
- for message and field rules, all the fields of the message are defined

Furthermore, every message including validation rules provides the `validate()` and `validateWithMask(google.protobuf.FieldMask)` methods, allowing nested validation calls.

### Example

> An complete example is located at [protocel-example](https://github.com/nlachfr/protocel-example) repository.

1. Create protobuf definition

```protobuf
syntax = "proto3";

package testdata.basic;
option go_package = "github.com/nlachfr/protocel/testdata/validate/basic";

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

## `protocel-gateway`

### Configuration file

In order to use the gateway, you must provided a configuration file. It will allows you to :
- define listening address and upstreams
- specify proto files location
- configure protocel validation

All the configuration fields are specified in the [`protocel.gateway.Configuration`](./gateway/gateway.proto) message specification.
### Example

1. Create your configuration file

```yml
servers:
  - listen:
      - 127.0.0.1:8888
    upstreams:
      '*':
        address: https://demo.connect.build    
        protocol: grpc
files:
  sources:
    - proto/*.proto
  imports:
    - "."
validate:
  rule:
    servicerules:
      buf.connect.demo.eliza.v1.ElizaService:
        rule:
          programs:
            - expr: 'attribute_context.request.headers["ok"] == "ok"'      
      connect.ping.v1.PingService:
        rule:
          programs:
            - expr: 'attribute_context.request.headers["ok"] == "ok"'
```

2. Call protocel-gateway with your newly created configuration

```bash
$ protocel-gateway -config ./config.yml
```
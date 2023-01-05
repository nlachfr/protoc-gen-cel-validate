package validate

import (
	"context"
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/Neakxs/protocel/testdata/validate/option"
	"github.com/google/cel-go/cel"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestBuildServiceValidateProgram(t *testing.T) {
	tests := []struct {
		Name    string
		Desc    protoreflect.ServiceDescriptor
		Config  *ValidateOptions
		EnvOpt  cel.EnvOption
		WantErr bool
	}{
		{
			Name:    "No validation",
			Desc:    validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("Service")),
			WantErr: false,
		},
		{
			Name:    "Service level expr",
			Desc:    validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceExpr")),
			WantErr: false,
		},
		{
			Name:    "Service level expr with missing const",
			Desc:    validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceOptions")),
			WantErr: true,
		},
		{
			Name: "Service level expr with global const",
			Desc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceOptions")),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"isAdmHdr": "x-is-admin",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:    "Service level expr with local const",
			Desc:    validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceLocalOptions")),
			WantErr: false,
		},
		{
			Name: "Service level expr with const conflict",
			Desc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceLocalOptions")),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"isAdmHdr": "x-is-admin",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:    "Method level expr",
			Desc:    validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodExpr")),
			WantErr: false,
		},
		{
			Name:    "Method level with missing const",
			Desc:    validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodOptions")),
			WantErr: true,
		},
		{
			Name: "Method level with global const",
			Desc: validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodOptions")),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"isAdmHdr": "x-is-admin",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:    "Method level with local const",
			Desc:    validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodLocalOptions")),
			WantErr: false,
		},
		{
			Name: "Method level with const conflict",
			Desc: validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodLocalOptions")),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"isAdmHdr": "x-is-admin",
						},
					},
				},
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := BuildServiceValidateProgram(tt.Config, tt.Desc, tt.EnvOpt)
			if (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

func TestBuildMethodValidateProgram(t *testing.T) {
	tests := []struct {
		Name    string
		Expr    string
		Method  string
		Desc    protoreflect.MethodDescriptor
		Context *attribute_context.AttributeContext
		Request proto.Message
		EnvOpt  cel.EnvOption
		WantErr bool
	}{
		{
			Name:   "Method name mistmatch",
			Expr:   `false`,
			Method: "operation",
			Desc:   validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{
					Operation: "another.Operation",
				},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: false,
		},
		{
			Name: "Headers map is nil",
			Expr: `attribute_context.request.headers["ok"] == "ok"`,
			Desc: validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: true,
		},
		{
			Name: "Nil request",
			Expr: `request.ref == "ref"`,
			Desc: validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: nil,
			WantErr: true,
		},
		{
			Name: "Attribute context validation failed",
			Expr: `attribute_context.request.headers["ok"] == "ok"`,
			Desc: validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
				Request: &attribute_context.AttributeContext_Request{
					Headers: map[string]string{"ok": ""},
				},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: true,
		},
		{
			Name: "Request validation failed",
			Expr: `request.ref == "ref"`,
			Desc: validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: true,
		},
		{
			Name: "Attribute Context & Request validation succeeded",
			Expr: `attribute_context.request.headers["ref"] == request.ref`,
			Desc: validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
				Request: &attribute_context.AttributeContext_Request{
					Headers: map[string]string{"ref": "myRef"},
				},
			},
			Request: &validate.TestRpcRequest{Ref: "myRef"},
		},
		{
			Name: "Request validation with missing variable",
			Expr: `myVariable == request.ref`,
			Desc: validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &validate.TestRpcRequest{Ref: "myRef"},
			WantErr: true,
		},
		{
			Name: "Request validation with variable succeeded",
			Expr: `myVariable == request.ref`,
			Desc: validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			EnvOpt: cel.Lib(&options.Library{
				EnvOpts: []cel.EnvOption{
					cel.Variable("myVariable", cel.StringType),
				},
				PgrOpts: []cel.ProgramOption{
					cel.Globals(map[string]interface{}{
						"myVariable": "myRef",
					}),
				},
			}),
			Request: &validate.TestRpcRequest{Ref: "myRef"},
		},
		{
			Name: "Request validation call succeeded",
			Expr: `request.validate()`,
			Desc: validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			EnvOpt:  nil,
			Request: &validate.TestRpcRequest{Ref: "refs/myref"},
		},
		{
			Name: "Request validation with method defined options",
			Expr: `request.name == myMethodConst`,
			Desc: option.File_testdata_validate_option_option_proto.Services().Get(0).Methods().Get(0),
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			EnvOpt:  nil,
			Request: &option.OptionRequest{Name: "name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var ferr error
			if pgr, err := BuildMethodValidateProgram([]string{tt.Expr}, nil, tt.Desc, tt.EnvOpt); err != nil {
				ferr = err
			} else {
				ferr = (&serviceValidateProgram{map[string]ValidateProgram{tt.Method: pgr}}).Validate(context.Background(), tt.Context, tt.Request)
			}
			if (ferr != nil && !tt.WantErr) || (ferr == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, ferr)
			}
		})
	}
}

package validate

import (
	"context"
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
)

func TestValidateInterceptor(t *testing.T) {
	env, err := cel.NewEnv(
		cel.TypeDescs(validate.File_testdata_validate_test_proto, attribute_context.File_google_rpc_context_attribute_context_proto),
		cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
		cel.Types((&validate.TestRpcRequest{}).ProtoReflect().New().Interface()),
		cel.Variable("request", cel.ObjectType(string((&validate.TestRpcRequest{}).ProtoReflect().Descriptor().FullName()))),
	)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		Name    string
		Expr    string
		Method  string
		Context *attribute_context.AttributeContext
		Request *validate.TestRpcRequest
		WantErr bool
	}{
		{
			Name:   "Method name mistmatch",
			Expr:   `false`,
			Method: "operation",
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{
					Operation: "another.Operation",
				},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: false,
		},
		{
			Name:    "Attribute context is nil",
			Expr:    `false`,
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: false,
		},
		{
			Name: "Headers map is nil",
			Expr: `attribute_context.request.headers["ok"] == "ok"`,
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: true,
		},
		{
			Name: "Nil request",
			Expr: `request.ref == "ref"`,
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: nil,
			WantErr: true,
		},
		{
			Name: "Attribute context validation failed",
			Expr: `attribute_context.request.headers["ok"] == "ok"`,
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
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: true,
		},
		{
			Name: "Attribute Context & Request validation succeeded",
			Expr: `attribute_context.request.headers["ref"] == request.ref`,
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
				Request: &attribute_context.AttributeContext_Request{
					Headers: map[string]string{"ref": "myRef"},
				},
			},
			Request: &validate.TestRpcRequest{Ref: "myRef"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if ast, err := env.Compile(tt.Expr); err != nil {
				t.Error(err)
			} else if pgr, err := env.Program(ast); err != nil {
				t.Error(err)
			} else {
				err := NewValidateInterceptor(map[string]*Program{tt.Method: {rules: []cel.Program{pgr}}}).Validate(context.Background(), tt.Context, tt.Request)
				if (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
					t.Errorf("wantErr %v, got %v", tt.WantErr, err)
				}
			}
		})
	}
}

func TestBuildMethodValidateProgram(t *testing.T) {
	tests := []struct {
		Name    string
		Expr    string
		Method  string
		Context *attribute_context.AttributeContext
		Request *validate.TestRpcRequest
		EnvOpt  cel.EnvOption
		WantErr bool
	}{
		{
			Name:   "Method name mistmatch",
			Expr:   `false`,
			Method: "operation",
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
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: true,
		},
		{
			Name: "Nil request",
			Expr: `request.ref == "ref"`,
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: nil,
			WantErr: true,
		},
		{
			Name: "Attribute context validation failed",
			Expr: `attribute_context.request.headers["ok"] == "ok"`,
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
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &validate.TestRpcRequest{Ref: ""},
			WantErr: true,
		},
		{
			Name: "Attribute Context & Request validation succeeded",
			Expr: `attribute_context.request.headers["ref"] == request.ref`,
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
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &validate.TestRpcRequest{Ref: "myRef"},
			WantErr: true,
		},
		{
			Name: "Request validation with variable succeeded",
			Expr: `myVariable == request.ref`,
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
			Context: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			EnvOpt:  nil,
			Request: &validate.TestRpcRequest{Ref: "refs/myref"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var ferr error
			if pgr, err := BuildMethodValidateProgram([]string{tt.Expr}, nil, validate.File_testdata_validate_test_proto.Services().Get(0).Methods().Get(0), tt.EnvOpt); err != nil {
				ferr = err
			} else {
				ferr = NewValidateInterceptor(map[string]*Program{tt.Method: pgr}).Validate(context.Background(), tt.Context, tt.Request)
			}
			if (ferr != nil && !tt.WantErr) || (ferr == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, ferr)
			}
		})
	}
}

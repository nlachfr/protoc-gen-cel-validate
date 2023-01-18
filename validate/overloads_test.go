package validate

import (
	"context"
	"testing"

	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestDefaultOverloadBuilder(t *testing.T) {
	tests := []struct {
		Name           string
		Message        proto.Message
		Config         *Options
		EnvOpt         cel.EnvOption
		WantCompileErr bool
		WantEvalErr    bool
		WantFalse      bool
	}{
		{
			Name:        "NOK",
			Message:     &validate.TestRpcRequest{},
			WantEvalErr: true,
		},
		{
			Name:    "OK ",
			Message: &validate.TestRpcRequest{Ref: "refs/myref", Raw: "raw"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			opts := (&defaultOverloadBuilder{}).buildOverloads(tt.Message.ProtoReflect().Descriptor())
			opts = append(opts, cel.TypeDescs(tt.Message.ProtoReflect().Descriptor().ParentFile()),
				cel.TypeDescs(fieldmaskpb.File_google_protobuf_field_mask_proto),
				cel.Variable("fm", cel.ObjectType(string((&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName()))),
				cel.Variable("myvar", cel.ObjectType(string(tt.Message.ProtoReflect().Descriptor().FullName()))),
			)
			env, err := cel.NewEnv(opts...)
			if err != nil {
				t.Error(err)
			}
			for _, expr := range []string{`myvar.validate()`, `myvar.validateWithMask(fm)`} {
				ast, issues := env.Compile(expr)
				if issues != nil {
					err = issues.Err()
				}
				if (tt.WantCompileErr && err == nil) || (!tt.WantCompileErr && err != nil) {
					t.Errorf("wantCompileErr %v, got %v", tt.WantCompileErr, err)
				}
				pgr, err := env.Program(ast)
				if err != nil {
					t.Error(err)
				} else {
					val, _, err := pgr.ContextEval(context.Background(), map[string]interface{}{"myvar": tt.Message, "fm": &fieldmaskpb.FieldMask{Paths: []string{"*"}}})
					if err == nil {
						if e, ok := val.Value().(error); ok {
							err = e
						}
					}
					if (tt.WantEvalErr && err == nil) || (!tt.WantEvalErr && err != nil) {
						t.Errorf("wantEvalErr %v, got %v", tt.WantEvalErr, err)
					} else if err == nil {
						if (tt.WantFalse && val.Value().(bool)) || (!tt.WantFalse && !val.Value().(bool)) {
							t.Errorf("wantFalse %v, got %v", tt.WantFalse, val.Value().(bool))
						}
					}
				}
			}
		})
	}
}
func TestFallbackOverloadBuilder(t *testing.T) {
	tests := []struct {
		Name           string
		Message        proto.Message
		Config         *Options
		EnvOpt         cel.EnvOption
		WantCompileErr bool
		WantEvalErr    bool
		WantFalse      bool
	}{
		{
			Name:        "NOK",
			Message:     &validate.TestRpcRequest{Ref: "r"},
			WantEvalErr: true,
		},
		{
			Name:    "OK ",
			Message: &validate.TestRpcRequest{Ref: "refs/myref", Raw: "raw"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			opts := (&fallbackOverloadBuilder{
				Builder: newBuilder(),
			}).buildOverloads(tt.Message.ProtoReflect().Descriptor())
			opts = append(opts, cel.TypeDescs(tt.Message.ProtoReflect().Descriptor().ParentFile()),
				cel.TypeDescs(fieldmaskpb.File_google_protobuf_field_mask_proto),
				cel.Variable("fm", cel.ObjectType(string((&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName()))),
				cel.Variable("myvar", cel.ObjectType(string(tt.Message.ProtoReflect().Descriptor().FullName()))),
			)
			env, err := cel.NewEnv(opts...)
			if err != nil {
				t.Error(err)
			}
			for _, expr := range []string{`myvar.validate()`, `myvar.validateWithMask(fm)`} {
				ast, issues := env.Compile(expr)
				if issues != nil {
					err = issues.Err()
				}
				if (tt.WantCompileErr && err == nil) || (!tt.WantCompileErr && err != nil) {
					t.Errorf("wantCompileErr %v, got %v", tt.WantCompileErr, err)
				}
				pgr, err := env.Program(ast)
				if err != nil {
					t.Error(err)
				} else {
					val, _, err := pgr.ContextEval(context.Background(), map[string]interface{}{"myvar": tt.Message, "fm": &fieldmaskpb.FieldMask{Paths: []string{"*"}}})
					if err == nil {
						if e, ok := val.Value().(error); ok {
							err = e
						}
					}
					if (tt.WantEvalErr && err == nil) || (!tt.WantEvalErr && err != nil) {
						t.Errorf("wantEvalErr %v, got %v", tt.WantEvalErr, err)
					} else if err == nil {
						if (tt.WantFalse && val.Value().(bool)) || (!tt.WantFalse && !val.Value().(bool)) {
							t.Errorf("wantFalse %v, got %v", tt.WantFalse, val.Value().(bool))
						}
					}
				}
			}
		})
	}
}

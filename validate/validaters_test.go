package validate

import (
	"context"
	"testing"

	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var tests = []struct {
	Name           string
	Message        proto.Message
	Config         *ValidateOptions
	EnvOpt         cel.EnvOption
	WantCompileErr bool
	WantEvalErr    bool
	WantFalse      bool
}{
	{
		Name:        "NOK (Validater implemented)",
		Message:     &validate.TestRpcRequest{},
		WantEvalErr: true,
	},
	{
		Name:        "NOK (Validater not implemented)",
		Message:     &validate.FieldExpr{Name: "name"},
		WantEvalErr: true,
	},
	{
		Name:    "OK (Validater not implemented)",
		Message: &validate.FieldExpr{Name: "notname"},
	},
	{
		Name:    "OK (Validater implemented)",
		Message: &validate.TestRpcRequest{Ref: "refs/myref", Raw: "raw"},
	},
}

func TestValidateFunctionOpt(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.TypeDescs(tt.Message.ProtoReflect().Descriptor().ParentFile()),
				cel.Variable("myvar", cel.ObjectType(string(tt.Message.ProtoReflect().Descriptor().FullName()))),
				(&validateOverloadBuilder{
					config: tt.Config,
					envOpt: tt.EnvOpt,
				}).buildValidateFunction(tt.Message.ProtoReflect().Descriptor()),
			)
			if err != nil {
				t.Error(err)
			}
			ast, issues := env.Compile(`myvar.validate()`)
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
				val, _, err := pgr.ContextEval(context.Background(), map[string]interface{}{"myvar": tt.Message})
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
		})
	}
}

func TestValidateWithMaskFunctionOpt(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.TypeDescs(tt.Message.ProtoReflect().Descriptor().ParentFile()),
				cel.TypeDescs(fieldmaskpb.File_google_protobuf_field_mask_proto),
				cel.Variable("myvar", cel.ObjectType(string(tt.Message.ProtoReflect().Descriptor().FullName()))),
				cel.Variable("fm", cel.ObjectType(string((&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName()))),
				(&validateOverloadBuilder{
					config: tt.Config,
					envOpt: tt.EnvOpt,
				}).buildValidateWithMaskFunction(tt.Message.ProtoReflect().Descriptor()),
			)
			if err != nil {
				t.Error(err)
			}
			ast, issues := env.Compile(`myvar.validateWithMask(fm)`)
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
		})
	}
}

package validate

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/nlachfr/protocel/testdata/validate"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestBuildRuleValidater(t *testing.T) {
	tests := []struct {
		Name      string
		Rule      *Rule
		Desc      protoreflect.MessageDescriptor
		EnvOption cel.EnvOption
		WantErr   bool
	}{
		{
			Name: "Unknown field",
			Rule: &Rule{
				Programs: []*Rule_Program{{Expr: `name`}},
			},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: true,
		},
		{
			Name: "Invalid return type",
			Rule: &Rule{
				Programs: []*Rule_Program{{Expr: `name`}},
			},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: true,
		},
		{
			Name: "Invalid validate call on standard type",
			Rule: &Rule{
				Programs: []*Rule_Program{{Expr: `ref.validate()`}},
			},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: true,
		},
		{
			Name: "Unknown field in macro",
			Rule: &Rule{
				Options: &Options{
					Globals: &Options_Globals{
						Functions: map[string]string{
							"macro": `name == "name"`,
						},
					},
				},
				Programs: []*Rule_Program{{Expr: `macro()`}},
			},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: true,
		},
		{
			Name: "Regexp error",
			Rule: &Rule{
				Options: &Options{
					Overloads: &Options_Overloads{
						Variables: map[string]*Options_Overloads_Type{
							"myVariable": {Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							}},
						},
					},
				},
				Programs: []*Rule_Program{{Expr: `ref.matches("[")`}},
			},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: true,
		},
		{
			Name: "OK",
			Rule: &Rule{
				Programs: []*Rule_Program{{Expr: `ref == "ref"`}},
			},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: false,
		},
		{
			Name: "OK (with constant)",
			Rule: &Rule{
				Options: &Options{
					Globals: &Options_Globals{
						Constants: map[string]string{
							"constRef": "ref",
						},
					},
				},
				Programs: []*Rule_Program{{Expr: `ref == constRef`}},
			},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: false,
		},
		{
			Name: "OK (with macro)",
			Rule: &Rule{
				Options: &Options{
					Globals: &Options_Globals{
						Functions: map[string]string{
							"rule": `ref`,
						},
					},
				},
				Programs: []*Rule_Program{{Expr: `rule() == ref`}},
			},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: false,
		},
		{
			Name: "OK (with variable)",
			Rule: &Rule{
				Options: &Options{
					Overloads: &Options_Overloads{
						Variables: map[string]*Options_Overloads_Type{
							"myVariable": {Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							}},
						},
					},
				},
				Programs: []*Rule_Program{{Expr: `ref == myVariable`}},
			},
			Desc: (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			EnvOption: cel.Lib(&Library{
				PgrOpts: []cel.ProgramOption{cel.Globals(map[string]interface{}{"myVariable": "ref"})},
			}),
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			lib := &Library{}
			if tt.EnvOption != nil {
				lib.EnvOpts = append(lib.EnvOpts, tt.EnvOption)
			}
			lib.EnvOpts = append(lib.EnvOpts, cel.DeclareContextProto(tt.Desc))
			if tt.Rule != nil {
				lib.EnvOpts = append(lib.EnvOpts, BuildEnvOption(tt.Rule.Options, tt.Desc))
			}
			_, err := BuildRuleValidater(tt.Rule, cel.Lib(lib))
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

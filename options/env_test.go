package options

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestBuildEnvOption(t *testing.T) {
	tests := []struct {
		Name    string
		Expr    string
		Desc    protoreflect.MessageDescriptor
		EnvOpt  cel.EnvOption
		Config  *Options
		WantErr bool
	}{
		{
			Name:    "No options",
			Expr:    "1 == 2",
			WantErr: false,
		},
		{
			Name:    "Constant (error)",
			Expr:    `name == "name"`,
			WantErr: true,
		},
		{
			Name: "Constant",
			Expr: `name == "name"`,
			Config: &Options{
				Globals: &Options_Globals{
					Constants: map[string]string{"name": "name"},
				},
			},
			WantErr: false,
		},
		{
			Name: "Overload variable (error)",
			Expr: `name == "name"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Variables: map[string]*Options_Overloads_Type{"name": {
						Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}},
					},
				},
			},
			WantErr: true,
		},
		{
			Name: "Overload variable",
			Expr: `name == "name"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Variables: map[string]*Options_Overloads_Type{"name": {
						Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}},
					},
				},
			},
			EnvOpt: cel.Lib(&Library{
				PgrOpts: []cel.ProgramOption{cel.Globals(map[string]interface{}{"name": "name"})},
			}),
			WantErr: false,
		},
		{
			Name: "Overload function (error)",
			Expr: `getName() == "name"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"getName()": {Result: &Options_Overloads_Type{Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}}},
					},
				},
			},
			WantErr: true,
		},
		{
			Name: "Overload function",
			Expr: `getName() == "name"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"getName": {Result: &Options_Overloads_Type{Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}}},
					},
				},
			},
			EnvOpt: cel.Lib(&Library{
				EnvOpts: []cel.EnvOption{cel.Function("getName", cel.Overload("getName", []*cel.Type{}, cel.StringType, cel.FunctionBinding(func(values ...ref.Val) ref.Val { return types.String("name") })))},
				PgrOpts: []cel.ProgramOption{},
			}),
			WantErr: false,
		},
		{
			Name:    "Stdlib override (error)",
			Expr:    `type == "type"`,
			WantErr: true,
		},
		{
			Name: "Stdlib override (const)",
			Expr: `type == "type"`,
			Config: &Options{
				Globals: &Options_Globals{
					Constants: map[string]string{"type": "name"},
				},
				StdlibOverridingEnabled: true,
			},
			WantErr: false,
		},
		{
			Name: "Stdlib override (variable)",
			Expr: `type == "type"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Variables: map[string]*Options_Overloads_Type{"type": {
						Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}},
					},
				},
				StdlibOverridingEnabled: true,
			},
			EnvOpt: cel.Lib(&Library{
				PgrOpts: []cel.ProgramOption{cel.Globals(map[string]interface{}{"type": "name"})},
			}),
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			opts := []cel.EnvOption{BuildEnvOption(tt.Config, tt.Desc)}
			if tt.EnvOpt != nil {
				opts = append(opts, tt.EnvOpt)
			}
			env, err := cel.NewCustomEnv(opts...)
			if err == nil {
				ast, issues := env.Compile(tt.Expr)
				if issues != nil {
					err = issues.Err()
				} else {
					pgr, perr := env.Program(ast, cel.EvalOptions(cel.OptOptimize))
					if perr != nil {
						err = perr
					} else {
						_, _, err = pgr.Eval(map[string]interface{}{})
					}
				}
			}
			if (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

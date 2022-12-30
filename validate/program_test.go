package validate

import (
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestBuildValidateProgram(t *testing.T) {
	tests := []struct {
		Name      string
		Exprs     []string
		Config    *ValidateOptions
		Desc      protoreflect.MessageDescriptor
		EnvOption cel.EnvOption
		Imports   []protoreflect.FileDescriptor
		WantErr   bool
	}{
		{
			Name:    "Unknown field",
			Exprs:   []string{`name`},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: true,
		},
		{
			Name:    "Invalid return type",
			Exprs:   []string{`"name""`},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: true,
		},
		{
			Name:    "Invalid validate call on standard type",
			Exprs:   []string{`ref.validate()`},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			WantErr: true,
		},
		{
			Name:  "Unknown field in macro",
			Exprs: []string{`macro()`},
			Desc:  (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Functions: map[string]string{
							"macro": `name == "name"`,
						},
					},
				},
			},
			WantErr: true,
		},
		{
			Name:  "Regexp error",
			Exprs: []string{`ref.matches("[")`},
			Desc:  (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Config: &ValidateOptions{
				Options: &options.Options{
					Overloads: &options.Options_Overloads{
						Variables: map[string]*options.Options_Overloads_Type{
							"myVariable": {Type: &options.Options_Overloads_Type_Primitive_{
								Primitive: options.Options_Overloads_Type_STRING,
							}},
						},
					},
				},
			},
			WantErr: true,
		},
		{
			Name:    "OK",
			Exprs:   []string{`ref == "ref"`},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Config:  nil,
			WantErr: false,
		},
		{
			Name:  "OK (with constant)",
			Exprs: []string{`ref == constRef`},
			Desc:  (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"constRef": "ref",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:  "OK (with macro)",
			Exprs: []string{`rule() == ref`},
			Desc:  (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Functions: map[string]string{
							"rule": `ref`,
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:  "OK (with variable)",
			Exprs: []string{`ref == myVariable`},
			Desc:  (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Config: &ValidateOptions{
				Options: &options.Options{
					Overloads: &options.Options_Overloads{
						Variables: map[string]*options.Options_Overloads_Type{
							"myVariable": {Type: &options.Options_Overloads_Type_Primitive_{
								Primitive: options.Options_Overloads_Type_STRING,
							}},
						},
					},
				},
			},
			EnvOption: cel.Lib(&options.Library{
				PgrOpts: []cel.ProgramOption{cel.Globals(map[string]interface{}{"myVariable": "ref"})},
			}),
			WantErr: false,
		},
		{
			Name:    "OK (validate nested)",
			Exprs:   []string{`nested.validate()`},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Config:  nil,
			WantErr: false,
		},
		{
			Name:    "OK (validateWithMask nested)",
			Exprs:   []string{`nested.validateWithMask(fm)`},
			Desc:    (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Config:  nil,
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			lib := &options.Library{}
			if tt.EnvOption != nil {
				lib.EnvOpts = append(lib.EnvOpts, tt.EnvOption)
			}
			lib.EnvOpts = append(lib.EnvOpts, cel.DeclareContextProto(tt.Desc))
			lib.EnvOpts = append(lib.EnvOpts, buildValidatersFunctions(tt.Desc)...)
			if tt.Config == nil {
				lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(nil, tt.Desc))
			} else {
				lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(tt.Config.Options, tt.Desc))
			}
			_, err := BuildValidateProgram(tt.Exprs, tt.Config, cel.Lib(lib), tt.Imports...)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

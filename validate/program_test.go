package validate

import (
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
)

var tests = []struct {
	Name    string
	Expr    string
	Config  *ValidateOptions
	WantErr bool
}{
	{
		Name:    "Unknown field",
		Expr:    `name`,
		Config:  nil,
		WantErr: true,
	},
	{
		Name:    "Invalid return type",
		Expr:    `ref`,
		Config:  nil,
		WantErr: true,
	},
	{
		Name:    "Invalid validate call on standard type",
		Expr:    `ref.validate()`,
		Config:  nil,
		WantErr: true,
	},
	{
		Name:    "OK",
		Expr:    `ref == "ref"`,
		Config:  nil,
		WantErr: false,
	},
	{
		Name: "OK (with constant)",
		Expr: `ref == constRef`,
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
		Name: "OK (with macro)",
		Expr: `rule() == "ref"`,
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
		Name:    "OK (validate nested)",
		Expr:    `nested.validate()`,
		Config:  nil,
		WantErr: false,
	},
	{
		Name:    "OK (validateWithMask nested)",
		Expr:    `nested.validateWithMask(fm)`,
		Config:  nil,
		WantErr: false,
	},
}

func TestBuildValidateProgramFromDesk(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := BuildValidateProgramFromDesc(tt.Expr, nil, validate.File_testdata_validate_test_proto.Messages().Get(0), tt.Config)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

func TestBuildValidateProgram(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := BuildValidateProgram(tt.Expr, &validate.TestRpcRequest{}, tt.Config)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

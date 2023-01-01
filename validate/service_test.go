package validate

import (
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
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
			Name: "Service level expr with local const",
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

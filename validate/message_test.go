package validate

import (
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestBuildMessageValidateProgram(t *testing.T) {
	tests := []struct {
		Name    string
		Desc    protoreflect.MessageDescriptor
		Config  *ValidateOptions
		EnvOpt  cel.EnvOption
		WantErr bool
	}{
		{
			Name:    "No validation",
			Desc:    validate.File_testdata_validate_message_proto.Messages().ByName("Message"),
			WantErr: false,
		},
		{
			Name:    "Message level expr",
			Desc:    validate.File_testdata_validate_message_proto.Messages().ByName("MessageExpr"),
			WantErr: false,
		},
		// {
		// 	Name:    "Message level expr with missing const",
		// 	Desc:    validate.File_testdata_validate_message_proto.Messages().ByName("MessageOptions"),
		// 	WantErr: true,
		// },
		{
			Name: "Message level expr with global const",
			Desc: validate.File_testdata_validate_message_proto.Messages().ByName("MessageOptions"),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"emptyName": "",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:    "Message level expr with local const",
			Desc:    validate.File_testdata_validate_message_proto.Messages().ByName("MessageLocalOptions"),
			WantErr: false,
		},
		{
			Name: "Message level expr with const conflict",
			Desc: validate.File_testdata_validate_message_proto.Messages().ByName("MessageLocalOptions"),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"emptyName": "",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:    "Field level expr",
			Desc:    validate.File_testdata_validate_field_proto.Messages().ByName("FieldExpr"),
			WantErr: false,
		},
		{
			Name:    "Field resource reference wrong",
			Desc:    validate.File_testdata_validate_field_proto.Messages().ByName("FieldReferenceWrong"),
			WantErr: true,
		},
		{
			Name:    "Field resource reference type",
			Desc:    validate.File_testdata_validate_field_proto.Messages().ByName("FieldReferenceType"),
			WantErr: false,
		},
		{
			Name:    "Field resource reference child type",
			Desc:    validate.File_testdata_validate_field_proto.Messages().ByName("FieldReferenceChild"),
			WantErr: false,
		},
		{
			Name:    "Field resource reference type and child type",
			Desc:    validate.File_testdata_validate_field_proto.Messages().ByName("FieldReferenceTypeAndChild"),
			WantErr: true,
		},
		{
			Name:    "Field repeated resource reference",
			Desc:    validate.File_testdata_validate_field_proto.Messages().ByName("FieldRepeatedReferenceType"),
			WantErr: false,
		},
		{
			Name:    "Field level expr with missing const",
			Desc:    validate.File_testdata_validate_field_proto.Messages().ByName("FieldOptions"),
			WantErr: true,
		},
		{
			Name: "Field level expr with global const",
			Desc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldOptions"),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"emptyName": "",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name:    "Field level expr with local const",
			Desc:    validate.File_testdata_validate_field_proto.Messages().ByName("FieldLocalOptions"),
			WantErr: false,
		},
		{
			Name: "Field level expr with const conflict",
			Desc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldLocalOptions"),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"emptyName": "",
						},
					},
				},
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := BuildMessageValidateProgram(tt.Config, tt.Desc, tt.EnvOpt, tt.Desc.ParentFile())
			if (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

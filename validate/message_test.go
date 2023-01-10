package validate

import (
	"context"
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/sourcecontextpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/typepb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestMessaveValidateWithMask(t *testing.T) {
	tests := []struct {
		Name      string
		Message   proto.Message
		Fieldmask *fieldmaskpb.FieldMask
		WantErr   bool
	}{
		{
			Name:    "Wrong FieldExpr.name",
			Message: &validate.FieldExpr{Name: "name"},
			WantErr: true,
		},
		{
			Name:      "Wrong FieldExpr.name without mask",
			Message:   &validate.FieldExpr{Name: "name"},
			Fieldmask: &fieldmaskpb.FieldMask{Paths: []string{}},
			WantErr:   false,
		},
		{
			Name:    "Good FieldExpr.name",
			Message: &validate.FieldExpr{Name: "nam"},
			WantErr: false,
		},
		{
			Name:      "Good FieldExpr.name with * mask",
			Message:   &validate.FieldExpr{Name: "nam"},
			Fieldmask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			WantErr:   false,
		},
		{
			Name:    "Missing FieldExpr.name (not required)",
			Message: &validate.FieldExpr{},
			WantErr: false,
		},
		{
			Name:      "Missing FieldRequired.name (required)",
			Message:   &validate.FieldRequired{},
			Fieldmask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			WantErr:   true,
		},
		{
			Name:      "Good FieldRequired.name (required)",
			Message:   &validate.FieldRequired{Name: "name"},
			Fieldmask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			WantErr:   false,
		},
		{
			Name:      "Wrong FieldsExpr.display_name",
			Message:   &validate.FieldsExpr{Name: "nam", DisplayName: "nam"},
			Fieldmask: &fieldmaskpb.FieldMask{Paths: []string{"name", "display_name"}},
			WantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			p, err := BuildMessageValidateProgram(nil, tt.Message.ProtoReflect().Descriptor(), nil, tt.Message.ProtoReflect().Descriptor().ParentFile())
			if err != nil {
				t.Error(err)
			}
			if err = p.ValidateWithMask(context.Background(), tt.Message, tt.Fieldmask); (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

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
			Name: "Message config expr with const conflict",
			Desc: validate.File_testdata_validate_message_proto.Messages().ByName("Message"),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"emptyName": "",
						},
					},
				},
				Rules: map[string]*ValidateRule{
					string(validate.File_testdata_validate_message_proto.Messages().ByName("Message").FullName()): {
						Expr: `name != ""`,
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
		{
			Name: "Field config expr with const conflict",
			Desc: validate.File_testdata_validate_message_proto.Messages().ByName("Message"),
			Config: &ValidateOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"emptyName": "",
						},
					},
				},
				Rules: map[string]*ValidateRule{
					string(validate.File_testdata_validate_message_proto.Messages().ByName("Message").Fields().ByName("name").FullName()): {
						Expr: `name != ""`,
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

func TestIsDefaultValue(t *testing.T) {
	tests := []struct {
		Name       string
		Message    proto.Message
		Descriptor protoreflect.FieldDescriptor
		IsDefault  bool
	}{
		{
			Name:       "Bool (default)",
			Message:    &wrapperspb.BoolValue{},
			Descriptor: (&wrapperspb.BoolValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Bytes (default)",
			Message:    &wrapperspb.BytesValue{},
			Descriptor: (&wrapperspb.BytesValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Bytes (zero length)",
			Message:    &wrapperspb.BytesValue{Value: []byte{}},
			Descriptor: (&wrapperspb.BytesValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Double (default)",
			Message:    &wrapperspb.DoubleValue{},
			Descriptor: (&wrapperspb.DoubleValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Float (default)",
			Message:    &wrapperspb.FloatValue{},
			Descriptor: (&wrapperspb.FloatValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Int32 (default)",
			Message:    &wrapperspb.Int32Value{},
			Descriptor: (&wrapperspb.Int32Value{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Int64 (default)",
			Message:    &wrapperspb.Int64Value{},
			Descriptor: (&wrapperspb.Int64Value{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "String (default)",
			Message:    &wrapperspb.StringValue{},
			Descriptor: (&wrapperspb.StringValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "UInt32 (default)",
			Message:    &wrapperspb.UInt32Value{},
			Descriptor: (&wrapperspb.UInt32Value{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "UInt64 (default)",
			Message:    &wrapperspb.UInt64Value{},
			Descriptor: (&wrapperspb.UInt64Value{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "List (default)",
			Message:    &structpb.ListValue{},
			Descriptor: (&structpb.ListValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "List (zero length)",
			Message:    &structpb.ListValue{Values: []*structpb.Value{}},
			Descriptor: (&structpb.ListValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Map (default)",
			Message:    &structpb.Struct{},
			Descriptor: (&structpb.Struct{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Map (zero length)",
			Message:    &structpb.Struct{Fields: make(map[string]*structpb.Value)},
			Descriptor: (&structpb.Struct{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Enum (default)",
			Message:    &typepb.Field{},
			Descriptor: (&typepb.Field{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  true,
		},
		{
			Name:       "Message (default)",
			Message:    &typepb.Enum{},
			Descriptor: (&typepb.Enum{}).ProtoReflect().Descriptor().Fields().Get(3),
			IsDefault:  true,
		},
		{
			Name:       "Bool (non default)",
			Message:    &wrapperspb.BoolValue{Value: true},
			Descriptor: (&wrapperspb.BoolValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "Bytes (non default)",
			Message:    &wrapperspb.BytesValue{Value: []byte{1}},
			Descriptor: (&wrapperspb.BytesValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "Double (non default)",
			Message:    &wrapperspb.DoubleValue{Value: 42},
			Descriptor: (&wrapperspb.DoubleValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "Float (non default)",
			Message:    &wrapperspb.FloatValue{Value: 42},
			Descriptor: (&wrapperspb.FloatValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "Int32 (non default)",
			Message:    &wrapperspb.Int32Value{Value: 42},
			Descriptor: (&wrapperspb.Int32Value{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "Int64 (non default)",
			Message:    &wrapperspb.Int64Value{Value: 42},
			Descriptor: (&wrapperspb.Int64Value{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "String (non default)",
			Message:    &wrapperspb.StringValue{Value: "42"},
			Descriptor: (&wrapperspb.StringValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "UInt32 (non default)",
			Message:    &wrapperspb.UInt32Value{Value: 42},
			Descriptor: (&wrapperspb.UInt32Value{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "UInt64 (non default)",
			Message:    &wrapperspb.UInt64Value{Value: 42},
			Descriptor: (&wrapperspb.UInt64Value{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "List (non default)",
			Message:    &structpb.ListValue{Values: []*structpb.Value{{}}},
			Descriptor: (&structpb.ListValue{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "Map (non default)",
			Message:    &structpb.Struct{Fields: map[string]*structpb.Value{"": {}}},
			Descriptor: (&structpb.Struct{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "Enum (non default)",
			Message:    &typepb.Field{Kind: typepb.Field_TYPE_BOOL},
			Descriptor: (&typepb.Field{}).ProtoReflect().Descriptor().Fields().Get(0),
			IsDefault:  false,
		},
		{
			Name:       "Message (non default)",
			Message:    &typepb.Enum{SourceContext: &sourcecontextpb.SourceContext{}},
			Descriptor: (&typepb.Enum{}).ProtoReflect().Descriptor().Fields().Get(3),
			IsDefault:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if res := isDefaultValue(tt.Message, tt.Descriptor); res != tt.IsDefault {
				t.Errorf("want %v, got %v", tt.IsDefault, res)
			}
		})
	}
}

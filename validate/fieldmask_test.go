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

func TestValidateWithMask(t *testing.T) {
	lib := &options.Library{EnvOpts: []cel.EnvOption{cel.DeclareContextProto((&validate.TestRpcRequest{}).ProtoReflect().Descriptor())}}
	lib.EnvOpts = append(lib.EnvOpts, buildValidatersFunctions((&validate.TestRpcRequest{}).ProtoReflect().Descriptor())...)
	lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(nil, (&validate.TestRpcRequest{}).ProtoReflect().Descriptor()))
	noDepthMap := map[string]*Program{}
	noDepthMap["nested"], _ = BuildValidateProgram([]string{`true`}, nil, cel.Lib(lib))
	fmNoDepthMap := map[string]*Program{}
	fmNoDepthMap["nested"], _ = BuildValidateProgram([]string{`nested.validateWithMask(fm)`}, nil, cel.Lib(lib))
	tests := []struct {
		Name          string
		ValidationMap map[string]*Program
		Tests         []struct {
			Name    string
			Message proto.Message
			WantErr bool
		}
	}{
		{
			Name:          "No depth",
			ValidationMap: noDepthMap,
			Tests: []struct {
				Name    string
				Message proto.Message
				WantErr bool
			}{
				{
					Name:    "Empty ref",
					Message: &validate.TestRpcRequest{},
					WantErr: true,
				},
				{
					Name: "OK (ref)",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
					},
					WantErr: false,
				},
			},
		},
		{
			Name:          "Fieldmask validation without depth",
			ValidationMap: fmNoDepthMap,
			Tests: []struct {
				Name    string
				Message proto.Message
				WantErr bool
			}{
				{
					Name: "No fieldmask",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
					},
					WantErr: false,
				},
				{
					Name: "Fieldmask with empty struct fields",
					Message: &validate.TestRpcRequest{
						Ref:    "ref",
						Nested: &validate.Nested{},
						Fm:     &fieldmaskpb.FieldMask{Paths: []string{"name"}},
					},
					WantErr: true,
				},
				{
					Name: "Fieldmask with one invalid field",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
						Nested: &validate.Nested{
							Name: "name",
						},
						Fm: &fieldmaskpb.FieldMask{Paths: []string{"name", "value"}},
					},
					WantErr: true,
				},
				{
					Name: "Fieldmask with invalid nested field",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
						Nested: &validate.Nested{
							Ref: &validate.RefMultiple{
								Name: "noname",
							},
						},
						Fm: &fieldmaskpb.FieldMask{Paths: []string{"ref.name"}},
					},
					WantErr: true,
				},
				{
					Name: "Fieldmask with valid fields",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
						Nested: &validate.Nested{
							Name:  "name",
							Value: "value",
							Ref: &validate.RefMultiple{
								Name: "name",
							},
						},
						Fm: &fieldmaskpb.FieldMask{Paths: []string{"name", "value", "ref.name"}},
					},
					WantErr: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for _, ttt := range tt.Tests {
				t.Run(ttt.Name, func(t *testing.T) {
					err := ValidateWithMask(context.Background(), ttt.Message, &fieldmaskpb.FieldMask{Paths: []string{"*"}}, tt.ValidationMap, true)
					if (err != nil && !ttt.WantErr) || (err == nil && ttt.WantErr) {
						t.Errorf("wantErr %v, got %v", ttt.WantErr, err)
					}
				})
			}
		})
	}
	if err := ValidateWithMask(context.Background(), &validate.TestRpcRequest{Ref: "ref"}, nil, noDepthMap, true); err != nil {
		t.Errorf("wantErr false, got %v", err)
	}
	if err := ValidateWithMask(context.Background(), &validate.TestRpcRequest{Ref: "ref"}, &fieldmaskpb.FieldMask{Paths: []string{"ref", "nested"}}, noDepthMap, true); err != nil {
		t.Errorf("wantErr false, got %v", err)
	}
	if err := ValidateWithMask(context.Background(), nil, nil, nil, true); err == nil {
		t.Errorf("wantErr true, got <nil>")
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

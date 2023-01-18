package validate

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/sourcecontextpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/typepb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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
			if res := IsDefaultValue(tt.Message, tt.Descriptor); res != tt.IsDefault {
				t.Errorf("want %v, got %v", tt.IsDefault, res)
			}
		})
	}
}

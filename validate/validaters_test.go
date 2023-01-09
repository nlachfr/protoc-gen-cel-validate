package validate

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type optBuilder struct {
	got  map[string]bool
	want map[string]bool
}

func (b *optBuilder) Build(name string, t string) cel.FunctionOpt {
	if b.got == nil {
		b.got = map[string]bool{}
	}
	b.got[t] = true
	return buildValidateFunctionOpt(name, t)
}

func (b *optBuilder) Error() error {
	if len(b.got) == 0 && len(b.want) == 0 {
		return nil
	} else if !reflect.DeepEqual(b.got, b.want) {
		return fmt.Errorf("want %v, got %v", b.want, b.got)
	}
	return nil
}

func TestBuildValidaterFunctionOpts(t *testing.T) {
	tests := []struct {
		Name    string
		Message proto.Message
		Builder *optBuilder
	}{
		{
			Name:    "Message without rules",
			Message: &timestamppb.Timestamp{},
			Builder: &optBuilder{want: map[string]bool{}},
		},
		{
			Name:    "Field rules",
			Message: &validate.FieldExpr{},
			Builder: &optBuilder{want: map[string]bool{
				string((&validate.FieldExpr{}).ProtoReflect().Descriptor().FullName()): true,
			}},
		},
		{
			Name:    "Message rule",
			Message: &validate.MessageExpr{},
			Builder: &optBuilder{want: map[string]bool{
				string((&validate.MessageExpr{}).ProtoReflect().Descriptor().FullName()): true,
			}},
		},
		{
			Name:    "Message with rules, and nested message with rules",
			Message: &validate.MessageNestedExpr{},
			Builder: &optBuilder{want: map[string]bool{
				string((&validate.MessageNestedExpr{}).ProtoReflect().Descriptor().FullName()): true,
				string((&validate.MessageExpr{}).ProtoReflect().Descriptor().FullName()):       true,
			}},
		},
		{
			Name:    "Message without rules, and nested message with rules",
			Message: &validate.MessageNested{},
			Builder: &optBuilder{want: map[string]bool{
				string((&validate.MessageExpr{}).ProtoReflect().Descriptor().FullName()): true,
			}},
		},
		{
			Name:    "Recursive message",
			Message: &structpb.Struct{},
			Builder: &optBuilder{want: map[string]bool{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buildValidaterFunctionOpts("validate", tt.Builder.Build, tt.Message.ProtoReflect().Descriptor(), map[string]bool{})
			if err := tt.Builder.Error(); err != nil {
				t.Error(err)
			}
		})
	}
}

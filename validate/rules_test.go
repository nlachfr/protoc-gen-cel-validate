package validate

import (
	"context"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/nlachfr/protocel/options"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServiceRuleValidater(t *testing.T) {
	tests := []struct {
		Name             string
		Validater        func() ServiceRuleValidater
		AttributeContext *attribute_context.AttributeContext
		Request          proto.Message
		WantErr          bool
	}{
		{
			Name: "Method name mismatch",
			Validater: func() ServiceRuleValidater {
				mds := map[string]protoreflect.MethodDescriptor{}
				mrvs := map[string]MethodRuleValidater{}
				return &serviceRuleValidater{methodDescs: mds, methodRulesValidaters: mrvs}
			},
			AttributeContext: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{
					Operation: "unknown.Operation",
				},
			},
			WantErr: false,
		},
		{
			Name: "Headers map is nil",
			Validater: func() ServiceRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.TypeDescs(desc.ParentFile()),
						cel.TypeDescs(attribute_context.File_google_rpc_context_attribute_context_proto),
						cel.Variable("request", cel.ObjectType(string(desc.FullName()))),
						cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
						options.BuildEnvOption(nil),
					},
				}
				rv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `attribute_context.request.headers["ok"] == "ok"`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				mds := map[string]protoreflect.MethodDescriptor{}
				mrvs := map[string]MethodRuleValidater{}
				return &serviceRuleValidater{ruleValidater: rv, methodDescs: mds, methodRulesValidaters: mrvs}
			},
			AttributeContext: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &timestamppb.Timestamp{},
			WantErr: true,
		},
		{
			Name: "Nil request",
			Validater: func() ServiceRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.TypeDescs(desc.ParentFile()),
						cel.TypeDescs(attribute_context.File_google_rpc_context_attribute_context_proto),
						cel.Variable("request", cel.ObjectType(string(desc.FullName()))),
						cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
						options.BuildEnvOption(nil),
					},
				}
				rv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `request.getSeconds() == 1`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				mds := map[string]protoreflect.MethodDescriptor{}
				mrvs := map[string]MethodRuleValidater{
					"": &methodRuleValidater{validater: rv},
				}
				return &serviceRuleValidater{ruleValidater: nil, methodDescs: mds, methodRulesValidaters: mrvs}
			},
			AttributeContext: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: nil,
			WantErr: true,
		},
		{
			Name: "Attribute context validation failed",
			Validater: func() ServiceRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.TypeDescs(desc.ParentFile()),
						cel.TypeDescs(attribute_context.File_google_rpc_context_attribute_context_proto),
						cel.Variable("request", cel.ObjectType(string(desc.FullName()))),
						cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
						options.BuildEnvOption(nil),
					},
				}
				rv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `attribute_context.request.headers["ok"] == "ok"`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				mds := map[string]protoreflect.MethodDescriptor{}
				mrvs := map[string]MethodRuleValidater{}
				return &serviceRuleValidater{ruleValidater: rv, methodDescs: mds, methodRulesValidaters: mrvs}
			},
			AttributeContext: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
				Request: &attribute_context.AttributeContext_Request{
					Headers: map[string]string{"ok": ""},
				},
			},
			Request: nil,
			WantErr: true,
		},
		{
			Name: "Request validation failed",
			Validater: func() ServiceRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.TypeDescs(desc.ParentFile()),
						cel.TypeDescs(attribute_context.File_google_rpc_context_attribute_context_proto),
						cel.Variable("request", cel.ObjectType(string(desc.FullName()))),
						cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
						options.BuildEnvOption(nil),
					},
				}
				rv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `request.getSeconds() == 1`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				mds := map[string]protoreflect.MethodDescriptor{}
				mrvs := map[string]MethodRuleValidater{
					"": &methodRuleValidater{validater: rv},
				}
				return &serviceRuleValidater{ruleValidater: nil, methodDescs: mds, methodRulesValidaters: mrvs}
			},
			AttributeContext: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &timestamppb.Timestamp{},
			WantErr: true,
		},
		{
			Name: "OK",
			Validater: func() ServiceRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.TypeDescs(desc.ParentFile()),
						cel.TypeDescs(attribute_context.File_google_rpc_context_attribute_context_proto),
						cel.Variable("request", cel.ObjectType(string(desc.FullName()))),
						cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
						options.BuildEnvOption(nil),
					},
				}
				rv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `request.getSeconds() == 1`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				mds := map[string]protoreflect.MethodDescriptor{}
				mrvs := map[string]MethodRuleValidater{
					"": &methodRuleValidater{validater: rv},
				}
				return &serviceRuleValidater{ruleValidater: nil, methodDescs: mds, methodRulesValidaters: mrvs}
			},
			AttributeContext: &attribute_context.AttributeContext{
				Api: &attribute_context.AttributeContext_Api{},
			},
			Request: &timestamppb.Timestamp{Seconds: 1},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			v := tt.Validater()
			err := v.Validate(context.Background(), tt.AttributeContext, tt.Request)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

func TestMessageRuleValidater(t *testing.T) {
	tests := []struct {
		Name          string
		Validater     func() MessageRuleValidater
		Request       proto.Message
		FieldMask     *fieldmaskpb.FieldMask
		HasValidaters bool
		WantErr       bool
	}{
		{
			Name: "Field rule failure",
			Validater: func() MessageRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.DeclareContextProto(desc),
						options.BuildEnvOption(nil),
					},
				}
				frvs := map[string]FieldRuleValidater{}
				frv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `seconds > 10`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				frvs["seconds"] = &fieldRuleValidater{validater: frv}
				return &messageRuleValidater{ruleValidater: nil, fieldRulesValidaters: frvs}
			},
			HasValidaters: true,
			Request:       &timestamppb.Timestamp{Seconds: 1, Nanos: 5},
			FieldMask:     &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			WantErr:       true,
		},
		{
			Name: "Field rule failure (required)",
			Validater: func() MessageRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.DeclareContextProto(desc),
						options.BuildEnvOption(nil),
					},
				}
				frvs := map[string]FieldRuleValidater{}
				frv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `seconds > 10`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				frvs["seconds"] = &fieldRuleValidater{validater: frv, required: true}
				return &messageRuleValidater{ruleValidater: nil, fieldRulesValidaters: frvs}
			},
			HasValidaters: true,
			Request:       &timestamppb.Timestamp{Nanos: 5},
			FieldMask:     &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			WantErr:       true,
		},
		{
			Name: "Message rule failure",
			Validater: func() MessageRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.DeclareContextProto(desc),
						options.BuildEnvOption(nil),
					},
				}
				rv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `false`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				return &messageRuleValidater{ruleValidater: rv}
			},
			HasValidaters: true,
			Request:       &timestamppb.Timestamp{Nanos: 5},
			FieldMask:     &fieldmaskpb.FieldMask{Paths: []string{"*"}},
			WantErr:       true,
		},
		{
			Name: "OK (field rule with nil fieldmask)",
			Validater: func() MessageRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.DeclareContextProto(desc),
						options.BuildEnvOption(nil),
					},
				}
				frvs := map[string]FieldRuleValidater{}
				frv1, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `seconds > 1`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				frv2, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `nanos > 1`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				frvs["seconds"] = &fieldRuleValidater{validater: frv1}
				frvs["nanos"] = &fieldRuleValidater{validater: frv2}
				return &messageRuleValidater{ruleValidater: nil, fieldRulesValidaters: frvs}
			},
			HasValidaters: true,
			Request:       &timestamppb.Timestamp{Seconds: 10},
		},
		{
			Name: "OK (field rule with specified fieldmask)",
			Validater: func() MessageRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.DeclareContextProto(desc),
						options.BuildEnvOption(nil),
					},
				}
				frvs := map[string]FieldRuleValidater{}
				frv1, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `seconds > 1`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				frv2, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `nanos > 1`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				frvs["seconds"] = &fieldRuleValidater{validater: frv1}
				frvs["nanos"] = &fieldRuleValidater{validater: frv2}
				return &messageRuleValidater{ruleValidater: nil, fieldRulesValidaters: frvs}
			},
			HasValidaters: true,
			Request:       &timestamppb.Timestamp{Seconds: 10},
			FieldMask:     &fieldmaskpb.FieldMask{Paths: []string{"seconds"}},
		},
		{
			Name: "OK (message rule only)",
			Validater: func() MessageRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.DeclareContextProto(desc),
						options.BuildEnvOption(nil),
					},
				}
				rv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `nanos < 10`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				return &messageRuleValidater{ruleValidater: rv}
			},
			HasValidaters: true,
			Request:       &timestamppb.Timestamp{Nanos: 5},
			FieldMask:     &fieldmaskpb.FieldMask{Paths: []string{"*"}},
		},
		{
			Name: "OK (message rule and field rules)",
			Validater: func() MessageRuleValidater {
				desc := (&timestamppb.Timestamp{}).ProtoReflect().Descriptor()
				lib := &options.Library{
					EnvOpts: []cel.EnvOption{
						cel.DeclareContextProto(desc),
						options.BuildEnvOption(nil),
					},
				}
				rv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `nanos < 10`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				frvs := map[string]FieldRuleValidater{}
				frv, err := BuildRuleValidater(&Rule{
					Programs: []*Rule_Program{{Expr: `seconds > 1`}},
				}, cel.Lib(lib))
				if err != nil {
					panic(err)
				}
				frvs["seconds"] = &fieldRuleValidater{validater: frv}
				return &messageRuleValidater{ruleValidater: rv, fieldRulesValidaters: frvs}
			},
			HasValidaters: true,
			Request:       &timestamppb.Timestamp{Seconds: 10, Nanos: 5},
			FieldMask:     &fieldmaskpb.FieldMask{Paths: []string{"*"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			v := tt.Validater()
			if tt.HasValidaters != v.HasValidaters() {
				t.Errorf("want %v, got %v", tt.HasValidaters, v.HasValidaters())
			}
			err := v.ValidateWithMask(context.Background(), tt.Request, tt.FieldMask)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

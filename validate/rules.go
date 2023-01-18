package validate

import (
	"context"
	"fmt"
	"strings"

	"github.com/Neakxs/protocel/validate/errors"
	"github.com/google/cel-go/common/types"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type ServiceRuleValidater interface {
	Validate(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error
}
type serviceRuleValidater struct {
	ruleValidater         RuleValidater
	methodDescs           map[string]protoreflect.MethodDescriptor
	methodRulesValidaters map[string]MethodRuleValidater
}

func (v *serviceRuleValidater) Validate(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error {
	if attr == nil || attr.Api == nil {
		return nil
	} else {
		req := map[string]interface{}{
			"attribute_context": attr,
		}
		if v.ruleValidater != nil {
			for _, pgr := range v.ruleValidater.Programs() {
				if val, _, err := pgr.Program.ContextEval(ctx, req); err != nil {
					return errors.Wrap(err, m, v.methodDescs[attr.Api.Operation], attr)
				} else if !types.IsBool(val) || !val.Value().(bool) {
					return errors.New(m, v.methodDescs[attr.Api.Operation], attr)
				}
			}
		}
		req["request"] = m
		if methodValidater, ok := v.methodRulesValidaters[attr.Api.Operation]; ok && methodValidater != nil {
			if validater := methodValidater.Validater(); validater != nil {
				for _, pgr := range validater.Programs() {
					if val, _, err := pgr.Program.ContextEval(ctx, req); err != nil {
						return errors.Wrap(err, m, v.methodDescs[attr.Api.Operation], attr)
					} else if !types.IsBool(val) || !val.Value().(bool) {
						return errors.New(m, v.methodDescs[attr.Api.Operation], attr)
					}
				}
			}
		}
	}
	return nil
}

type MethodRuleValidater interface {
	Validater() RuleValidater
}
type methodRuleValidater struct {
	validater RuleValidater
}

func (v *methodRuleValidater) Validater() RuleValidater { return v.validater }

type MessageRuleValidater interface {
	ValidateWithMask(ctx context.Context, m proto.Message, fm *fieldmaskpb.FieldMask) error
	HasValidaters() bool
}

type messageRuleValidater struct {
	ruleValidater        RuleValidater
	fieldRulesValidaters map[string]FieldRuleValidater
}

func (v *messageRuleValidater) ValidateWithMask(ctx context.Context, m proto.Message, fm *fieldmaskpb.FieldMask) error {
	if v.fieldRulesValidaters == nil && v.ruleValidater == nil {
		return fmt.Errorf("validation failed")
	}
	vars := map[string]interface{}{}
	for i := 0; i < m.ProtoReflect().Descriptor().Fields().Len(); i++ {
		field := m.ProtoReflect().Descriptor().Fields().Get(i)
		vars[field.TextName()] = m.ProtoReflect().Get(field)
	}
	pathsMap := map[string][]string{}
	mdesc := m.ProtoReflect().Descriptor()
	if fm == nil {
		m.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
			pathsMap[fd.TextName()] = []string{""}
			return true
		})
	} else if len(fm.Paths) == 1 && fm.Paths[0] == "*" {
		if v.ruleValidater != nil {
			for _, p := range v.ruleValidater.Programs() {
				if val, _, err := p.Program.ContextEval(ctx, vars); err != nil {
					return errors.Wrap(err, m, mdesc, nil)
				} else if !types.IsBool(val) || !val.Value().(bool) {
					return errors.New(m, mdesc, nil)
				}
			}
		}
		for i := 0; i < mdesc.Fields().Len(); i++ {
			fdesc := mdesc.Fields().Get(i)
			pathsMap[fdesc.TextName()] = []string{""}
		}
	} else {
		for i := 0; i < len(fm.Paths); i++ {
			parts := strings.SplitN(fm.Paths[i], ".", 2)
			k := parts[0]
			v, ok := pathsMap[k]
			if !ok {
				pathsMap[k] = []string{}
			}
			if len(parts) > 1 {
				pathsMap[k] = append(v, parts[1:]...)
			} else {
				pathsMap[k] = append(v, "")
			}
		}
	}
	for i := 0; i < mdesc.Fields().Len(); i++ {
		fdesc := mdesc.Fields().Get(i)
		if paths, ok := pathsMap[fdesc.TextName()]; ok {
			subs := []string{}
			for j := 0; j < len(paths); j++ {
				if paths[j] == "" {
					if fieldValidater, ok := v.fieldRulesValidaters[string(fdesc.Name())]; ok {
						if !IsDefaultValue(m, fdesc) {
							if fieldValidater.Validater() != nil {
								for _, p := range fieldValidater.Validater().Programs() {
									if val, _, err := p.Program.ContextEval(ctx, vars); err != nil {
										return errors.Wrap(err, m, fdesc, nil)
									} else if !types.IsBool(val) || !val.Value().(bool) {
										return errors.New(m, fdesc, nil)
									}
								}
							}
						} else if fieldValidater.IsRequired() {
							return errors.New(m, fdesc, nil)
						}
					}
				} else if paths[j] != "*" {
					subs = append(subs, paths[j])
				}
			}
			if len(subs) > 0 && fdesc.Kind() == protoreflect.MessageKind {
				if v, ok := m.ProtoReflect().Get(fdesc).Message().Interface().(Validater); ok {
					if err := v.ValidateWithMask(ctx, &fieldmaskpb.FieldMask{Paths: subs}); err != nil {
						return errors.Wrap(err, m, fdesc, nil)
					}
				}
			}
		}
	}
	return nil
}

func (v *messageRuleValidater) HasValidaters() bool {
	return v.ruleValidater != nil || len(v.fieldRulesValidaters) > 0
}

type FieldRuleValidater interface {
	Validater() RuleValidater
	IsRequired() bool
}

type fieldRuleValidater struct {
	validater RuleValidater
	required  bool
}

func (v *fieldRuleValidater) Validater() RuleValidater {
	return v.validater
}
func (v *fieldRuleValidater) IsRequired() bool {
	return v.required
}

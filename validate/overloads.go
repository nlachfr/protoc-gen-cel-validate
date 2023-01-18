package validate

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func buildOverloads(desc protoreflect.MessageDescriptor, validateUnary func(ref.Val) ref.Val, validateWithMaskBinary func(ref.Val, ref.Val) ref.Val) []cel.EnvOption {
	res := []cel.EnvOption{}
	if opts := buildFunctionOpts(desc, "validate", func(name, t string) cel.FunctionOpt {
		return cel.MemberOverload(
			fmt.Sprintf("%s_%s", t, name),
			[]*cel.Type{cel.ObjectType(t)},
			cel.BoolType,
			cel.UnaryBinding(validateUnary),
		)
	}); len(opts) > 0 {
		res = append(res, cel.Function("validate", opts...))
	}
	if opts := buildFunctionOpts(desc, "validateWithMask", func(name, t string) cel.FunctionOpt {
		return cel.MemberOverload(
			fmt.Sprintf("%s_%s", t, name),
			[]*cel.Type{cel.ObjectType(t), cel.ObjectType(string((&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName()))},
			cel.BoolType,
			cel.BinaryBinding(validateWithMaskBinary),
		)
	}); len(opts) > 0 {
		res = append(res, cel.Function("validateWithMask", opts...))
	}
	return res
}

func buildFunctionOpts(desc protoreflect.MessageDescriptor, name string, optBuilder func(name, t string) cel.FunctionOpt, m ...map[string]bool) []cel.FunctionOpt {
	if len(m) == 0 {
		m = append(m, map[string]bool{})
	}
	functionOpts := []cel.FunctionOpt{}
	messageType := string(desc.FullName())
	if _, ok := m[0][messageType]; !ok {
		m[0][messageType] = true
		buildValidate := false
		if GetExtension(desc.Options(), E_Message).(*MessageRule) != nil {
			buildValidate = true
		}
		for i := 0; i < desc.Fields().Len(); i++ {
			fd := desc.Fields().Get(i)
			if GetExtension(fd.Options(), E_Field).(*FieldRule) != nil {
				buildValidate = true
			}
			if fd.Kind() == protoreflect.MessageKind {
				functionOpts = append(functionOpts, buildFunctionOpts(fd.Message(), name, optBuilder, m...)...)
			}
		}
		if buildValidate {
			functionOpts = append(functionOpts, optBuilder(name, messageType))
		}
	}
	return functionOpts
}

type overloadBuilder interface {
	buildOverloads(desc protoreflect.MessageDescriptor) []cel.EnvOption
}

type defaultOverloadBuilder struct{}

func (b *defaultOverloadBuilder) buildOverloads(desc protoreflect.MessageDescriptor) []cel.EnvOption {
	return buildOverloads(desc, b.validate, b.validateWithMask)
}

func (b *defaultOverloadBuilder) validate(value ref.Val) ref.Val {
	var err error
	if v, ok := value.Value().(Validater); ok {
		err = v.Validate(context.TODO())
	} else {
		return types.Bool(false)
	}
	if err == nil {
		return types.Bool(true)
	} else if vErr, ok := err.(ref.Val); ok {
		return vErr
	} else {
		return types.NewErr(err.Error())
	}
}

func (b *defaultOverloadBuilder) validateWithMask(lhs, rhs ref.Val) ref.Val {
	var err error
	fm := rhs.Value().(*fieldmaskpb.FieldMask)
	if v, ok := lhs.Value().(Validater); ok {
		err = v.ValidateWithMask(context.TODO(), fm)
	} else {
		return types.Bool(false)
	}
	if err == nil {
		return types.Bool(true)
	} else if vErr, ok := err.(ref.Val); ok {
		return vErr
	} else {
		return types.NewErr(err.Error())
	}
}

type fallbackOverloadBuilder struct {
	Builder *builder
}

func (b *fallbackOverloadBuilder) buildOverloads(desc protoreflect.MessageDescriptor) []cel.EnvOption {
	return buildOverloads(desc, b.validate, b.validateWithMask)
}

func (b *fallbackOverloadBuilder) validate(value ref.Val) ref.Val {
	msg, ok := value.Value().(proto.Message)
	if ok {
		desc := msg.ProtoReflect().Descriptor()
		messageValidater, err := b.Builder.BuildMessageRuleValidater(desc)
		if err != nil {
			return types.NewErr(err.Error())
		}
		if err = messageValidater.ValidateWithMask(context.TODO(), msg, &fieldmaskpb.FieldMask{Paths: []string{"*"}}); err != nil {
			if vErr, ok := err.(ref.Val); ok {
				return vErr
			}
			return types.NewErr(err.Error())
		}
		return types.Bool(true)
	}
	return types.Bool(false)
}

func (b *fallbackOverloadBuilder) validateWithMask(lhs, rhs ref.Val) ref.Val {
	fm := rhs.Value().(*fieldmaskpb.FieldMask)
	msg, ok := lhs.Value().(proto.Message)
	if ok {
		desc := msg.ProtoReflect().Descriptor()
		messageValidater, err := b.Builder.BuildMessageRuleValidater(desc)
		if err != nil {
			return types.NewErr(err.Error())
		}
		if err = messageValidater.ValidateWithMask(context.TODO(), msg, fm); err != nil {
			if vErr, ok := err.(ref.Val); ok {
				return vErr
			}
			return types.NewErr(err.Error())
		}
		return types.Bool(true)
	}
	return types.Bool(false)
}

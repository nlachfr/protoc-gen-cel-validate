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

func buildValidatersFunctions(config *ValidateOptions, desc protoreflect.MessageDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) []cel.EnvOption {
	res := []cel.EnvOption{}
	builder := &validateOverloadBuilder{
		config:  config,
		envOpt:  envOpt,
		imports: imports,
	}
	if r := builder.buildValidateFunction(desc); r != nil {
		res = append(res, r)
	}
	if r := builder.buildValidateWithMaskFunction(desc); r != nil {
		res = append(res, r)
	}
	return res
}

type validateOverloadBuilder struct {
	config   *ValidateOptions
	envOpt   cel.EnvOption
	imports  []protoreflect.FileDescriptor
	fallback map[string]MessageValidateProgram
}

func (b *validateOverloadBuilder) buildValidateFunction(desc protoreflect.MessageDescriptor) cel.EnvOption {
	if opts := b.buildFunctionOpts(desc, "validate", b.ValidateFunctionOpt, map[string]bool{}); len(opts) > 0 {
		return cel.Function("validate", opts...)
	}
	return nil
}

func (b *validateOverloadBuilder) buildValidateWithMaskFunction(desc protoreflect.MessageDescriptor) cel.EnvOption {
	if opts := b.buildFunctionOpts(desc, "validateWithMask", b.ValidateWithMaskFunctionOpt, map[string]bool{}); len(opts) > 0 {
		return cel.Function("validateWithMask", opts...)
	}
	return nil
}

func (b *validateOverloadBuilder) buildFunctionOpts(desc protoreflect.MessageDescriptor, name string, optBuilder func(name, t string) cel.FunctionOpt, m map[string]bool) []cel.FunctionOpt {
	functionOpts := []cel.FunctionOpt{}
	messageType := string(desc.FullName())
	if _, ok := m[messageType]; !ok {
		m[messageType] = true
		buildValidate := false
		if proto.GetExtension(desc.Options(), E_Message).(*ValidateRule) != nil {
			buildValidate = true
		}
		for i := 0; i < desc.Fields().Len(); i++ {
			fd := desc.Fields().Get(i)
			if proto.GetExtension(fd.Options(), E_Field).(*ValidateRule) != nil {
				buildValidate = true
			}
			if fd.Kind() == protoreflect.MessageKind {
				functionOpts = append(functionOpts, b.buildFunctionOpts(fd.Message(), name, optBuilder, m)...)
			}
		}
		if buildValidate {
			functionOpts = append(functionOpts, optBuilder(name, messageType))
		}
	}
	fmt.Println(m)
	return functionOpts
}

func (b *validateOverloadBuilder) ValidateFunctionOpt(name, t string) cel.FunctionOpt {
	return cel.MemberOverload(
		fmt.Sprintf("%s_%s", t, name),
		[]*cel.Type{cel.ObjectType(t)},
		cel.BoolType,
		cel.UnaryBinding(func(value ref.Val) ref.Val {
			var err error
			if v, ok := value.Value().(Validater); ok {
				err = v.Validate(context.TODO())
			} else if msg, ok := value.Value().(proto.Message); ok {
				desc := msg.ProtoReflect().Descriptor()
				pgr, ok := b.fallback[string(desc.FullName())]
				if !ok {
					if fbPgr, err := BuildMessageValidateProgram(b.config, desc, b.envOpt, b.imports...); err != nil {
						return types.NewErr(err.Error())
					} else {
						pgr = fbPgr
					}
				}
				err = pgr.ValidateWithMask(context.TODO(), msg, &fieldmaskpb.FieldMask{Paths: []string{"*"}})
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
		}),
	)
}

func (b *validateOverloadBuilder) ValidateWithMaskFunctionOpt(name, t string) cel.FunctionOpt {
	return cel.MemberOverload(
		fmt.Sprintf("%s_%s", t, name),
		[]*cel.Type{cel.ObjectType(t), cel.ObjectType(string((&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName()))},
		cel.BoolType,
		cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
			var err error
			fm, ok := rhs.Value().(*fieldmaskpb.FieldMask)
			if !ok {
				return types.Bool(false)
			}
			if v, ok := lhs.Value().(Validater); ok {
				err = v.ValidateWithMask(context.TODO(), fm)
			} else if msg, ok := lhs.Value().(proto.Message); ok {
				desc := msg.ProtoReflect().Descriptor()
				pgr, ok := b.fallback[string(desc.FullName())]
				if !ok {
					if fbPgr, err := BuildMessageValidateProgram(b.config, desc, b.envOpt, b.imports...); err != nil {
						return types.NewErr(err.Error())
					} else {
						pgr = fbPgr
					}
				}
				err = pgr.ValidateWithMask(context.TODO(), msg, fm)
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
		}),
	)
}

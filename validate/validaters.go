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

func buildValidateFunctionOpt(name string, t string) cel.FunctionOpt {
	return cel.MemberOverload(
		fmt.Sprintf("%s_%s", t, name),
		[]*cel.Type{cel.ObjectType(t)},
		cel.BoolType,
		cel.UnaryBinding(func(value ref.Val) ref.Val {
			if v, ok := value.Value().(Validater); ok {
				if err := v.Validate(context.TODO()); err == nil {
					return types.Bool(true)
				} else if validationErr, ok := err.(*ValidationError); ok {
					return validationErr
				} else {
					return types.NewErr(err.Error())
				}
			}
			return types.Bool(false)
		}),
	)
}

func buildValidateWithMaskFunctionOpt(name string, t string) cel.FunctionOpt {
	return cel.MemberOverload(
		fmt.Sprintf("%s_%s", t, name),
		[]*cel.Type{cel.ObjectType(t), cel.ObjectType(string((&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName()))},
		cel.BoolType,
		cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
			if v, ok := lhs.Value().(Validater); ok {
				if fm, ok := rhs.Value().(*fieldmaskpb.FieldMask); ok {
					if err := v.ValidateWithMask(context.TODO(), fm); err == nil {
						return types.Bool(true)
					} else if validationErr, ok := err.(*ValidationError); ok {
						return validationErr
					} else {
						return types.NewErr(err.Error())
					}
				}
			}
			return types.Bool(false)
		}),
	)
}

func buildValidatersFunctions(desc protoreflect.MessageDescriptor) []cel.EnvOption {
	res := []cel.EnvOption{}
	if r := buildValidaterFunction("validate", buildValidateFunctionOpt, desc); r != nil {
		res = append(res, r)
	}
	if r := buildValidaterFunction("validateWithMask", buildValidateWithMaskFunctionOpt, desc); r != nil {
		res = append(res, r)
	}
	return res
}

func buildValidaterFunction(name string, optBuilder func(name string, t string) cel.FunctionOpt, desc protoreflect.MessageDescriptor) cel.EnvOption {
	opts := buildValidaterFunctionOpts(name, optBuilder, desc, map[string]bool{})
	if len(opts) > 0 {
		return cel.Function(name, opts...)
	}
	return nil
}

func buildValidaterFunctionOpts(name string, optBuilder func(name string, t string) cel.FunctionOpt, desc protoreflect.MessageDescriptor, m map[string]bool) []cel.FunctionOpt {
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
				functionOpts = append(functionOpts, buildValidaterFunctionOpts(name, optBuilder, fd.Message(), m)...)
			}
		}
		if buildValidate {
			functionOpts = append(functionOpts, optBuilder(name, messageType))
		}
	}
	return functionOpts
}

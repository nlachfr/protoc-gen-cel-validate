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
	return []cel.EnvOption{
		buildValidaterFunction("validate", buildValidateFunctionOpt, desc),
		buildValidaterFunction("validateWithMask", buildValidateWithMaskFunctionOpt, desc),
	}
}

func buildValidaterFunction(name string, optBuilder func(name string, t string) cel.FunctionOpt, desc protoreflect.MessageDescriptor) cel.EnvOption {
	return cel.Function(name, buildValidaterFunctionOpts(name, optBuilder, desc, map[string]bool{})...)
}

func buildValidaterFunctionOpts(name string, optBuilder func(name string, t string) cel.FunctionOpt, desc protoreflect.MessageDescriptor, m map[string]bool) []cel.FunctionOpt {
	functionOpts := []cel.FunctionOpt{}
	for i := 0; i < desc.Fields().Len(); i++ {
		fd := desc.Fields().Get(i)
		if proto.GetExtension(fd.Options(), E_Field) != nil {
			t := string(desc.FullName())
			if _, ok := m[t]; !ok {
				functionOpts = append(functionOpts, optBuilder(name, t))
				m[t] = true
			}
		}
		if fd.Kind() == protoreflect.MessageKind {
			functionOpts = append(functionOpts, buildValidaterFunctionOpts(name, optBuilder, fd.Message(), m)...)
		}
	}
	return functionOpts
}
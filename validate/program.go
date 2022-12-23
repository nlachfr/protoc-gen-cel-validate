package validate

import (
	"context"
	"fmt"

	options "github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func BuildValidateProgram(expr string, config *ValidateOptions, desc protoreflect.MessageDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (cel.Program, error) {
	envOpts := []cel.EnvOption{
		cel.Types(&fieldmaskpb.FieldMask{}),
		cel.DeclareContextProto(desc),
	}
	if envOpt != nil {
		envOpts = append(envOpts, envOpt)
	}
	for _, imp := range imports {
		envOpts = append(envOpts, cel.TypeDescs(imp))
	}
	if config != nil {
		envOpts = append(envOpts, options.BuildEnvOption(config.Options, desc))
		if config.Options != nil {
			if macros, err := options.BuildMacros(config.Options, expr, envOpts); err != nil {
				return nil, fmt.Errorf("build macros error: %v", err)
			} else {
				envOpts = append(envOpts, cel.Macros(macros...))
			}
		}
	} else {
		envOpts = append(envOpts, options.BuildEnvOption(nil, desc))
	}
	envOpts = append(envOpts, buildValidatersFunctions(desc)...)
	env, err := cel.NewCustomEnv(envOpts...)
	if err != nil {
		return nil, fmt.Errorf("new env error: %w", err)
	}
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("compile error: %w", issues.Err())
	}
	if !ast.OutputType().IsAssignableType(cel.BoolType) {
		return nil, fmt.Errorf("output type not bool")
	}
	pgr, err := env.Program(ast, cel.OptimizeRegex(interpreter.MatchesRegexOptimization))
	if err != nil {
		return nil, fmt.Errorf("program error: %w", err)
	}
	return pgr, nil
}

func buildValidatersFunctions(desc protoreflect.MessageDescriptor) []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function("validate", buildValidatersValidateFunctionOpts(map[string]bool{}, desc)...),
		cel.Function("validateWithMask", buildValidatersValidateWithMaskFunctionOpts(map[string]bool{}, desc)...),
	}
}

func buildValidatersValidateFunctionOpts(m map[string]bool, desc protoreflect.MessageDescriptor) []cel.FunctionOpt {
	functionOpts := []cel.FunctionOpt{}
	for i := 0; i < desc.Fields().Len(); i++ {
		fd := desc.Fields().Get(i)
		if proto.GetExtension(fd.Options(), E_Field) != nil {
			t := string(desc.FullName())
			if _, ok := m[t]; !ok {
				functionOpts = append(functionOpts, cel.MemberOverload(
					fmt.Sprintf("%s_validate", t),
					[]*cel.Type{cel.ObjectType(t)},
					cel.BoolType,
					cel.UnaryBinding(func(value ref.Val) ref.Val {
						if v, ok := value.Value().(Validater); ok {
							if err := v.Validate(context.TODO()); err == nil {
								return types.Bool(true)
							} else {
								return types.NewErr(err.Error())
							}
						}
						return types.Bool(false)
					}),
				))
				m[t] = true
			}
		}
		if fd.Kind() == protoreflect.MessageKind {
			functionOpts = append(functionOpts, buildValidatersValidateFunctionOpts(m, fd.Message())...)
		}
	}
	return functionOpts
}

func buildValidatersValidateWithMaskFunctionOpts(m map[string]bool, desc protoreflect.MessageDescriptor) []cel.FunctionOpt {
	functionOpts := []cel.FunctionOpt{}
	for i := 0; i < desc.Fields().Len(); i++ {
		fd := desc.Fields().Get(i)
		if proto.GetExtension(fd.Options(), E_Field) != nil {
			t := string(desc.FullName())
			if _, ok := m[t]; !ok {
				functionOpts = append(functionOpts, cel.MemberOverload(
					fmt.Sprintf("%s_validateWithMask", t),
					[]*cel.Type{cel.ObjectType(t), cel.ObjectType(string((&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName()))},
					cel.BoolType,
					cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
						if v, ok := lhs.Value().(Validater); ok {
							if fm, ok := rhs.Value().(*fieldmaskpb.FieldMask); ok {
								if err := v.ValidateWithMask(context.TODO(), fm); err == nil {
									return types.Bool(true)
								} else {
									return types.NewErr(err.Error())
								}
							}
						}
						return types.Bool(false)
					}),
				))
				m[t] = true
			}
		}
		if fd.Kind() == protoreflect.MessageKind {
			functionOpts = append(functionOpts, buildValidatersValidateWithMaskFunctionOpts(m, fd.Message())...)
		}
	}
	return functionOpts
}

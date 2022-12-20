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

func BuildValidateProgramFromDesc(expr string, imports []protoreflect.FileDescriptor, msgDesc protoreflect.MessageDescriptor, config *ValidateOptions, libs ...cel.Library) (cel.Program, error) {
	envOpts := []cel.EnvOption{
		cel.TypeDescs(msgDesc.Parent()),
	}
	for i := 0; i < len(imports); i++ {
		envOpts = append(envOpts, cel.TypeDescs(imports[i]))
	}
	for i := 0; i < len(libs); i++ {
		envOpts = append(envOpts, cel.Lib(libs[i]))
	}
	return buildValidateProgram(expr, msgDesc, config, envOpts...)
}

func BuildValidateProgram(expr string, msg proto.Message, config *ValidateOptions, libs ...cel.Library) (cel.Program, error) {
	envOpts := []cel.EnvOption{
		cel.Types(msg),
	}
	for i := 0; i < len(libs); i++ {
		envOpts = append(envOpts, cel.Lib(libs[i]))
	}
	return buildValidateProgram(expr, msg.ProtoReflect().Descriptor(), config, envOpts...)
}

func buildValidateProgram(expr string, desc protoreflect.MessageDescriptor, config *ValidateOptions, envOpts ...cel.EnvOption) (cel.Program, error) {
	envOpts = append(envOpts,
		cel.Types(&fieldmaskpb.FieldMask{}),
		cel.DeclareContextProto(desc),
	)
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

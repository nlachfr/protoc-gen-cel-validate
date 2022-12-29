package validate

import (
	"fmt"

	options "github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/interpreter"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Program struct {
	rules []cel.Program
}

func BuildValidateProgram(exprs []string, config *ValidateOptions, desc protoreflect.MessageDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (*Program, error) {
	envOpts := []cel.EnvOption{cel.Types(&fieldmaskpb.FieldMask{})}
	if envOpt != nil {
		envOpts = append(envOpts, envOpt)
	}
	for _, imp := range imports {
		envOpts = append(envOpts, cel.TypeDescs(imp))
	}
	if desc != nil {
		envOpts = append(envOpts, cel.DeclareContextProto(desc))
		envOpts = append(envOpts, buildValidatersFunctions(desc)...)
		if msgOptions := proto.GetExtension(desc.Options(), E_Message).(*ValidateRule); msgOptions != nil {
			if config == nil {
				config = &ValidateOptions{}
			}
			if config.Options != nil {
				proto.Merge(config.Options, msgOptions.Options)
			} else {
				config.Options = msgOptions.Options
			}
		}
	}
	pgrs := []cel.Program{}
	for _, expr := range exprs {
		customEnvOpts := envOpts
		if config != nil {
			customEnvOpts = append(customEnvOpts, options.BuildEnvOption(config.Options, desc))
			if config.Options != nil {
				if macros, err := options.BuildMacros(config.Options, expr, customEnvOpts); err != nil {
					return nil, fmt.Errorf("build macros error: %v", err)
				} else {
					customEnvOpts = append(customEnvOpts, cel.Macros(macros...))
				}
			}
		} else {
			customEnvOpts = append(customEnvOpts, options.BuildEnvOption(nil, desc))
		}
		env, err := cel.NewCustomEnv(customEnvOpts...)
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
		pgrs = append(pgrs, pgr)
	}
	return &Program{rules: pgrs}, nil
}

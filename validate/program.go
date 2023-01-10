package validate

import (
	"context"
	"fmt"

	options "github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Validater interface {
	Validate(ctx context.Context) error
	ValidateWithMask(ctx context.Context, fm *fieldmaskpb.FieldMask) error
}

type ValidateProgram interface {
	CEL() map[string]cel.Program
}

type validateProgram struct {
	rules map[string]cel.Program
}

func (p *validateProgram) CEL() map[string]cel.Program { return p.rules }

func BuildValidateProgram(exprs []string, config *ValidateOptions, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (ValidateProgram, error) {
	envOpts := []cel.EnvOption{cel.Types(&fieldmaskpb.FieldMask{})}
	if envOpt != nil {
		envOpts = append(envOpts, envOpt)
	}
	for _, imp := range imports {
		envOpts = append(envOpts, cel.TypeDescs(imp))
	}
	rules := map[string]cel.Program{}
	for _, expr := range exprs {
		customEnvOpts := envOpts
		if config != nil && config.Options != nil {
			if macros, err := options.BuildMacros(config.Options, expr, customEnvOpts); err != nil {
				return nil, fmt.Errorf("build macros error: %v", err)
			} else {
				customEnvOpts = append(customEnvOpts, cel.Macros(macros...))
			}
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
		pgr, err := env.Program(ast, cel.EvalOptions(cel.OptOptimize))
		if err != nil {
			return nil, fmt.Errorf("program error: %w", err)
		}
		rules[expr] = pgr
	}
	return &validateProgram{rules: rules}, nil
}
